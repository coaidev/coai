package channel

import (
	"chat/globals"
	"chat/utils"
	"errors"
	"time"

	"github.com/spf13/viper"
)

var ConduitInstance *Manager
var ChargeInstance *ChargeManager
var SystemInstance *SystemConfig
var PlanInstance *PlanManager

func InitManager() {
	ConduitInstance = NewChannelManager()
	ChargeInstance = NewChargeManager()
	SystemInstance = NewSystemConfig()
	PlanInstance = NewPlanManager()
}

func NewChannelManager() *Manager {
	var seq Sequence
	if err := viper.UnmarshalKey("channel", &seq); err != nil {
		panic(err)
	}

	// sort by priority

	manager := &Manager{
		Sequence:          seq,
		Models:            []string{},
		PreflightSequence: map[string]Sequence{},
	}
	manager.Load()

	return manager
}

func (m *Manager) Load() {
	// load channels
	for _, channel := range m.Sequence {
		if channel != nil {
			channel.Load()
		}
	}

	// init support models
	m.Models = []string{}
	for _, channel := range m.GetActiveSequence() {
		for _, model := range channel.GetHitModels() {
			if !utils.Contains(model, m.Models) {
				m.Models = append(m.Models, model)
			}
		}
	}

	// init preflight sequence
	m.PreflightSequence = map[string]Sequence{}
	for _, model := range m.Models {
		var seq Sequence
		for _, channel := range m.GetActiveSequence() {
			if channel.IsHit(model) {
				seq = append(seq, channel)
			}
		}
		seq.Sort()
		m.PreflightSequence[model] = seq
	}

	stamp := time.Now().Unix()

	globals.SupportModels = m.Models
	globals.V1ListModels = globals.ListModels{
		Object: "list",
		Data: utils.Each(m.Models, func(model string) globals.ListModelsItem {
			return globals.ListModelsItem{
				Id:      model,
				Object:  "model",
				Created: stamp,
				OwnedBy: "system",
			}
		}),
	}
}

func (m *Manager) GetSequence() Sequence {
	return m.Sequence
}

func (m *Manager) GetActiveSequence() Sequence {
	var seq Sequence
	for _, channel := range m.Sequence {
		if channel.GetState() {
			seq = append(seq, channel)
		}
	}
	seq.Sort()
	return seq
}

func (m *Manager) GetModels() []string {
	return m.Models
}

func (m *Manager) GetPreflightSequence() map[string]Sequence {
	return m.PreflightSequence
}

// HitSequence returns the preflight sequence of the model
func (m *Manager) HitSequence(model string) Sequence {
	return m.PreflightSequence[model]
}

// HasChannel returns whether the channel exists
func (m *Manager) HasChannel(model string) bool {
	return utils.Contains(model, m.Models)
}

func (m *Manager) GetTicker(model, group string) *Ticker {
	if !m.HasChannel(model) {
		return nil
	}

	return NewTicker(m.HitSequence(model), group)
}

func (m *Manager) Len() int {
	return len(m.Sequence)
}

func (m *Manager) GetMaxId() int {
	var max int
	for _, channel := range m.Sequence {
		if channel.Id > max {
			max = channel.Id
		}
	}
	return max
}

func (m *Manager) SaveConfig() error {
	return utils.SaveConfig("channel", m.Sequence)
}

func (m *Manager) CreateChannel(channel *Channel) error {
	channel.Id = m.GetMaxId() + 1
	m.Sequence = append(m.Sequence, channel)
	return m.SaveConfig()
}

func (m *Manager) UpdateChannel(id int, channel *Channel) error {
	for i, item := range m.Sequence {
		if item.Id == id {
			m.Sequence[i] = channel
			return m.SaveConfig()
		}
	}
	return errors.New("channel not found")
}

func (m *Manager) DeleteChannel(id int) error {
	for i, item := range m.Sequence {
		if item.Id == id {
			m.Sequence = append(m.Sequence[:i], m.Sequence[i+1:]...)
			return m.SaveConfig()
		}
	}
	return errors.New("channel not found")
}

func (m *Manager) ActivateChannel(id int) error {
	for i, item := range m.Sequence {
		if item.Id == id {
			m.Sequence[i].State = true
			return m.SaveConfig()
		}
	}
	return errors.New("channel not found")
}

func (m *Manager) DeactivateChannel(id int) error {
	for i, item := range m.Sequence {
		if item.Id == id {
			m.Sequence[i].State = false
			return m.SaveConfig()
		}
	}
	return errors.New("channel not found")
}
