package plugin_manager

import (
	"fmt"

	"github.com/yeswearebot/yes-core/core"
)

// BotPlatform 统一的通信平台接口
type BotPlatform interface {
	PlatformName() string
	SendMessage(channelID string, msg string) error
}

type AdapterManager struct {
	platforms map[string]BotPlatform
}

// 在包加载时自动注册
func init() {
	core.Register(func() core.Plugin {
		return &AdapterManager{platforms: make(map[string]BotPlatform)}
	})
}

func (m *AdapterManager) Name() string                        { return "adapter-manager" }
func (m *AdapterManager) Init(ctx *core.SystemContext) error  { return nil }
func (m *AdapterManager) Start(ctx *core.SystemContext) error { return nil }
func (m *AdapterManager) Stop(ctx *core.SystemContext) error  { return nil }

func (m *AdapterManager) RegisterPlatform(p BotPlatform) {
	m.platforms[p.PlatformName()] = p
	fmt.Printf("[Manager] 已注册平台: %s\n", p.PlatformName())
}

func (m *AdapterManager) SendMessage(platform, channelID, msg string) {
	if p, ok := m.platforms[platform]; ok {
		p.SendMessage(channelID, msg)
	} else {
		fmt.Printf("[Manager] 找不到平台: %s\n", platform)
	}
}
