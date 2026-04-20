package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	authula "github.com/Authula/authula"
	authulaconfig "github.com/Authula/authula/config"
	authulaenv "github.com/Authula/authula/env"
	authulaevents "github.com/Authula/authula/events"
	authulamodels "github.com/Authula/authula/models"

	csrfplugin "github.com/Authula/authula/plugins/csrf"
	emailplugin "github.com/Authula/authula/plugins/email"
	emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
	emailpasswordplugintypes "github.com/Authula/authula/plugins/email-password/types"
	emailplugintypes "github.com/Authula/authula/plugins/email/types"

	// bearerplugin "github.com/Authula/authula/plugins/bearer"
	// jwtplugin "github.com/Authula/authula/plugins/jwt"
	// jwtplugintypes "github.com/Authula/authula/plugins/jwt/types"

	oauth2plugin "github.com/Authula/authula/plugins/oauth2"
	oauth2plugintypes "github.com/Authula/authula/plugins/oauth2/types"
	ratelimitplugin "github.com/Authula/authula/plugins/rate-limit"
	secondarystorageplugin "github.com/Authula/authula/plugins/secondary-storage"
	sessionplugin "github.com/Authula/authula/plugins/session"

	loggerplugin "github.com/Authula/authula-playground/plugins/logger"
	loggerplugintypes "github.com/Authula/authula-playground/plugins/logger/types"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// -------------------------------------
	// Init Authula config
	// -------------------------------------

	config := authulaconfig.NewConfig(
		authulaconfig.WithAppName("AuthulaPlayground"),
		authulaconfig.WithBasePath("/api/auth"),
		authulaconfig.WithDatabase(authulamodels.DatabaseConfig{
			Provider: "postgres",
			URL:      os.Getenv(authulaenv.EnvDatabaseURL),
		}),
		authulaconfig.WithLogger(authulamodels.LoggerConfig{
			Level: "debug",
		}),
		authulaconfig.WithSession(authulamodels.SessionConfig{
			CookieName:         "authula.session_token",
			ExpiresIn:          time.Hour,
			UpdateAge:          5 * time.Minute,
			CookieMaxAge:       24 * time.Hour,
			Secure:             false,
			HttpOnly:           true,
			SameSite:           "lax",
			MaxSessionsPerUser: 5,
			AutoCleanup:        true,
			CleanupInterval:    time.Minute,
		}),
		authulaconfig.WithVerification(authulamodels.VerificationConfig{
			AutoCleanup:     true,
			CleanupInterval: time.Minute,
		}),
		authulaconfig.WithSecurity(authulamodels.SecurityConfig{
			TrustedOrigins: []string{"http://localhost:3000"},
			CORS: authulamodels.CORSConfig{
				AllowCredentials: true,
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Authorization", "Content-Type", "Set-Cookie", "Cookie", "X-AUTHULA-CSRF-TOKEN"},
				ExposedHeaders:   []string{"X-AUTHULA-CSRF-TOKEN"},
				MaxAge:           24 * time.Hour,
			},
		}),
		authulaconfig.WithEventBus(authulamodels.EventBusConfig{
			Provider: authulaevents.ProviderKafka,
			Kafka: &authulamodels.KafkaConfig{
				Brokers:       os.Getenv(authulaenv.EnvKafkaBrokers),
				ConsumerGroup: os.Getenv(authulaenv.EnvEventBusConsumerGroup),
			},
		}),
		authulaconfig.WithRouteMappings([]authulamodels.RouteMapping{
			// Core Routes
			{
				Paths:   []string{"GET:/me"},
				Plugins: []string{sessionplugin.HookIDSessionAuth.String()},
			},
			{
				Paths: []string{"POST:/sign-out"},
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// Email-Password Routes
			{
				Paths: []string{
					"POST:/email-password/sign-in",
					"POST:/email-password/sign-up",
				},
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Paths:   []string{"GET:/email-password/verify-email"},
				Plugins: []string{sessionplugin.HookIDSessionAuthOptional.String()},
			},
			{
				Paths: []string{
					"POST:/email-password/send-email-verification",
					"POST:/email-password/request-password-reset",
					"POST:/email-password/change-password",
					"POST:/email-password/request-email-change",
				},
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// Custom Routes
			{
				Paths:   []string{"GET:/api/v1/health"},
				Plugins: []string{},
			},
		}),
	)

	// -------------------------------------
	// Init Authula instance
	// -------------------------------------

	authula := authula.New(&authula.AuthConfig{
		Config: config,
		Plugins: []authulamodels.Plugin{
			// Built-in plugins
			// Secondary storage plugin MUST be registered before rate-limit plugin
			// This allows rate-limit to optionally use Redis/database for distributed rate limiting
			secondarystorageplugin.New(secondarystorageplugin.SecondaryStoragePluginConfig{
				Enabled:  true,
				Provider: secondarystorageplugin.SecondaryStorageProviderRedis,
				Redis: &secondarystorageplugin.RedisStorageConfig{
					URL: os.Getenv(authulaenv.EnvRedisURL),
				},
			}),
			csrfplugin.New(csrfplugin.CSRFPluginConfig{
				Enabled: false,
			}),
			emailplugin.New(emailplugintypes.EmailPluginConfig{
				Enabled:     true,
				Provider:    emailplugintypes.ProviderSMTP,
				FromAddress: "noreply@example.com",
				TLSMode:     emailplugintypes.SMTPTLSModeOff,
			}),
			emailpasswordplugin.New(emailpasswordplugintypes.EmailPasswordPluginConfig{
				Enabled:                     true,
				MinPasswordLength:           8,
				MaxPasswordLength:           32,
				DisableSignUp:               false,
				RequireEmailVerification:    true,
				AutoSignIn:                  true,
				SendEmailOnSignUp:           true,
				SendEmailOnSignIn:           false,
				EmailVerificationExpiresIn:  24 * time.Hour,
				PasswordResetExpiresIn:      time.Hour,
				RequestEmailChangeExpiresIn: time.Hour,
			}),
			oauth2plugin.New(oauth2plugintypes.OAuth2PluginConfig{
				Enabled: true,
				Providers: map[string]oauth2plugintypes.ProviderConfig{
					"discord": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvDiscordClientID),
						ClientSecret: os.Getenv(authulaenv.EnvDiscordClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/discord", config.BaseURL, config.BasePath),
					},
					"github": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvGithubClientID),
						ClientSecret: os.Getenv(authulaenv.EnvGithubClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/github", config.BaseURL, config.BasePath),
					},
					"google": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvGoogleClientID),
						ClientSecret: os.Getenv(authulaenv.EnvGoogleClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/google", config.BaseURL, config.BasePath),
					},
				},
			}),
			sessionplugin.New(sessionplugin.SessionPluginConfig{
				Enabled: true,
			}),
			// jwtplugin.New(jwtplugintypes.JWTPluginConfig{
			// 	Enabled:   true,
			// 	Algorithm: jwtplugintypes.JWTAlgEdDSA,
			// }),
			// bearerplugin.New(bearerplugin.BearerPluginConfig{
			// 	Enabled: true,
			// }),
			ratelimitplugin.New(ratelimitplugin.RateLimitPluginConfig{
				Enabled:  true,
				Provider: ratelimitplugin.RateLimitProviderRedis,
			}),

			// Custom plugins
			loggerplugin.New(loggerplugintypes.LoggerPluginConfig{
				Enabled:     true,
				MaxLogCount: 10,
			}),
		},
	})

	// You can uncomment the following 2 lines to drop all migrations (i.e., reset the database).
	// ctx := context.Background()
	// if err := authula.PluginRegistry.DropMigrations(ctx); err != nil {
	// 	slog.Error("failed to drop plugin migrations", "error", err)
	// 	return
	// }
	// if err := authula.DropCoreMigrations(ctx); err != nil {
	// 	slog.Error("failed to drop core migrations", "error", err)
	// 	return
	// }
	// slog.Info("all migrations dropped successfully")
	// return

	// -------------------------------------
	// Add custom routes to the router
	// Note: Call RegisterCustomRoute() BEFORE Handler() to ensure routes are registered before handler is served
	// Custom routes are registered without the /api/auth prefix
	// -------------------------------------

	// Health check endpoint
	authula.RegisterCustomRoute(authulamodels.Route{
		Method: "GET",
		Path:   "/api/v1/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx, _ := authulamodels.GetRequestContext(r.Context())
			reqCtx.SetJSONResponse(http.StatusOK, map[string]any{
				"status": "ok",
			})
		}),
	})

	// authula.RegisterHook(authulamodels.Hook{
	// 	Stage: authulamodels.HookBefore,
	// 	Matcher: func(ctx *authulamodels.RequestContext) bool {
	// 		return ctx.UserID != nil && *ctx.UserID != "" && slices.Contains(
	// 			[]string{
	// 				"/api/protected",
	// 				"/path/to/more/routes...",
	// 			},
	// 			ctx.Path,
	// 		)
	// 	},
	// 	Handler: func(ctx *authulamodels.RequestContext) error {
	// 		// Do as you wish before the request is processed by the route handler...
	// 		return nil
	// 	},
	// })

	// -------------------------------------
	// Attach Authula handler to your chosen framework and run your server
	// All hooks (CORS, auth, rate limiting, etc.) are applied via the plugin system
	// -------------------------------------

	port := os.Getenv(authulaenv.EnvPort)
	slog.Debug(fmt.Sprintf("Server running on http://localhost:%s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), authula.Handler()); err != nil {
		slog.Error("Server error", "err", err)
	}
}
