package routers

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
	"net/http"
	"scopelens-server/config"
	"scopelens-server/middleware"
	"scopelens-server/utils/logger"
	"time"
)

func InitRouters() *gin.Engine {
	r := gin.Default()

	// static resources
	r.Static("/assets", "./assets")

	// middleware
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger.Logger, time.RFC3339, false))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))

	// CORS
	r.Use(middleware.Cors())
	if config.Mode == "release" {
		r.Use(middleware.TlsHandler())

		// limit access rate by custom key (here is IP) and rate for POST
		limiterMiddleware := limit.NewRateLimiter(func(c *gin.Context) string {
			return c.ClientIP() // limit rate by client ip
		}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
			// limit 1/60 qps/clientIp and permit bursts of at most 3 tokens,
			// and the limiter liveness time duration is 1 day
			// https://www.cyhone.com/articles/usage-of-golang-rate/
			return rate.NewLimiter(rate.Every(time.Minute), 3), 24 * time.Hour
		}, func(c *gin.Context) {
			if c.Request.Method == http.MethodPost {
				c.AbortWithStatus(429) // handle exceed rate limit request
			}
		})
		r.Use(func(c *gin.Context) {
			if c.Request.Method == "POST" {
				limiterMiddleware(c)
				return
			}
			c.Next()
		})
	}

	apiGroups := r.Group("api")
	InitUserRouter(apiGroups)
	InitAuthRouter(apiGroups)
	InitAuthAUTHRouter(apiGroups)
	InitFormatRouter(apiGroups)
	InitTeamRouter(apiGroups)
	InitTeamAUTHRouter(apiGroups)

	return r
}
