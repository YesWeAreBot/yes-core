package core

import (
	"context"
)

// Plugin 是整个 yes-core 微内核的基石接口。
// 无论是底层的网络适配器，还是上层的 AI 拟人逻辑，都必须实现这个接口。
// 框架通过这个接口统一管理所有组件的生命周期。
type Plugin interface {
	// Name 返回插件的唯一标识符。
	// 用于服务注册表 查找和依赖声明。
	Name() string

	// Init 初始化阶段。
	// 调用时机：所有插件注册完毕后，Start 之前。
	// 推荐操作：读取配置、向 Registry 暴露自己的服务、向 EventBus 订阅事件。
	// 禁止操作：不要在这里启动阻塞的网络请求或死循环。
	Init(ctx *SystemContext) error

	// Start 启动阶段。
	// 调用时机：所有插件 Init 完毕后。
	// 推荐操作：启动 Goroutine，建立 WebSocket 连接，开始监听端口等。
	Start(ctx *SystemContext) error

	// Stop 停止阶段。
	// 调用时机：收到系统中断信号 (Ctrl+C / Kill) 时，按启动的逆序调用。
	// 推荐操作：保存数据、关闭网络连接，实现优雅退出。
	Stop(ctx *SystemContext) error
}

// Dependent 依赖声明接口 (可选实现)。
// 如果插件 A 必须在插件 B 之前启动，A 就可以实现这个接口。
// 框架会根据这个接口返回的切片，使用 Kahn 算法进行拓扑排序，保证启动绝对安全。
type Dependent interface {
	DependsOn() []string
}

// SystemContext 是微内核传给插件的“万能钥匙”。
// 框架在调用生命周期的每一个阶段时，都会把这个上下文传给插件。
// 插件只能通过它与外部世界（其他插件）交互，实现了绝对的解耦。
type SystemContext struct {
	// GoContext 原生上下文，用于控制超时、传递取消信号。
	GoContext context.Context

	// Events 事件总线，负责跨插件的异步、弱耦合消息分发。
	Events EventBus

	// Registry 服务注册表，负责跨插件的同步、强耦合服务调用。
	Registry Registry
}
