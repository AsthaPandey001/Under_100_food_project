package authrepo

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/profile"
	"github.com/mikestefanello/pagoda/ent/user"
	"github.com/mikestefanello/pagoda/pkg/services/sms"
)

var (
	ErrOTPExpired             = errors.New("OTP verification code has expired")
	ErrOTPInvalid             = errors.New("invalid OTP code. Please check and try again")
	ErrOTPMaxAttemptsExceeded = errors.New("maximum verification attempts exceeded. Please request a new OTP")
	ErrOTPCooldown            = errors.New("please wait 60 seconds before requesting another OTP")
	ErrOTPNotFound            = errors.New("no active OTP request found for this phone number")
)

// PhoneOTPRepo manages cryptographically secure hashed OTP verification challenges
type PhoneOTPRepo struct {
	db  *sql.DB
	orm *ent.Client
}

// NewPhoneOTPRepo creates a new PhoneOTPRepo
func NewPhoneOTPRepo(db *sql.DB, orm *ent.Client) *PhoneOTPRepo {
	return &PhoneOTPRepo{
		db:  db,
		orm: orm,
	}
}

// CreateChallenge generates a 6-digit OTP, stores its SHA-256 hash with 5-minute expiry, and enforces a 60s cooldown
func (r *PhoneOTPRepo) CreateChallenge(ctx context.Context, phoneInput, portal, ipAddress string) (string, string, error) {
	phoneE164 := sms.NormalizeE164(phoneInput)
	if len(phoneE164) < 11 {
		return "", "", fmt.Errorf("invalid phone number format")
	}

	if portal == "" {
		portal = "customer"
	}

	// 1. Enforce 60-second cooldown
	var recentCount int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM phone_otp_challenges 
		WHERE phone_e164 = $1 AND portal = $2 AND created_at > NOW() - INTERVAL '60 seconds'
	`, phoneE164, portal).Scan(&recentCount)
	if err == nil && recentCount > 0 {
		return "", "", ErrOTPCooldown
	}

	// 2. Invalidate any unconsumed previous challenges for this phone & portal
	_, _ = r.db.ExecContext(ctx, `
		UPDATE phone_otp_challenges 
		SET consumed_at = NOW() 
		WHERE phone_e164 = $1 AND portal = $2 AND consumed_at IS NULL
	`, phoneE164, portal)

	// 3. Generate cryptographically secure 6-digit OTP (100000 - 999999)
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	var code string
	if err != nil {
		code = fmt.Sprintf("%06d", (time.Now().UnixNano()%900000)+100000)
	} else {
		code = fmt.Sprintf("%06d", n.Int64()+100000)
	}

	// 4. Store only SHA-256 hash in database
	otpHash := HashToken(code)
	expiresAt := time.Now().Add(5 * time.Minute)

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO phone_otp_challenges (phone_e164, portal, otp_hash, attempts, max_attempts, expires_at, created_at, ip_address)
		VALUES ($1, $2, $3, 0, 5, $4, NOW(), $5)
	`, phoneE164, portal, otpHash, expiresAt, ipAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to store OTP challenge: %w", err)
	}

	return phoneE164, code, nil
}

// VerifyChallenge verifies the entered OTP against stored hash, enforcing attempt limits and expiration
func (r *PhoneOTPRepo) VerifyChallenge(ctx context.Context, phoneInput, portal, inputCode string) error {
	phoneE164 := sms.NormalizeE164(phoneInput)
	cleanCode := strings.TrimSpace(inputCode)
	if cleanCode == "" {
		return ErrOTPInvalid
	}

	if portal == "" {
		portal = "customer"
	}

	var challengeID int64
	var storedHash string
	var attempts int
	var maxAttempts int
	var expiresAt time.Time
	var consumedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, otp_hash, attempts, max_attempts, expires_at, consumed_at 
		FROM phone_otp_challenges 
		WHERE phone_e164 = $1 AND portal = $2 AND consumed_at IS NULL 
		ORDER BY id DESC LIMIT 1
	`, phoneE164, portal).Scan(&challengeID, &storedHash, &attempts, &maxAttempts, &expiresAt, &consumedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOTPNotFound
		}
		return fmt.Errorf("database error checking OTP challenge: %w", err)
	}

	if consumedAt.Valid {
		return ErrOTPNotFound
	}

	if time.Now().After(expiresAt) {
		_, _ = r.db.ExecContext(ctx, `UPDATE phone_otp_challenges SET consumed_at = NOW() WHERE id = $1`, challengeID)
		return ErrOTPExpired
	}

	if attempts >= maxAttempts {
		_, _ = r.db.ExecContext(ctx, `UPDATE phone_otp_challenges SET consumed_at = NOW() WHERE id = $1`, challengeID)
		return ErrOTPMaxAttemptsExceeded
	}

	// Compare hashes with constant-time equality check to prevent timing attacks
	inputHash := HashToken(cleanCode)
	match := subtle.ConstantTimeCompare([]byte(storedHash), []byte(inputHash)) == 1

	if !match {
		// Increment attempts
		newAttempts := attempts + 1
		if newAttempts >= maxAttempts {
			_, _ = r.db.ExecContext(ctx, `UPDATE phone_otp_challenges SET attempts = $1, consumed_at = NOW() WHERE id = $2`, newAttempts, challengeID)
			return ErrOTPMaxAttemptsExceeded
		}
		_, _ = r.db.ExecContext(ctx, `UPDATE phone_otp_challenges SET attempts = $1 WHERE id = $2`, newAttempts, challengeID)
		return ErrOTPInvalid
	}

	// Mark consumed
	_, err = r.db.ExecContext(ctx, `UPDATE phone_otp_challenges SET consumed_at = NOW() WHERE id = $1`, challengeID)
	if err != nil {
		return fmt.Errorf("failed to consume OTP challenge: %w", err)
	}

	return nil
}

// FindUserByPhone finds an existing user linked to this phone number
func (r *PhoneOTPRepo) FindUserByPhone(ctx context.Context, phoneInput string) (*ent.User, error) {
	phoneE164 := sms.NormalizeE164(phoneInput)
	nationalDigits := sms.CleanIndianPhone(phoneE164)
	syntheticEmail := fmt.Sprintf("phone_%s@fooddiscovery.com", nationalDigits)

	// 1. Check users.phone direct column
	var userID int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE phone = $1 OR phone = $2 LIMIT 1`, phoneE164, nationalDigits).Scan(&userID)
	if err == nil && userID > 0 {
		return r.orm.User.Get(ctx, userID)
	}

	// 2. Check profile.phone_number_e164
	p, err := r.orm.Profile.Query().
		Where(profile.Or(
			profile.PhoneNumberE164EQ(phoneE164),
			profile.PhoneNumberE164EQ(nationalDigits),
		)).
		WithUser().
		First(ctx)
	if err == nil && p != nil && p.Edges.User != nil {
		return p.Edges.User, nil
	}

	// 3. Check synthetic phone email
	u, err := r.orm.User.Query().
		Where(user.EmailEQ(syntheticEmail)).
		First(ctx)
	if err == nil && u != nil {
		return u, nil
	}

	return nil, ErrUserNotFound
}

// LinkPhoneToUser links verified phone number to the user and profile records
func (r *PhoneOTPRepo) LinkPhoneToUser(ctx context.Context, userID int, phoneInput string) error {
	phoneE164 := sms.NormalizeE164(phoneInput)

	// Update users.phone
	_, _ = r.db.ExecContext(ctx, `UPDATE users SET phone = $1 WHERE id = $2`, phoneE164, userID)

	// Update profile
	u, err := r.orm.User.Query().Where(user.ID(userID)).WithProfile().Only(ctx)
	if err == nil && u != nil {
		if u.Edges.Profile != nil {
			_, _ = u.Edges.Profile.Update().
				SetPhoneNumberE164(phoneE164).
				SetCountryCode("+91").
				SetPhoneVerified(true).
				Save(ctx)
		} else {
			_, _ = r.orm.Profile.Create().
				SetUser(u).
				SetPhoneNumberE164(phoneE164).
				SetCountryCode("+91").
				SetPhoneVerified(true).
				Save(ctx)
		}
	}
	return nil
}

// MaskEmail masks email for privacy preserving account recovery (e.g. aarav@example.com -> a***v@e***e.com)
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***.***"
	}
	name := parts[0]
	domainParts := strings.Split(parts[1], ".")

	var maskedName string
	if len(name) <= 2 {
		maskedName = name[:1] + "*"
	} else {
		maskedName = string(name[0]) + strings.Repeat("*", len(name)-2) + string(name[len(name)-1])
	}

	domainName := domainParts[0]
	var maskedDomain string
	if len(domainName) <= 2 {
		maskedDomain = domainName[:1] + "*"
	} else {
		maskedDomain = string(domainName[0]) + strings.Repeat("*", len(domainName)-2) + string(domainName[len(domainName)-1])
	}

	ext := "com"
	if len(domainParts) > 1 {
		ext = strings.Join(domainParts[1:], ".")
	}

	return fmt.Sprintf("%s@%s.%s", maskedName, maskedDomain, ext)
}
