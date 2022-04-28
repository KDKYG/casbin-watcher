package casbin

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kdkyg/casbin_watcher/config"
	watcher "github.com/kdkyg/casbin_watcher/redis-watcher"
	"log"
)

func StartWatcher() {
	redisAddress := fmt.Sprintf("%s:%d", config.WatcherConfig.RedisConf.Host, config.WatcherConfig.RedisConf.Port)
	w, err := watcher.NewWatcher(redisAddress, watcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: config.WatcherConfig.RedisConf.Password,
			DB:       config.WatcherConfig.RedisConf.Instance,
		},
		Channel:    "/casbin",
		IgnoreSelf: true,
	})
	if err != nil {
		panic(err)
	}
	globalEnforcer.SetWatcher(w)
	w.SetUpdateCallback(updateCallback)
}

func updateCallback(msg string) {
	fmt.Println("update callback:", msg)
	globalEnforcer.m.Lock()
	globalEnforcer.EnableAutoNotifyWatcher(false)
	defer func() {
		globalEnforcer.EnableAutoNotifyWatcher(true)
		globalEnforcer.m.Unlock()
	}()

	m := &watcher.MSG{}
	if err := m.UnmarshalBinary([]byte(msg)); err != nil {
		log.Println("UnmarshalBinary error:", err)
		return
	}
	if m.Method == "UpdateForAddPolicies" {
		rules := interface2rules(m.Params)
		log.Println("update rules:", rules)
		if m.Sec == m.Ptype && m.Sec == "p" {
			ok, err := globalEnforcer.AddPolicies(rules)
			if !ok || err != nil {
				log.Println("AddPolicies error, bool:", ok, "err:", err)
			}
		} else if m.Sec == m.Ptype && m.Sec == "g" {
			ok, err := globalEnforcer.AddGroupingPolicies(rules)
			if !ok || err != nil {
				log.Println("AddGroupingPolicies error, bool:", ok, "err:", err)
			}
		} else {
			log.Println("type error: ", m.Ptype)
		}
	} else if m.Method == "UpdateForRemovePolicies" {
		rules := interface2rules(m.Params)
		log.Println("update rules:", rules)
		if m.Sec == m.Ptype && m.Sec == "p" {
			ok, err := globalEnforcer.RemovePolicies(rules)
			if !ok || err != nil {
				log.Println("RemovePolicies error, bool:", ok, "err:", err)
			}
		} else if m.Sec == m.Ptype && m.Sec == "g" {
			ok, err := globalEnforcer.RemoveGroupingPolicies(rules)
			if !ok || err != nil {
				log.Println("RemoveGroupingPolicies error, bool:", ok, "err:", err)
			}
		} else {
			log.Println("type error: ", m.Ptype)
		}
	} else {
		log.Println("method error: ", m.Method)
	}
}

func interface2rules(i interface{}) [][]string {
	rules := make([][]string, 0)
	for _, item := range i.([]interface{}) {
		tmp := item.([]interface{})
		rule := make([]string, 0)
		for _, t := range tmp {
			rule = append(rule, t.(string))
		}
		rules = append(rules, rule)
	}
	return rules
}