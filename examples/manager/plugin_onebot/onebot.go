package plugin_onebot

import (
	"fmt"

	"github.com/yeswearebot/yes-core/core"
	"github.com/yeswearebot/yes-core/examples/manager/plugin_manager"
)

type OneBotAdapter struct{}

func init() {
	core.Register(func() core.Plugin { return &OneBotAdapter{} })
}

func (o *OneBotAdapter) DependsOn() []string {
	return []string{"adapter-manager"} // 告诉框架：我依赖 manager，请先启动它
}

func (o *OneBotAdapter) Name() string { return "onebot-adapter" }
func (o *OneBotAdapter) Init(ctx *core.SystemContext) error {
	rawManager, ok := ctx.Registry.Get("adapter-manager")
	if !ok {
		return fmt.Errorf("adapter-manager not found")
	}
	manager := rawManager.(*plugin_manager.AdapterManager)
	manager.RegisterPlatform(o)
	return nil
}
func (o *OneBotAdapter) Start(ctx *core.SystemContext) error { return nil }
func (o *OneBotAdapter) Stop(ctx *core.SystemContext) error  { return nil }

func (o *OneBotAdapter) PlatformName() string { return "onebot" }
func (o *OneBotAdapter) SendMessage(channelID string, msg string) error {
	fmt.Printf("[OneBot] 正在向频道 %s 发送消息: %s\n", channelID, msg)
	return nil
}
