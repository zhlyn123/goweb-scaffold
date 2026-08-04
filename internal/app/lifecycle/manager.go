package lifecycle

import (
	"context"
	"errors"
	"fmt"
)

// Hook 表示一个参与应用生命周期管理的组件。
type Hook struct {
	Name  string
	Start func(context.Context) error
	Stop  func(context.Context) error
}

// Manager 按注册顺序启动组件，并按相反顺序关闭组件。
// 这样依赖关系会更可控：例如数据库先于 HTTP 启动，HTTP 先于数据库关闭。
type Manager struct {
	hooks   []Hook
	started []Hook
}

// NewManager 创建一个空的生命周期管理器。
func NewManager() *Manager {
	return &Manager{}
}

// Register 添加一个生命周期钩子。
func (m *Manager) Register(hook Hook) {
	m.hooks = append(m.hooks, hook)
}

// Start 按顺序启动所有已注册的钩子。
func (m *Manager) Start(ctx context.Context) error {
	for _, hook := range m.hooks {
		if hook.Start == nil {
			m.started = append(m.started, hook)
			continue
		}

		if err := hook.Start(ctx); err != nil {
			return fmt.Errorf("%s: %w", hook.Name, err)
		}

		m.started = append(m.started, hook)
	}

	return nil
}

// Stop 按相反顺序关闭所有已启动的钩子。
func (m *Manager) Stop(ctx context.Context) error {
	var errs []error

	for i := len(m.started) - 1; i >= 0; i-- {
		hook := m.started[i]
		if hook.Stop == nil {
			continue
		}

		if err := hook.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", hook.Name, err))
		}
	}

	m.started = nil
	return errors.Join(errs...)
}
