package core

import "sync"

// Registry 服务注册表接口。
// 解决依赖注入 (DI) 的核心。如果 A 插件想要直接调用 B 插件的方法，它们通过 Registry 寻找彼此。
type Registry interface {
	RegisterPlugin(plugin Plugin)
	Get(pluginName string) (Plugin, bool)
}

// defaultRegistry 注册表的默认实现。
type defaultRegistry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

// NewRegistry 创建一个新的注册表实例。
func NewRegistry() Registry {
	return &defaultRegistry{
		plugins: make(map[string]Plugin),
	}
}

func (r *defaultRegistry) RegisterPlugin(plugin Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 以插件的 Name() 作为 key 存储
	r.plugins[plugin.Name()] = plugin
}

func (r *defaultRegistry) Get(pluginName string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[pluginName]
	return p, ok
}
