package main

import (
	"github.com/kdkyg/casbin_watcher/casbin"
	"github.com/kdkyg/casbin_watcher/config"
	"github.com/kdkyg/casbin_watcher/router"
)

func main() {
	config.Init()

	casbin.Init()

	go casbin.StartWatcher()

	router.InitRouter()
}
