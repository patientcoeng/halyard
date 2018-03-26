package alerting

type Alert interface {
	Trigger(summary string, detail string)
}

type Manager struct {
	alerts []Alert
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddAlert(alert Alert) {
	m.alerts = append(m.alerts, alert)
}

func (m *Manager) Trigger(summary string, detail string) {
	for _, alert := range m.alerts {
		alert.Trigger(summary, detail)
	}
}
