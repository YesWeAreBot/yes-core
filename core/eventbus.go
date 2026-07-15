package core

import "sync"

// EventHandler 事件处理函数的统一签名。
// payload 使用 any (泛型)，意味着事件总线不关心业务数据结构，完全由业务层自行断言。
type EventHandler func(payload any)

// EventBus 事件总线接口。
// 极其轻量的发布/订阅 中心，系统神经的总枢纽。
type EventBus interface {
	Subscribe(topic string, handler EventHandler)
	Publish(topic string, payload any)
}

// defaultEventBus 事件总线的默认实现。
type defaultEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

// NewEventBus 创建一个新的事件总线实例。
func NewEventBus() EventBus {
	return &defaultEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *defaultEventBus) Subscribe(topic string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// 将订阅者追加到对应 topic 的切片中
	b.handlers[topic] = append(b.handlers[topic], handler)
}

func (b *defaultEventBus) Publish(topic string, payload any) {
	b.mu.RLock()
	// 【并发安全核心细节】：
	// 这里必须深拷贝一份 handlers 列表。
	// 如果不拷贝，在 for 循环异步执行 handler 时，如果有其他协程刚好 Subscribe 修改了原切片，会导致 panic。
	handlers := make([]EventHandler, len(b.handlers[topic]))
	copy(handlers, b.handlers[topic])
	b.mu.RUnlock()

	// 异步广播给所有订阅者，防止某个慢处理者阻塞发布者
	for _, h := range handlers {
		go h(payload)
	}
}
