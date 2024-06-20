package marker

import (
	"fmt"
	"strings"
	"sync"
)

type Marker struct {
	mu         sync.Mutex
	Strategies map[string]Strategy
}

type Strategy struct {
	Login    string
	Company  string
	Location string
}

func New() *Marker {
	return &Marker{
		Strategies: make(map[string]Strategy),
	}
}

func (m *Marker) AddStrategies(strategies ...string) error {
	for _, strategy := range strategies {
		parts := strings.Split(strategy, ",")
		if len(parts) != 3 {
			return fmt.Errorf("invalid marker strategies format: %s", strategy)
		}

		login := strings.Trim(parts[0], "` ")
		company := strings.Trim(parts[1], "` ")
		location := strings.Trim(parts[2], "` ")

		m.mu.Lock()
		m.Strategies[login] = Strategy{
			Login:    login,
			Company:  company,
			Location: location,
		}
		m.mu.Unlock()
	}
	return nil
}

func (m *Marker) DeleteStrategy(login string) {
	m.mu.Lock()
	delete(m.Strategies, login)
	m.mu.Unlock()
}

func (m *Marker) Marks() []Strategy {
	var ss []Strategy
	for _, strategy := range m.Strategies {
		ss = append(ss, strategy)
	}
	return ss
}
