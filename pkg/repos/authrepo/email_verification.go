package authrepo

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/user"
)

var (
	ErrTokenExpired       = errors.New("email verification link has expired")
	ErrTokenInvalid       = errors.New("email verification link is invalid")
	ErrTokenAlreadyUsed   = errors.New("email verification link has already been used")
	ErrResendCooldown     = errors.New("please wait 60 seconds before requesting another verification email")
	ErrUserAlreadyVerified = errors.New("account is already verified")
	ErrUserNotFound       = errors.New("account not found")
)

// EmailVerificationRepo manages secure single-use email verification tokens
type EmailVerificationRepo struct {
	db  *sql.DB
	orm *ent.Client
}

// NewEmailVerificationRepo creates a new EmailVerificationRepo
func NewEmailVerificationRepo(db *sql.DB, orm *ent.Client) *EmailVerificationRepo {
	return &EmailVerificationRepo{
		db:  db,
		orm: orm,
	}
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreateToken creates a cryptographically secure 32-byte hex token, stores its SHA-256 hash, and enforces a 60s cooldown
func (r *EmailVerificationRepo) CreateToken(ctx context.Context, userID int) (string, error) {
	// 1. Enforce 60-second cooldown on consecutive token creation
	var recentCount int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM email_verification_tokens 
		WHERE user_id = $1 AND created_at > NOW() - INTERVAL '60 seconds'
	`, userID).Scan(&recentCount)
	if err == nil && recentCount > 0 {
		return "", ErrResendCooldown
	}

	// 2. Invalidate any previous unconsumed tokens for this user
	_, _ = r.db.ExecContext(ctx, `
		UPDATE email_verification_tokens 
		SET consumed_at = NOW() 
		WHERE user_id = $1 AND consumed_at IS NULL
	`, userID)

	// 3. Generate 32-byte cryptographically secure random token (64 hex characters)
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate secure random token: %w", err)
	}
	rawToken := hex.EncodeToString(buf)
	tokenHash := HashToken(rawToken)

	// 4. Token expiration: 30 minutes
	expiresAt := time.Now().Add(30 * time.Minute)

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO email_verification_tokens (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
	`, userID, tokenHash, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to persist email verification token: %w", err)
	}

	return rawToken, nil
}

// VerifyToken validates the provided raw token against stored SHA-256 hash, marks it consumed, and verifies the user
func (r *EmailVerificationRepo) VerifyToken(ctx context.Context, rawToken string) (int, error) {
	cleanToken := strings.TrimSpace(rawToken)
	if cleanToken == "" {
		return 0, ErrTokenInvalid
	}

	tokenHash := HashToken(cleanToken)

	var tokenID int64
	var userID int
	var expiresAt time.Time
	var consumedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, expires_at, consumed_at 
		FROM email_verification_tokens 
		WHERE token_hash = $1
	`, tokenHash).Scan(&tokenID, &userID, &expiresAt, &consumedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTokenInvalid
		}
		return 0, fmt.Errorf("error querying verification token: %w", err)
	}

	if consumedAt.Valid {
		return 0, ErrTokenAlreadyUsed
	}

	if time.Now().After(expiresAt) {
		return 0, ErrTokenExpired
	}

	// Mark token consumed
	_, err = r.db.ExecContext(ctx, `
		UPDATE email_verification_tokens 
		SET consumed_at = NOW() 
		WHERE id = $1
	`, tokenID)
	if err != nil {
		return 0, fmt.Errorf("failed to consume verification token: %w", err)
	}

	// Mark user verified
	_, err = r.orm.User.UpdateOneID(userID).
		SetVerified(true).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to mark user verified: %w", err)
	}

	return userID, nil
}

// PrepareResend validates that user exists and is not verified, then generates a fresh verification token
func (r *EmailVerificationRepo) PrepareResend(ctx context.Context, email string) (string, *ent.User, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if cleanEmail == "" {
		return "", nil, ErrUserNotFound
	}

	u, err := r.orm.User.Query().
		Where(user.EmailEQ(cleanEmail)).
		Only(ctx)
	if err != nil {
		return "", nil, ErrUserNotFound
	}

	if u.Verified {
		return "", nil, ErrUserAlreadyVerified
	}

	rawToken, err := r.CreateToken(ctx, u.ID)
	if err != nil {
		return "", nil, err
	}

	return rawToken, u, nil
}
