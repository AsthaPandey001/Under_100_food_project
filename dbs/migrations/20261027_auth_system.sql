-- 20261027_auth_system.sql
-- Complete, Secure, Zero-Cost-First Authentication Schema Migration

-- 1. Add phone column to users table if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(32);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);

-- 2. Email verification tokens (stores only token hash)
CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user ON email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_hash ON email_verification_tokens(token_hash);

-- 3. Phone OTP challenges (stores only OTP hash, enforces max attempts, expiration, and cooldown)
CREATE TABLE IF NOT EXISTS phone_otp_challenges (
    id BIGSERIAL PRIMARY KEY,
    phone_e164 VARCHAR(32) NOT NULL,
    portal VARCHAR(32) NOT NULL DEFAULT 'customer',
    otp_hash VARCHAR(255) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(64)
);
CREATE INDEX IF NOT EXISTS idx_phone_otp_challenges_lookup ON phone_otp_challenges(phone_e164, portal, consumed_at, expires_at);
