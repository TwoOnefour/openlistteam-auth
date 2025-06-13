package auth

import (
	"fmt"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	frontEndBaseUrl string
)

func initVar() {
	s := os.Getenv("ALI_LIMIT_MINUTES")
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	aliMinutes = v
	s = os.Getenv("ALI_LIMIT_MAX")
	v, err = strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	aliMax = v
	// client
	aliClientID = os.Getenv("ALI_DRIVE_CLIENT_ID")
	aliClientSecret = os.Getenv("ALI_DRIVE_CLIENT_SECRET")

	baiduClientId = os.Getenv("BAIDU_CLIENT_ID")
	baiduClientSecret = os.Getenv("BAIDU_CLIENT_SECRET")
	baiduCallbackUri = frontEndBaseUrl + "/tool/baidu/callback"

	frontEndBaseUrl = os.Getenv("API_BASE")
	if strings.TrimSpace(frontEndBaseUrl) == "" {
		panic(fmt.Errorf("ENV API_BASE is empty"))
	}
	oneDriveCallBackUri = frontEndBaseUrl + "/tool/onedrive/callback"

	dropBoxAppId = os.Getenv("DROPBOX_APP_ID")
	dropBoxAppSecret = os.Getenv("DROPBOX_APP_SECRET")
	if dropBoxAppId == "" || dropBoxAppSecret == "" {
		panic(fmt.Errorf("DROPBOX_APP_ID/APP_SECRET is empty"))
	}

}

func SetupOnedriveApi(g *gin.RouterGroup) {
	g.POST("/onedrive/get_refresh_token", onedriveToken)
	// CORS should be settle by caddy, nginx or other reverse proxy middleware, but for test purposes, I set it here
	g.OPTIONS("/onedrive/get_refresh_token", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Status(204)
	})
	g.POST("/onedrive/get_site_id", spSiteID)
}

func SetUpAliApi(g *gin.RouterGroup) {
	g.GET("/ali/qr", Qr)
	g.POST("/ali/ck", Ck)
	aliOpen := g.Group("/ali_open")
	aliOpen.Any("/limit", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"minutes": aliMinutes,
			"max":     aliMax,
		})
	})
	aliOpenLimit := aliOpen.Group("")
	aliOpenLimit.Use(ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
		LimitKey: func(ctx *gin.Context) string {
			return ctx.ClientIP()
		},
		LimitedHandler: func(ctx *gin.Context) {
			ctx.JSON(429, gin.H{
				"code":    "Too Many Requests",
				"message": "Too Many Requests",
				"error":   "Too Many Requests",
			})
			ctx.Abort()
		},
		TokenBucketConfig: func(context *gin.Context) (time.Duration, int) {
			return time.Duration(aliMinutes) * time.Minute, aliMax
		},
	}))
	aliOpenLimit.Any("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ClientIP": c.ClientIP(),
			"RemoteIP": c.RemoteIP(),
		})
	})
	aliOpenLimit.Any("/token", aliAccessToken)
	aliOpenLimit.Any("/refresh", aliAccessToken)
	aliOpenLimit.Any("/code", aliAccessToken)
	aliOpenLimit.Any("/qr", aliQrcode)
}

func SetupBaiduApi(g *gin.RouterGroup) {
	g.GET("/baidu/get_refresh_token", baiduToken)
}

func Setup115Api(g *gin.RouterGroup) {
	g.GET("/115/auth_device_code", Open115Qrcode)
	g.POST("/115/get_token", Open115Token)
}

func SetupDropboxApi(g *gin.RouterGroup) {
	g.POST("/dropbox/token", getDropBoxToken)
}

func Setup(g *gin.RouterGroup) {
	initVar()
	SetUpAliApi(g)
	SetupOnedriveApi(g)
	SetupBaiduApi(g)
	Setup115Api(g)
	SetupDropboxApi(g)
}
