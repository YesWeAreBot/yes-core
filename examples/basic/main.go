package main

import (
	"fmt"
	"time"

	"github.com/yeswearebot/yes-core/core"
)

// --- 倾听者插件 ---
type ListenerPlugin struct{}

func (l *ListenerPlugin) Name() string { return "listener" }
func (l *ListenerPlugin) Init(ctx *core.SystemContext) error {
	// 订阅 "greet" 事件
	ctx.Events.Subscribe("greet", func(payload any) {
		fmt.Printf("[Listener] 收到问候: %s\n", payload)
	})
	return nil
}
func (l *ListenerPlugin) Start(ctx *core.SystemContext) error { return nil }
func (l *ListenerPlugin) Stop(ctx *core.SystemContext) error  { return nil }

// --- 说话者插件 ---
type SpeakerPlugin struct{}

func (s *SpeakerPlugin) Name() string                       { return "speaker" }
func (s *SpeakerPlugin) Init(ctx *core.SystemContext) error { return nil }
func (s *SpeakerPlugin) Start(ctx *core.SystemContext) error {
	// 启动后 1 秒，发送一个问候
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("[Speaker] 正在发送问候...")
		ctx.Events.Publish("greet", "Hello from Speaker!")
	}()
	return nil
}
func (s *SpeakerPlugin) Stop(ctx *core.SystemContext) error { return nil }

func main() {
	app := core.NewApp()
	app.Register(&ListenerPlugin{})
	app.Register(&SpeakerPlugin{})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
