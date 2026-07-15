package plugin_ai

import (
	"fmt"
	"time"

	"github.com/yeswearebot/yes-core/core"
	"github.com/yeswearebot/yes-core/examples/manager/plugin_manager"
)

type YesimbotAI struct{}

func init() {
	core.Register(func() core.Plugin { return &YesimbotAI{} })
}

func (y *YesimbotAI) DependsOn() []string {
	return []string{"adapter-manager"} // AI 也依赖 manager，但不依赖 OneBot (因为它是通过 Manager 发消息的)
}

func (y *YesimbotAI) Name() string                       { return "yesimbot-ai" }
func (y *YesimbotAI) Init(ctx *core.SystemContext) error { return nil }
func (y *YesimbotAI) Start(ctx *core.SystemContext) error {
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("[YesimbotAI] 思考完毕，决定回复群友。")

		rawManager, ok := ctx.Registry.Get("adapter-manager")
		if !ok {
			return
		}
		manager := rawManager.(*plugin_manager.AdapterManager)
		manager.SendMessage("onebot", "群123456", "你们好，我是群友！(通过全局注册)")
	}()
	return nil
}
func (y *YesimbotAI) Stop(ctx *core.SystemContext) error { return nil }
