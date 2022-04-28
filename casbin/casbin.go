package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"log"
)

// PERM
const modelConfig = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`

var (
	globalEnforcer *WatcherEnforcer
)

func GetEnforcer() *WatcherEnforcer {
	return globalEnforcer
}

func Init() {
	m, _ := model.NewModelFromString(modelConfig)
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		log.Fatalln("NewEnforcer error:", err)
	}
	globalEnforcer = &WatcherEnforcer{Enforcer: enforcer}
	//
	LoadPoliciesData()
}

func LoadPoliciesData() {

	// first time load your policies data from adapter

	// ....

	//
}