package core

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

// App 代表整个微内核实例，是系统的大脑。
type App struct {
	plugins  []Plugin
	registry Registry
	eventBus EventBus
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewApp 创建一个新的内核实例。
// 它会自动监听系统的中断信号，并收集全局注册的插件。
func NewApp() *App {
	// 监听 Ctrl+C (SIGINT) 和 kill 命令 (SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	app := &App{
		ctx:      ctx,
		cancel:   cancel,
		registry: NewRegistry(),
		eventBus: NewEventBus(),
	}

	// 自动收集通过包 init() 机制全局注册的插件工厂，实例化它们
	for _, factory := range globalFactories {
		app.plugins = append(app.plugins, factory())
	}

	return app
}

// Register 提供手动注册插件的能力 (多用于测试)。
func (a *App) Register(p Plugin) {
	a.plugins = append(a.plugins, p)
}

// resolveDependencies 使用 Kahn 算法进行有向无环图 (DAG) 拓扑排序。
// 目的：根据插件声明的 DependsOn()，重排启动顺序，保证被依赖的插件先初始化。
func (a *App) resolveDependencies() error {
	graph := make(map[string][]string) // 邻接表：记录谁依赖谁
	inDegree := make(map[string]int)   // 入度表：记录该插件被多少个插件依赖
	pluginMap := make(map[string]Plugin)

	// 初始化数据结构
	for _, p := range a.plugins {
		name := p.Name()
		pluginMap[name] = p
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
	}

	// 构建依赖图
	for _, p := range a.plugins {
		// 使用类型断言检查插件是否实现了 Dependent 接口
		if dep, ok := p.(Dependent); ok {
			for _, req := range dep.DependsOn() {
				if _, exists := pluginMap[req]; !exists {
					return fmt.Errorf("plugin %s depends on %s, but %s is not registered", p.Name(), req, req)
				}
				// req -> p (p 依赖 req，所以 req 指向 p)
				graph[req] = append(graph[req], p.Name())
				inDegree[p.Name()]++
			}
		}
	}

	// 将所有入度为 0 的插件 (没有依赖的插件) 加入队列作为起点
	var queue []string
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	// BFS 遍历
	var sorted []Plugin
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		sorted = append(sorted, pluginMap[curr])

		for _, neighbor := range graph[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 如果排序后的数量不等于总数量，说明图中存在环 (循环依赖)
	if len(sorted) != len(a.plugins) {
		return fmt.Errorf("circular dependency detected among plugins")
	}

	a.plugins = sorted
	return nil
}

// Run 轰鸣引擎！这是框架最核心的生命周期调度方法。
func (a *App) Run() error {
	// 阶段 1: 依赖解析与排序
	if err := a.resolveDependencies(); err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}
	fmt.Printf("[yes-core] Plugin load order: ")
	for _, p := range a.plugins {
		fmt.Printf("%s ", p.Name())
	}
	fmt.Println()

	// 阶段 2: 注册到 Registry
	// 必须在 Init 之前注册，这样插件在 Init 时就能互相发现了
	for _, p := range a.plugins {
		a.registry.RegisterPlugin(p)
	}

	// 构建传递给插件的全局上下文
	sysCtx := &SystemContext{
		GoContext: a.ctx,
		Events:    a.eventBus,
		Registry:  a.registry,
	}

	// 阶段 3: 依次执行 Init (按拓扑排序)
	fmt.Println("[yes-core] Initializing plugins...")
	for _, p := range a.plugins {
		if err := p.Init(sysCtx); err != nil {
			return fmt.Errorf("plugin %s init failed: %w", p.Name(), err)
		}
	}

	// 阶段 4: 依次执行 Start
	fmt.Println("[yes-core] Starting plugins...")
	for _, p := range a.plugins {
		if err := p.Start(sysCtx); err != nil {
			return fmt.Errorf("plugin %s start failed: %w", p.Name(), err)
		}
	}
	fmt.Println("[yes-core] All plugins started successfully. System is running.")

	// 阶段 5: 阻塞等待退出信号
	// 当用户按下 Ctrl+C，ctx.Done() 通道会关闭，解除阻塞
	<-a.ctx.Done()
	fmt.Println("\n[yes-core] Shutdown signal received, stopping plugins...")

	// 阶段 6: 优雅退出 (反向 Stop)
	// 为什么反向？因为后启动的插件往往依赖先启动的插件。
	// 比如适配器先启动，AI 后启动。退出时应该先让 AI 停止处理，再断开适配器网络连接。
	for i := len(a.plugins) - 1; i >= 0; i-- {
		p := a.plugins[i]
		if err := p.Stop(sysCtx); err != nil {
			fmt.Printf("[yes-core] Warning: plugin %s stop failed: %v\n", p.Name(), err)
		}
	}

	fmt.Println("[yes-core] All plugins stopped. Goodbye!")
	return nil
}
