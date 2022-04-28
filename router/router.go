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

type ReqParam struct {
	Operation string `json:"operation"`
	policies [][]string `json:"policies"`
}

type ReadReq struct {
	policy []string `json:"policy"`
}

func InitRouter() {
	WatcherRouter = gin.Default()
	WatcherRouter.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	WatcherRouter.POST("policies/operation", func(context *gin.Context) {
		var req ReqParam
		err := context.ShouldBindJSON(&req)
		if err != nil {
			log.Println("ShouldBindJSON error:", err)
			context.JSON(http.StatusOK,"params error")
			return
		}
		if req.Operation == "add" {
			ok, err := casbin.GetEnforcer().ModifyBatchPolicies(casbin.GetEnforcer().AddPolicies, req.policies)
			if !ok || err != nil {
				log.Println("AddPolicies error:", err)
				context.JSON(http.StatusOK,"add policies error")
				return
			}
		} else if req.Operation == "remove" {
			ok, err := casbin.GetEnforcer().ModifyBatchPolicies(casbin.GetEnforcer().RemovePolicies, req.policies)
			if !ok || err != nil {
				log.Println("RemovePolicies error:", err)
				context.JSON(http.StatusOK,"add policies error")
				return
			}
		}
		context.JSON(http.StatusOK,"success")
	})
	WatcherRouter.GET("policies/read", func(context *gin.Context) {
		var req ReadReq
		err := context.ShouldBindJSON(&req)
		if err != nil {
			log.Println("ShouldBindJSON error:", err)
			context.JSON(http.StatusOK,"params error")
			return
		}
		ok,err := casbin.GetEnforcer().WEnforce(req.policy)
		if ok || err == nil {
			context.JSON(http.StatusOK,"success")
		}
		context.JSON(http.StatusForbidden,"you have not auth")
	})
	if err := WatcherRouter.Run(config.WatcherConfig.ListenPort); err != nil {
		log.Fatalln("Run error:", err)
	}
}