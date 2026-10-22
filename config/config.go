package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// TemplateExt stores the extension used for the template files
	TemplateExt = ".gohtml"

	// StaticDir stores the name of the directory that will serve static files
	StaticDir = "static"

	// StaticPrefix stores the URL prefix used when serving static files
	StaticPrefix = "files"
)

type app string
type environment string
type dbmode string

const (
	// EnvLocal represents the local environment
	EnvLocal environment = "local"

	// EnvTest represents the test environment
	EnvTest environment = "test"

	// EnvDevelop represents the development environment
	EnvDevelop environment = "dev"

	// EnvStaging represents the staging environment
	EnvStaging environment = "staging"

	// EnvQA represents the qa environment
	EnvQA environment = "qa"

	// EnvProduction represents the production environment
	EnvProduction environment = "prod"

	// DBModeEmbedded represents an embedded DB being used as a storage backend
	DBModeEmbedded dbmode = "embedded"

	// DBModeStandalone represents a standalone DB as being used as a storage backend
	DBModeStandalone dbmode = "standalone"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env environment) {
	if err := os.Setenv("PAGODA_APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete configuration
	Config struct {
		HTTP        HTTPConfig
		App         AppConfig
		Cache       CacheConfig
		Database    DatabaseConfig
		Mail        MailConfig
		Phone       PhoneConfig
		Recommender RecommenderConfig
		Storage     StorageConfig
		RateLimit   RateLimitConfig
	}

	// HTTPConfig stores HTTP configuration
	HTTPConfig struct {
		Hostname     string
		Port         uint16
		Domain       string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
		SseKeepAlive time.Duration
		TLS          struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	// AppConfig stores application configuration
	AppConfig struct {
		Name          app
		SupportEmail  string
		Environment   environment
		EncryptionKey string
		Timeout       time.Duration
		PasswordToken struct {
			Expiration time.Duration
			Length     int
		}
		EmailVerificationTokenExpiration time.Duration
		OperationalConstants             OperationalConstants
		PageSize                         int
		VapidPublicKey                   string
		VapidPrivateKey                  string
		SentryDsn                        string
		TestSentryUrl                    string
		PublicStripeKey                  string
		PrivateStripeKey                 string
		StripeWebhookSecret              string
		AppEncryptionKey                 string
		FirebaseBase64AccessKeys         string
		FirebaseJSONAccessKeys           []byte
	}

	OperationalConstants struct {
		NewsletterSignupEnabled                           bool
		UserSignupEnabled                                 bool
		UserSignupEnabledOnLandingPage                    bool
		QuestionInteractionValidLifetimeInDays            int
		NumMinSharedAnswersForPrivateMessages             int
		NotifEmojiDebounceTime                            time.Duration
		NotifyNewAnswerFromUnansweredQuestionDebounceTime time.Duration
		MinAnswerLen                                      int
		PaymentsEnabled                                   bool
		ProTrialTimespanInDays                            int
		ProductProCode                                    string
		ProductProPrice                                   float32
		PaymentFailedGracePeriodInDays                    int
		DeleteStaleNotificationAfterDays                  int
		MaxLikedQuestionHistoryFreePlan                   int
	}

	// CacheConfig stores the cache configuration
	CacheConfig struct {
		Hostname     string
		Port         uint16
		Password     string
		Database     int
		TestDatabase int
		Expiration   struct {
			StaticFile time.Duration
			Page       time.Duration
		}
	}

	// DatabaseConfig stores the database configuration
	DatabaseConfig struct {
		DbMode                 dbmode
		EmbeddedDriver         string
		EmbeddedConnection     string
		EmbeddedTestConnection string
		// TODO: eventually separate in-memory (SQLite) and standalone DB configs
		Hostname          string
		Port              uint16
		User              string
		Password          string
		DatabaseNameLocal string
		DatabaseNameProd  string
		TestDatabase      string
		SslCertPath       string
		SslMode           string
	}

	// MailConfig stores the mail configuration
	MailConfig struct {
		Hostname     string
		HttpPort     uint16
		SmtpPort     uint16
		User         string
		Password     string
		FromAddress  string
		ResendAPIKey string
	}

	PhoneConfig struct {
		SenderID                        string
		Region                          string
		ValidationCodeExpirationMinutes int
	}

	RecommenderConfig struct {
		NumProfilesToMatchAtOnce int
	}

	StorageConfig struct {
		AppBucketName             string
		StaticFilesBucketName     string
		S3Endpoint                string
		S3AccessKey               string
		S3SecretKey               string
		S3UseSSL                  bool
		ProfilePhotoMaxFileSizeMB int64
		PhotosMaxFileSizeMB       int64
	}

	RateLimitConfig struct {
		Enabled             bool
		LoginLimit          int
		RegisterLimit       int
		PasswordForgotLimit int
		PasswordResetLimit  int
	}
)

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var c Config

	// Common config loading
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../../../config")

	// Load the config file
	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	// Check the environment variable PAGODA_APP_ENVIRONMENT
	env := os.Getenv("PAGODA_APP_ENVIRONMENT")
	if env == string(EnvProduction) || c.App.Environment == EnvProduction {
		// Load env variables for production
		viper.SetEnvPrefix("pagoda")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// Load the config file
	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	if c.App.FirebaseBase64AccessKeys != "" {
		jsonCreds, err := base64.StdEncoding.DecodeString(c.App.FirebaseBase64AccessKeys)
		if err != nil {
			return c, fmt.Errorf("error decoding firebase credentials: %w", err)
		}
		c.App.FirebaseJSONAccessKeys = jsonCreds
	}

	if c.RateLimit.LoginLimit <= 0 {
		c.RateLimit.LoginLimit = 60
	}
	if c.RateLimit.RegisterLimit <= 0 {
		c.RateLimit.RegisterLimit = 30
	}
	if c.RateLimit.PasswordForgotLimit <= 0 {
		c.RateLimit.PasswordForgotLimit = 30
	}
	if c.RateLimit.PasswordResetLimit <= 0 {
		c.RateLimit.PasswordResetLimit = 30
	}

	if err := c.Validate(); err != nil {
		return c, err
	}

	return c, nil
}

// Validate validates configuration settings and fails safely if critical production secrets are missing or weak
func (c *Config) Validate() error {
	if c.App.Environment == EnvProduction {
		if c.App.EncryptionKey == "" || c.App.EncryptionKey == "?E(G+KbPeShVmYq3t6w9z$C&F)J@McQf" {
			return fmt.Errorf("production security failure: PAGODA_APP_ENCRYPTIONKEY must be set to a secure unique key and cannot use the development default")
		}
		if len(c.App.EncryptionKey) < 16 {
			return fmt.Errorf("production security failure: PAGODA_APP_ENCRYPTIONKEY must be at least 16 characters")
		}
		if c.Database.DbMode == DBModeStandalone && c.Database.Password == "" {
			return fmt.Errorf("production security failure: standalone database password must be configured in production")
		}
		if c.HTTP.Domain == "" {
			return fmt.Errorf("production security failure: HTTP domain must be configured in production for trusted link generation")
		}
	}
	return nil
}
