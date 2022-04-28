package casbin

import (
	"github.com/casbin/casbin/v2"
	"sync"
)

type WatcherEnforcer struct {
	*casbin.Enforcer
	m sync.RWMutex
}

type ModifyBatchPoliciesFunc func(rules [][]string) (bool, error)

func (m *WatcherEnforcer) WEnforce (rvals ...interface{}) (bool, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	return m.Enforcer.Enforce(rvals...)
}

func (m *WatcherEnforcer) ModifyBatchPolicies(f ModifyBatchPoliciesFunc, rules [][]string) (bool, error) {
	m.m.Lock()
	defer m.m.Unlock()
	return f(rules)
}