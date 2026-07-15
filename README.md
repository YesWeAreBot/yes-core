# yes-core

> A very tiny framework for any extendable event-driven system.

yes-core 是一个用 Go 语言编写的、极度纯粹的微内核。它的设计灵感来源于 Linux 内核的模块化思想与 Koishi 的插件机制。

## 架构图

graph TD
    subgraph 业务层
        A[业务插件 A: 群管]
        B[业务插件 B: AI 对话]
        C[业务插件 C: 自定义渲染器]
    end

    subgraph 抽象与调度层
        D[adapter-manager]
    end

    subgraph 协议适配层
        E[adapter-onebot]
        F[adapter-telegram...]
    end

    subgraph 微内核层
        G[yes-core]
    end

    A -. 依赖 .-> D
    B -. 依赖 .-> D
    C -. 重载渲染器 .-> D

    D -->|转换/分发事件| A
    D -->|调用接口| E

    E -->|WebSocket/HTTP| H(QQ/NapCat)
    G -->|管理生命周期| D
    G -->|管理生命周期| E


## 快速开始

### 1. 编写一个插件

任何结构体只要实现了 `Plugin` 接口，就可以成为 `yes-core` 的插件。推荐在包的 `init()` 函数中完成自动注册：

```go 
package my_plugin

import (
    "fmt"
    "github.com/yeswearebot/yes-core/core"
)

type MyPlugin struct{}

func init() {
    core.Register(func() core.Plugin { return &MyPlugin{} })
}

func (p *MyPlugin) Name() string { return “my-plugin” }

// 声明依赖 (可选): 保证 adapter-manager 先于自己启动
func (p *MyPlugin) DependsOn() []string { return []string{"adapter-manager"} }

func (p *MyPlugin) Init(ctx *core.SystemContext) error {
    // 1. 订阅事件
    ctx.Events.Subscribe("greet", func(payload any) {
      fmt.Printf("[MyPlugin] 收到问候: %s\n", payload)
    })
    return nil
}

func (p *MyPlugin) Start(ctx *core.SystemContext) error {
    // 2. 获取其他插件实例并调用 (强耦合同步调用)
    if rawManager, ok := ctx.Registry.Get(“adapter-manager”); ok {
        // manager := rawManager.(*manager.AdapterManager)
        // manager.SendMessage(…)
        _ = rawManager
    }

    // 3. 发布事件 (弱耦合异步广播)
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("[MyPlugin] 正在广播问候事件...")
        ctx.Events.Publish("greet", "Hello from MyPlugin!")
    }()

    return nil
}

func (p *MyPlugin) Stop(ctx *core.SystemContext) error { return nil }
```

### 2. 像搭积木一样组装系统

在 `main.go` 中，你只需要匿名引入需要的插件包，框架会自动完成实例化和启动：

```go
package main

import (
    "github.com/yeswearebot/yes-core/core"

    _ "github.com/yeswearebot/yes-core/plugins/adapter-manager" // 这些都是举例
    _ "github.com/yeswearebot/yes-core/plugins/onebot-adapter"
    _ "github.com/yeswearebot/yes-core/plugins/yesimbot-go"
    _ "github.com/yourname/my-plugin"
)

func main() {
    app := core.NewApp()
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

然后运行你的 `main.go`，启动你的系统！
