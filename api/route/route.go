package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"muxi-thief/api/response"
	"muxi-thief/middleware"
	"time"
)

var ProviderSet = wire.NewSet(
	NewApp,
	NewRouter,
)

type App struct {
	r *gin.Engine
}

func NewApp(r *gin.Engine) App {
	return App{
		r: r,
	}
}

// 启动
func (a *App) Run() {
	a.r.Run("0.0.0.0:8000")
}

type ControllerProxy interface {
	Start(ctx *gin.Context)
	Login(ctx *gin.Context)
	Paper(ctx *gin.Context)
	Attack(ctx *gin.Context)
	Eyes(ctx *gin.Context)
	GetStrategy(ctx *gin.Context)
}

func NewRouter(authController ControllerProxy, m *middleware.Middleware) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 添加 CORS 中间件
	r.Use(corsHdl())
	g := r.Group("/api/v1")

	g.GET("/start", authController.Start)
	g.POST("/login", authController.Login)
	g.GET("/getStrategy", m.AuthMiddleware(), authController.GetStrategy)
	g.GET("/eyes", m.AuthMiddleware(), authController.Eyes)
	g.PUT("/attack", m.AuthMiddleware(), authController.Attack)
	g.GET("/paper", m.AuthMiddleware(), authController.Paper)

	// 处理未定义的路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(403, response.Err{Err: "警告!访问未知路由,请及时撤离以防被XXHBGS发现"})
	})

	return r
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 列出常用请求头或者设置为空的列表以允许任意头
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin", "X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// 允许所有请求来源
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}
