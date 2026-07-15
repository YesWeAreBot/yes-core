package core

// PluginFactory 插件工厂函数签名。
// 使用工厂函数而不是直接注册实例，是为了保证每次启动 App 时，插件都是全新干净的内存状态。
type PluginFactory func() Plugin

// globalFactories 全局插件工厂列表。
// 在 Go 中，利用包的 init() 机制，不同包可以在编译时将自己注册到这里。
var globalFactories []PluginFactory

// Register 全局注册插件工厂。
// 推荐在插件包的 init() 函数中调用此方法。
// 这样在 main.go 中只需匿名 import 插件包，即可完成插件的自动挂载。
func Register(factory PluginFactory) {
	globalFactories = append(globalFactories, factory)
}
