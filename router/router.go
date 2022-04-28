package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kdkyg/casbin_watcher/casbin"
	"github.com/kdkyg/casbin_watcher/config"
	"log"
	"net/http"
)

var (
	WatcherRouter *gin.Engine
)

func InitRouter() {
	WatcherRouter = gin.Default()
	WatcherRouter.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	WatcherRouter.PUT("/policies", func(context *gin.Context) {
		var data interface{}
		context.ShouldBindJSON(&data)
		err := context.ShouldBindJSON(&data)
		if err != nil {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = casbin.GetEnforcer().AddPolicies(casbin.Interface2rules(data))
		if err != nil {
			http.Error(context.Writer, err.Error(), http.StatusServiceUnavailable)
			return
		}
	})
	WatcherRouter.DELETE("/policies",func(context *gin.Context) {
		var data interface{}
		context.ShouldBindJSON(&data)
		err := context.ShouldBindJSON(&data)
		if err != nil {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = casbin.GetEnforcer().RemovePolicy(casbin.Interface2rules(data))
		if err != nil {
			http.Error(context.Writer, err.Error(), http.StatusServiceUnavailable)
			return
		}
	})
	WatcherRouter.GET("/enforcer",func(context *gin.Context) {
		var data []interface{}
		err := context.ShouldBindJSON(&data)
		if err != nil {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}
		ok,err := casbin.GetEnforcer().WEnforce(data...)
		if ok || err == nil {
			context.JSON(http.StatusOK,"Authorized")
		}
		context.JSON(http.StatusBadRequest,"Unauthorized")
	})
	if err := WatcherRouter.Run(config.WatcherConfig.ListenPort); err != nil {
		log.Fatalln("Run error:", err)
	}
}