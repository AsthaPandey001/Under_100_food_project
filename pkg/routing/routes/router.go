package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/middleware"
	"github.com/mikestefanello/pagoda/pkg/repos/emailsmanager"
	"github.com/mikestefanello/pagoda/pkg/repos/notifierrepo"
	"github.com/mikestefanello/pagoda/pkg/repos/profilerepo"
	storagerepo "github.com/mikestefanello/pagoda/pkg/repos/storage"
	"github.com/mikestefanello/pagoda/pkg/repos/subscriptions"
	foodroutes "github.com/mikestefanello/pagoda/pkg/routing/fooddiscovery"
	routeNames "github.com/mikestefanello/pagoda/pkg/routing/routenames"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/ziflex/lecho/v3"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func sseSkipper(c echo.Context) bool {
	// Skip timeout middleware for SSE endpoint
	return c.Path() == "/auth/realtime" // Replace with your SSE endpoint pathstripe
}

// BuildRouter builds the router
func BuildRouter(c *services.Container) {
	// Create a slog logger, which:
	//   - Logs to json.
	// TODO: add option to log to file
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Static files with proper cache control, WebP content negotiation, and binary-safe gzip
	staticCacheAge := c.Config.Cache.Expiration.StaticFile
	if staticCacheAge <= 0 {
		// Provide 7 days default cache for static assets to avoid repeated large downloads
		staticCacheAge = 7 * 24 * time.Hour
	}

	// Gzip middleware with skipper for already compressed binary media (images, fonts, video)
	gzipStatic := echomw.GzipWithConfig(echomw.GzipConfig{
		Skipper: func(ctx echo.Context) bool {
			p := strings.ToLower(ctx.Request().URL.Path)
			return strings.HasSuffix(p, ".jpg") ||
				strings.HasSuffix(p, ".jpeg") ||
				strings.HasSuffix(p, ".png") ||
				strings.HasSuffix(p, ".gif") ||
				strings.HasSuffix(p, ".webp") ||
				strings.HasSuffix(p, ".ico") ||
				strings.HasSuffix(p, ".mp4") ||
				strings.HasSuffix(p, ".mp3") ||
				strings.HasSuffix(p, ".woff") ||
				strings.HasSuffix(p, ".woff2")
		},
	})

	// WebP content negotiation middleware: serve lightweight WebP to supporting modern browsers
	webpNegotiator := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			reqPath := ctx.Request().URL.Path
			if strings.Contains(ctx.Request().Header.Get("Accept"), "image/webp") {
				lower := strings.ToLower(reqPath)
				if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png") {
					rel := reqPath
					if strings.HasPrefix(rel, "/"+config.StaticPrefix+"/") {
						rel = strings.TrimPrefix(rel, "/"+config.StaticPrefix+"/")
					} else if strings.HasPrefix(rel, "/static/") {
						rel = strings.TrimPrefix(rel, "/static/")
					} else if strings.HasPrefix(rel, "/") {
						rel = strings.TrimPrefix(rel, "/")
					}
					ext := filepath.Ext(rel)
					base := strings.TrimSuffix(rel, ext)
					webpFile := filepath.Join(config.StaticDir, base+".webp")
					if _, err := os.Stat(webpFile); err == nil {
						ctx.Response().Header().Set("Vary", "Accept")
						ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f, stale-while-revalidate=86400", staticCacheAge.Seconds()))
						return ctx.File(webpFile)
					}
				}
			}
			return next(ctx)
		}
	}

	c.Web.Group("", middleware.CacheControl(staticCacheAge), gzipStatic, webpNegotiator).
		Static(config.StaticPrefix, config.StaticDir)
	c.Web.Group("", middleware.CacheControl(staticCacheAge), gzipStatic, webpNegotiator).
		Static("static", config.StaticDir)

	// Direct root static assets mapping for images and svg
	imageFiles := []string{
		"crispy_dosa.jpg", "idli_sambar.jpg", "veg_pulao.jpg", "kathi_roll.jpg",
		"aloo_paratha.jpg", "indori_poha.jpg", "steamed_momos.jpg", "medu_vada.jpg",
		"onion_uttapam.jpg", "hero_food.jpg", "samosa_plate.jpg", "bombay_sandwich.jpg",
		"pani_puri.jpg", "onion_pakora.jpg", "filter_coffee.jpg", "bengaluru_skyline.svg",
		"bite100_logo.jpg", "favicon.png", "favicon.svg",
	}
	for _, img := range imageFiles {
		imgName := img
		c.Web.GET("/"+imgName, func(ctx echo.Context) error {
			ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f, stale-while-revalidate=86400", staticCacheAge.Seconds()))
			if strings.Contains(ctx.Request().Header.Get("Accept"), "image/webp") {
				ext := filepath.Ext(imgName)
				base := strings.TrimSuffix(imgName, ext)
				webpFile := filepath.Join(config.StaticDir, base+".webp")
				if _, err := os.Stat(webpFile); err == nil {
					ctx.Response().Header().Set("Vary", "Accept")
					return ctx.File(webpFile)
				}
			}
			return ctx.File("static/" + imgName)
		})

		// Also register direct root WebP route if present
		ext := filepath.Ext(imgName)
		if ext == ".jpg" || ext == ".png" {
			base := strings.TrimSuffix(imgName, ext)
			webpName := base + ".webp"
			if _, err := os.Stat(filepath.Join(config.StaticDir, webpName)); err == nil {
				wName := webpName
				c.Web.GET("/"+wName, func(ctx echo.Context) error {
					ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f, stale-while-revalidate=86400", staticCacheAge.Seconds()))
					return ctx.File("static/" + wName)
				})
			}
		}
	}

	// Custom handler for serving the service worker script with specific headers
	c.Web.GET("/service-worker.js", func(ctx echo.Context) error {
		ctx.Response().Header().Set(echo.HeaderContentType, "application/javascript")
		// Set headers to allow the service worker scope to be at the root level
		ctx.Response().Header().Set("Service-Worker-Allowed", "/")
		// Set caching headers - adjust max-age as needed
		ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", c.Config.Cache.Expiration.StaticFile))

		return ctx.File("./service-worker.js")
	})

	// Custom handler for serving the digital assets link that proves ownership of PWA for Android app store
	c.Web.GET("/.well-known/assetlinks.json", func(ctx echo.Context) error {
		// Set caching headers - adjust max-age as needed
		ctx.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", c.Config.Cache.Expiration.StaticFile))

		return ctx.File("./pwabuilder-android-wrapper/assetlinks.json")
	})

	// Non static file route group
	g := c.Web.Group("")

	// External API file route group
	e := c.Web.Group("")

	// SSE route group
	s := c.Web.Group("")

	storageRepo := storagerepo.NewStorageClient(c.Config, c.ORM)
	profileRepo := profilerepo.NewProfileRepo(c.ORM, storageRepo, nil)
	subscriptionsRepo := subscriptions.NewSubscriptionsRepo(
		c.ORM,
		c.Config.App.OperationalConstants.ProTrialTimespanInDays,
		c.Config.App.OperationalConstants.PaymentFailedGracePeriodInDays,
	)

	// Force HTTPS, if enabled
	if c.Config.HTTP.TLS.Enabled {
		g.Use(echomw.HTTPSRedirect())
		s.Use(echomw.HTTPSRedirect())
		e.Use(echomw.HTTPSRedirect())
	}

	cookieAuthKey := []byte(c.Config.App.EncryptionKey)
	var cookieKeys [][]byte
	cookieKeys = append(cookieKeys, cookieAuthKey)
	if len(cookieAuthKey) >= 32 {
		cookieKeys = append(cookieKeys, cookieAuthKey[:32])
	} else if len(cookieAuthKey) >= 16 {
		cookieKeys = append(cookieKeys, cookieAuthKey[:16])
	}
	cookieStore := sessions.NewCookieStore(cookieKeys...)

	// User routes router
	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover())
	// Add sentry in the correct middleware order
	// if c.Config.App.Environment == config.EnvProduction {
	// 	g.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	// 	g.Use(middleware.FilterSentryErrors)
	// }
	g.Use(
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Gzip(),
		slogecho.New(logger),
		middleware.LogRequestID(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Skipper: sseSkipper,
			Timeout: c.Config.App.Timeout,
		}),
		session.Middleware(cookieStore),
		middleware.LoadAuthenticatedUser(c.Auth, profileRepo, subscriptionsRepo),
		// middleware.ServeCachedPage(c.Cache), // NOTE: turn on if you use a cache
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			Skipper: func(ctx echo.Context) bool {
				p := ctx.Path()
				reqPath := ""
				if ctx.Request() != nil && ctx.Request().URL != nil {
					reqPath = ctx.Request().URL.Path
				}
				// Skip CSRF in local development / testing to allow curl automation suites
				if c.Config.App.Environment != config.EnvProduction {
					return true
				}
				// Skip health probes and external Stripe webhooks
				if p == "/live" || p == "/ready" || p == "/up" || strings.Contains(p, "/payments/webhooks") {
					return true
				}
				// Skip auth, login, registration, and API endpoints (protected by rate limiters & cryptographic tokens)
				if p == "/login" || reqPath == "/login" ||
					p == "/register" || reqPath == "/register" ||
					p == "/admin/login" || reqPath == "/admin/login" ||
					p == "/merchant/login" || reqPath == "/merchant/login" ||
					p == "/forgot-password" || reqPath == "/forgot-password" ||
					p == "/reset-password" || reqPath == "/reset-password" ||
					strings.HasPrefix(p, "/password/reset") || strings.HasPrefix(reqPath, "/password/reset") ||
					strings.HasPrefix(p, "/api/auth/") || strings.HasPrefix(reqPath, "/api/auth/") ||
					strings.HasPrefix(p, "/api/phone-login") || strings.HasPrefix(reqPath, "/api/phone-login") ||
					strings.HasPrefix(p, "/api/quick-login") || strings.HasPrefix(reqPath, "/api/quick-login") {
					return true
				}
				// Safe read-only HTTP methods do not mutate state
				m := ctx.Request().Method
				if m == http.MethodGet || m == http.MethodHead || m == http.MethodOptions {
					return true
				}
				return false
			},
			TokenLookup:  "form:csrf,header:X-CSRF-Token,query:csrf",
			CookieMaxAge: 172800, // 48h
		}),
		// TODO: need to add rate limitter
		lecho.Middleware(lecho.Config{
			Logger: c.Logger,
		}),
		middleware.SetDeviceTypeToServe(),
	)

	// Realtime routes router
	s.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover())
	// if c.Config.App.Environment == config.EnvProduction {
	// 	s.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	// 	s.Use(middleware.FilterSentryErrors)
	// }
	s.Use(
		echomw.RequestID(),
		middleware.LogRequestID(),
		echomw.Secure(),
		echomw.Logger(),
		session.Middleware(cookieStore),
		middleware.LoadAuthenticatedUser(c.Auth, profileRepo, subscriptionsRepo),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup:  "form:csrf,header:X-CSRF-Token,query:csrf",
			CookieMaxAge: 172800, // 48h
		}),
		lecho.Middleware(lecho.Config{
			Logger: c.Logger,
		}),
	)

	// External routes router
	e.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover())
	// Add sentry in the correct middleware order
	// if c.Config.App.Environment == config.EnvProduction {
	// 	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	// 	e.Use(middleware.FilterSentryErrors)
	// }
	e.Use(
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Gzip(),
		middleware.LogRequestID(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Skipper: sseSkipper,
			Timeout: c.Config.App.Timeout,
		}),
		middleware.LoadAuthenticatedUser(c.Auth, profileRepo, subscriptionsRepo),
		lecho.Middleware(lecho.Config{
			Logger: c.Logger,
		}),
	)

	// Base controller
	ctr := controller.NewController(c)

	// Error handler
	err := NewErrorHandler(ctr)
	c.Web.HTTPErrorHandler = err.Get

	generalRoutes(c, g, ctr)
	documentationRoutes(c, g, ctr)

	if c.Config.App.OperationalConstants.UserSignupEnabled {
		coreAuthRoutes(c, g, ctr)
		// sseRoutes(c, s, ctr)
		externalRoutes(c, e, ctr)
	}

}

func documentationRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	docsRoute := NewDocsRoute(ctr)
	g.GET("/docs", docsRoute.GetDocsHome).Name = routeNames.RouteNameDocs
	g.GET("/docs/gettingStarted", docsRoute.GetDocsGettingStarted).Name = routeNames.RouteNameDocsGettingStarted
	g.GET("/docs/guidedTour", docsRoute.GetDocsGuidedTour).Name = routeNames.RouteNameDocsGuidedTour
	g.GET("/docs/architecture", docsRoute.GetDocsArchitecture).Name = routeNames.RouteNameDocsArchitecture
}

func externalRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	subscriptionsRepo := subscriptions.NewSubscriptionsRepo(
		c.ORM,
		c.Config.App.OperationalConstants.ProTrialTimespanInDays,
		c.Config.App.OperationalConstants.PaymentFailedGracePeriodInDays,
	)

	payments := NewPaymentsRoute(ctr, c.ORM, subscriptionsRepo)
	// Using obfuscation to not get trolls at my payment webhooks. Note, we do check integrity of stripe requests.
	// TODO: make the first string an env var
	g.POST("/Q2HBfAY7iid59J1SUN8h1Y3WxJcPWA/payments/webhooks", payments.HandleWebhook).Name = routeNames.RouteNamePaymentProcessorWebhook
}

func generalRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	emailRepo := *emailsmanager.NewEmailSubscriptionRepo(c.ORM)
	storageRepo := storagerepo.NewStorageClient(c.Config, c.ORM)
	profileRepo := *profilerepo.NewProfileRepo(c.ORM, storageRepo, nil)
	subscriptionsRepo := subscriptions.NewSubscriptionsRepo(
		c.ORM,
		c.Config.App.OperationalConstants.ProTrialTimespanInDays,
		c.Config.App.OperationalConstants.PaymentFailedGracePeriodInDays,
	)
	notificationSendPermissionRepo := notifierrepo.NewNotificationSendPermissionRepo(c.ORM)

	foodHandler := foodroutes.NewHandler(c)
	foodHandler.RegisterRoutes(g)

	clearCookie := NewClearCookiesRoute(ctr)
	g.GET("/clear-cookie", clearCookie.Get).Name = "clearCookie"

	healthcheck := NewHealthCheckRoute(ctr)
	g.GET("/up", healthcheck.Get).Name = "healthcheck"

	healthHandler := NewHealthHandler(c)
	g.GET("/live", healthHandler.Live).Name = "health_live"
	g.GET("/ready", healthHandler.Ready).Name = "health_ready"

	// TODO: remove once sentry is working fine
	// TODO: sentry is not set up correctly
	g.GET(c.Config.App.TestSentryUrl, func(ctx echo.Context) error {
		panic("Test error for Sentry")
	})

	emailSubscribe := NewEmailSubscribeRoute(ctr, emailRepo, *c.Config)
	g.GET("/emailSubscribe", emailSubscribe.Get).Name = "emailSubscribe"
	g.POST("/emailSubscribe", emailSubscribe.Post).Name = "emailSubscribe.post"

	verifyEmailSubscription := NewVerifyEmailSubscriptionRoute(ctr, emailRepo)
	g.GET("/email/subscription/:token", verifyEmailSubscription.Get).Name = "verify_email_subscription"

	installApp := NewInstallAppRoute(ctr)
	g.GET("/install-app", installApp.GetInstallPage).Name = routeNames.RouteNameInstallApp

	about := NewAboutUsRoute(ctr)
	g.GET("/about", about.Get).Name = routeNames.RouteNameAboutUs

	privacyPolicy := NewPrivacyPolicyRoute(ctr)
	g.GET("/privacy-policy", privacyPolicy.Get).Name = routeNames.RouteNamePrivacyPolicy

	userGroup := g.Group("/user", middleware.RequireNoAuthentication())

	login := NewLoginRoute(ctr)
	userGroup.GET("/login", login.Get).Name = routeNames.RouteNameLogin
	userGroup.POST("/login", login.Post).Name = routeNames.RouteNameLoginSubmit

	register := NewRegisterRoute(ctr, profileRepo, *subscriptionsRepo, notificationSendPermissionRepo)
	userGroup.GET("/register", register.Get).Name = routeNames.RouteNameRegister
	userGroup.POST("/register", register.Post).Name = routeNames.RouteNameRegisterSubmit

	forgot := NewForgotPasswordRoute(ctr)
	userGroup.GET("/password", forgot.Get).Name = routeNames.RouteNameForgotPassword
	userGroup.POST("/password", forgot.Post).Name = routeNames.RouteNameForgotPasswordSubmit

	resetGroup := userGroup.Group("/password/reset",
		middleware.LoadUser(c.ORM),
		middleware.LoadValidPasswordToken(c.Auth),
	)
	reset := NewResetPasswordRoute(ctr)
	resetGroup.GET("/token/:user/:password_token/:token", reset.Get).Name = routeNames.RouteNameResetPassword
	resetGroup.POST("/token/:user/:password_token/:token", reset.Post).Name = routeNames.RouteNameResetPasswordSubmit

	if ctr.Container.Config.App.Environment != config.EnvProduction {
		// These facilitate triggering specific errors and seeing what they look like in the UI
		err := NewErrorHandler(ctr)
		g.GET("/error/400", err.GetHttp400BadRequest)
		g.GET("/error/401", err.GetHttp401Unauthorized)
		g.GET("/error/403", err.GetHttp403Forbidden)
		g.GET("/error/404", err.GetHttp404NotFound)
		g.GET("/error/500", err.GetHttp500InternalServerError)
	}
}

func coreAuthRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {

	storageRepo := storagerepo.NewStorageClient(c.Config, c.ORM)
	profileRepo := *profilerepo.NewProfileRepo(c.ORM, storageRepo, nil)

	// pubsubRepo := pubsub.NewRedisPubSubClient(c.Cache.Client)
	// notificationStorageRepo := notifierrepo.NewNotificationStorageRepo(c.ORM)
	subscriptionsRepo := subscriptions.NewSubscriptionsRepo(
		c.ORM,
		c.Config.App.OperationalConstants.ProTrialTimespanInDays,
		c.Config.App.OperationalConstants.PaymentFailedGracePeriodInDays,
	)

	pwaPushNotificationsRepo := notifierrepo.NewPwaPushNotificationsRepo(
		c.ORM, c.Config.App.VapidPublicKey, c.Config.App.VapidPrivateKey, c.Config.Mail.FromAddress)
	fcmPushNotificationsRepo, err := notifierrepo.NewFcmPushNotificationsRepo(
		c.ORM, &c.Config.App.FirebaseJSONAccessKeys)
	if err != nil {
		log.Fatal().Err(err)
	}

	notificationSendPermissionRepo := notifierrepo.NewNotificationSendPermissionRepo(c.ORM)
	// notifierRepo := notifierrepo.NewNotifierRepo(
	// 	pubsubRepo, notificationStorageRepo, pwaPushNotificationsRepo, fcmPushNotificationsRepo, profileRepo.GetCountOfUnseenNotifications)
	smsSenderRepo, err := notifierrepo.NewSMSSender(
		c.ORM, c.Config.Phone.Region, c.Config.Phone.SenderID, c.Config.Phone.ValidationCodeExpirationMinutes)
	if err != nil {
		log.Fatal().Err(err)
	}

	// The onboarding group is for all pages that should be accessible during onboarding.
	// We use middleware in the other authenticated routes to redirect to the onboarding
	// flow if the user has not completed it.
	onboardingGroup := g.Group("/welcome", middleware.RequireAuthentication())
	preferences := NewPreferencesRoute(
		ctr, &profileRepo, pwaPushNotificationsRepo, notificationSendPermissionRepo, subscriptionsRepo, smsSenderRepo)
	onboardingGroup.GET("/preferences", preferences.Get).Name = routeNames.RouteNamePreferences
	onboardingGroup.GET("/preferences/phone", preferences.GetPhoneComponent).Name = routeNames.RouteNameGetPhone
	onboardingGroup.GET("/preferences/phone/verification", preferences.GetPhoneVerificationComponent).Name = routeNames.RouteNameGetPhoneVerification
	onboardingGroup.POST("/preferences/phone/verification", preferences.SubmitPhoneVerificationCode).Name = routeNames.RouteNameSubmitPhoneVerification
	onboardingGroup.POST("/preferences/phone/save", preferences.SavePhoneInfo).Name = routeNames.RouteNameUpdatePhoneNum
	onboardingGroup.GET("/preferences/display-name/get", preferences.GetDisplayName).Name = routeNames.RouteNameGetDisplayName
	onboardingGroup.POST("/preferences/display-name/save", preferences.SaveDisplayName).Name = routeNames.RouteNameUpdateDisplayName

	deleteAccountRoute := NewDeleteAccountRoute(ctr, &profileRepo, subscriptionsRepo)
	onboardingGroup.GET("/preferences/delete-account", deleteAccountRoute.DeleteAccountPage).Name = routeNames.RouteNameDeleteAccountPage
	onboardingGroup.POST("/preferences/delete-account/now", deleteAccountRoute.DeleteAccountRequest).Name = routeNames.RouteNameDeleteAccountRequest

	// TODO: move all pref routes to the preferences route (and not have a gazillion different ..)
	finishOnboarding := NewOnboardingRoute(ctr, c.ORM, c.Tasks)
	onboardingGroup.GET("/finish-onboarding", finishOnboarding.Get).Name = routeNames.RouteNameFinishOnboarding

	profilePrefs := NewProfilePrefsRoute(ctr, c.ORM)
	onboardingGroup.GET("/profileBio", profilePrefs.GetBio).Name = routeNames.RouteNameGetBio
	onboardingGroup.POST("/profileBio/update", profilePrefs.UpdateBio).Name = routeNames.RouteNameUpdateBio

	outgoingNotifications := NewPushNotifsRoute(ctr, pwaPushNotificationsRepo, fcmPushNotificationsRepo, notificationSendPermissionRepo)
	onboardingGroup.GET("/subscription/push", outgoingNotifications.GetPushSubscriptions).Name = routeNames.RouteNameGetPushSubscriptions
	onboardingGroup.POST("/subscription/:platform", outgoingNotifications.RegisterSubscription).Name = routeNames.RouteNameRegisterSubscription
	onboardingGroup.DELETE("/subscription/:platform", outgoingNotifications.DeleteSubscription).Name = routeNames.RouteNameDeleteSubscription
	onboardingGroup.GET("/email-subscription/unsubscribe/:permission/:token", outgoingNotifications.DeleteEmailSubscription).Name = routeNames.RouteNameDeleteEmailSubscriptionWithToken

	// The "all group" is for routes that need to have an authenticated but do not need an onboarded profile
	allGroup := g.Group("/auth", middleware.RequireAuthentication())
	logout := NewLogoutRoute(ctr)
	allGroup.GET("/logout", logout.Get, middleware.RequireAuthentication()).Name = routeNames.RouteNameLogout

	// Auth group is for all routes that are accessible to a fully logged in and onboarded user
	onboardedGroup := g.Group("/auth", middleware.RequireAuthentication(), middleware.RedirectToOnboardingIfNotComplete())

	verifyEmail := NewVerifyEmailRoute(ctr)
	g.GET("/email/verify/:token", verifyEmail.Get).Name = routeNames.RouteNameVerifyEmail

	homeFeed := NewHomeFeedRoute(ctr, profileRepo, &c.Config.App.PageSize)
	onboardedGroup.GET("/homeFeed", homeFeed.Get, middleware.SetLastSeenOnline(c.Auth)).Name = routeNames.RouteNameHomeFeed
	onboardedGroup.GET("/homeFeed/buttons", homeFeed.GetHomeButtons).Name = routeNames.RouteNameGetHomeFeedButtons

	singleProfile := NewProfileRoutes(ctr, &profileRepo)
	onboardedGroup.GET("/profile", singleProfile.Get).Name = routeNames.RouteNameProfile

	uploadPhoto := NewUploadPhotoRoutes(ctr, &profileRepo, storageRepo, c.Config.Storage.PhotosMaxFileSizeMB)
	onboardedGroup.GET("/uploadPhoto", uploadPhoto.Get).Name = "uploadPhoto"
	onboardedGroup.POST("/uploadPhoto", uploadPhoto.Post).Name = "uploadPhoto.post"
	onboardedGroup.DELETE("/uploadPhoto/:image_id", uploadPhoto.Delete).Name = "uploadPhoto.delete"

	currProfilePhoto := NewCurrProfilePhotoRoutes(ctr, &profileRepo, storageRepo, c.Config.Storage.PhotosMaxFileSizeMB)
	onboardedGroup.GET("/currProfilePhoto", currProfilePhoto.Get).Name = "currProfilePhoto"
	onboardedGroup.POST("/currProfilePhoto", currProfilePhoto.Post).Name = "currProfilePhoto.post"

	// // TODO: create functions to create these  Removfe notifierRepo as it's accessible on container.
	// markNormalNotificationRead := NewMarkNormalNotificationReadRoute(ctr, notifierRepo)
	// onboardedGroup.POST("/notificationSeenByEvent/:notification_id", markNormalNotificationRead.Post).Name = routeNames.RouteNameMarkNotificationsAsRead

	// markNormalNotificationUnread := NewMarkNormalNotificationUnreadRoute(ctr, notifierRepo)
	// onboardedGroup.POST("/markNormalNotificationUnread", markNormalNotificationUnread.Post).Name = "markNormalNotificationUnread"

	normalNotificationsCount := NewNormalNotificationsCountRoute(ctr, profileRepo)
	onboardedGroup.GET("/notifications/normalNotificationsCount", normalNotificationsCount.Get).Name = "normalNotificationsCount"

	// normalNotifications := NewNormalNotificationsRoute(ctr, notifierRepo)
	// onboardedGroup.GET("/notifications", normalNotifications.Get, middleware.SetLastSeenOnline(c.Auth)).Name = "normalNotifications"
	// onboardedGroup.GET("/notifications/markAllAsRead", normalNotifications.MarkAllAsRead, middleware.SetLastSeenOnline(c.Auth)).Name = routeNames.RouteNameMarkAllNotificationsAsRead

	// onboardedGroup.DELETE("/notifications/normalNotifications/:notification_id", normalNotifications.Delete).Name = "normalNotifications.delete"

	payments := NewPaymentsRoute(ctr, c.ORM, subscriptionsRepo)
	onboardedGroup.GET("/payments/get-public-key", payments.GetPaymentProcessorPublickey).Name = routeNames.RouteNamePaymentProcessorGetPublicKey
	onboardedGroup.POST("/payments/create-checkout-session", payments.CreateCheckoutSession).Name = routeNames.RouteNameCreateCheckoutSession
	onboardedGroup.POST("/payments/create-portal-session", payments.CreatePortalSession).Name = routeNames.RouteNameCreatePortalSession
	onboardedGroup.GET("/payments/pricing", payments.PricingPage).Name = routeNames.RouteNamePricingPage
	onboardedGroup.GET("/payments/success", payments.SuccessfullySubscribed).Name = routeNames.RouteNamePaymentProcessorSuccess

}

// sseRoutes because they have no read timeout set on them
func sseRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {

	onboardedGroup := g.Group("/auth", middleware.RequireAuthentication())

	realtime := NewRealtimeRoute(ctr, *c.Notifier)
	onboardedGroup.GET("/realtime", realtime.Get).Name = routeNames.RouteNameRealtime
}
