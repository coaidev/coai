package channel

import (
	"chat/globals"
	"chat/utils"

	"github.com/spf13/viper"
)

func NewChargeManager() *ChargeManager {
	var seq ChargeSequence
	if err := viper.UnmarshalKey("charge", &seq); err != nil {
		panic(err)
	}

	m := &ChargeManager{
		Sequence:         seq,
		Models:           map[string]*Charge{},
		NonBillingModels: []string{},
	}
	m.Load()

	return m
}

func (m *ChargeManager) Load() {
	seq := make(ChargeSequence, 0)
	for _, charge := range m.Sequence {
		if charge == nil {
			continue
		}
		if charge.Id == -1 {
			charge.Id = m.GetMaxId() + 1
		}
		seq = append(seq, charge)
	}
	m.Sequence = seq

	// init support models
	m.Models = map[string]*Charge{}
	for _, charge := range m.Sequence {
		for _, model := range charge.Models {
			if _, ok := m.Models[model]; !ok {
				m.Models[model] = charge
			}
		}
	}

	m.NonBillingModels = []string{}
	for _, charge := range m.Sequence {
		if !charge.IsBilling() {
			for _, model := range charge.Models {
				m.NonBillingModels = append(m.NonBillingModels, model)
			}
		}
	}
}

func (m *ChargeManager) GetModels() map[string]*Charge {
	return m.Models
}

func (m *ChargeManager) GetNonBillingModels() []string {
	return m.NonBillingModels
}

func (m *ChargeManager) IsBilling(model string) bool {
	return !utils.Contains(model, m.NonBillingModels)
}

func (m *ChargeManager) GetCharge(model string) *Charge {
	if charge, ok := m.Models[model]; ok {
		return charge
	}
	return &Charge{
		Type:      globals.NonBilling,
		Anonymous: false,
		Unset:     true,
	}
}

func (m *ChargeManager) SaveConfig() error {
	return utils.SaveConfig("charge", m.Sequence)
}

func (m *ChargeManager) GetMaxId() int {
	max := 0
	for _, charge := range m.Sequence {
		if charge.Id > max {
			max = charge.Id
		}
	}
	return max
}

func (m *ChargeManager) AddRawRule(charge *Charge) {
	charge.Id = m.GetMaxId() + 1
	m.Sequence = append(m.Sequence, charge)
}

func (m *ChargeManager) AddRule(charge Charge) error {
	m.AddRawRule(&charge)
	return m.SaveConfig()
}

func (m *ChargeManager) UpdateRawRule(charge *Charge) {
	for _, item := range m.Sequence {
		if item.Id == charge.Id {
			*item = *charge
			break
		}
	}
}

func (m *ChargeManager) UpdateRule(charge Charge) error {
	m.UpdateRawRule(&charge)
	return m.SaveConfig()
}

func (m *ChargeManager) SetRawRule(charge *Charge) {
	if charge.Id == -1 {
		m.AddRawRule(charge)
	} else {
		m.UpdateRawRule(charge)
	}
}

func (m *ChargeManager) SetRule(charge Charge) error {
	m.SetRawRule(&charge)
	return m.SaveConfig()
}

func (m *ChargeManager) DeleteRawRule(id int) {
	for i, item := range m.Sequence {
		if item.Id == id {
			m.Sequence = append(m.Sequence[:i], m.Sequence[i+1:]...)
			break
		}
	}
}

func (m *ChargeManager) DeleteRule(id int) error {
	m.DeleteRawRule(id)
	return m.SaveConfig()
}

func (m *ChargeManager) SyncRules(charge ChargeSequence, overwrite bool) error {
	for _, item := range charge {
		m.SyncRule(item, overwrite)
	}

	return m.SaveConfig()
}

func (m *ChargeManager) SyncRule(charge *Charge, overwrite bool) {
	if overwrite {
		m.SyncRuleWithOverwrite(charge)
	} else {
		m.SyncRuleWithoutOverwrite(charge)
	}
}

func (m *ChargeManager) SyncRuleWithOverwrite(charge *Charge) {
	if len(charge.Models) == 0 {
		return
	}

	for _, model := range charge.GetModels() {
		if raw := m.GetRuleByModel(model); raw != nil {
			if len(raw.Models) == 1 {
				// rule is already exist and only contains this model, just delete it

				m.DeleteRawRule(raw.Id)
			} else {
				// rule is already exist and contains other models, delete this model from it and add a new rule
				// delete model from raw rule
				raw.Models = utils.Filter(raw.Models, func(m string) bool {
					return m != model
				})
				m.UpdateRawRule(raw)
			}
		}
	}

	instance := charge.New("")
	instance.Models = charge.Models
	m.AddRawRule(instance)
}

func (m *ChargeManager) SyncRuleWithoutOverwrite(charge *Charge) {
	models := utils.Filter(charge.GetModels(), func(model string) bool {
		return !m.Contains(model)
	})

	if len(models) > 0 {
		charge.Models = models
		m.AddRawRule(charge)
	}
}

func (m *ChargeManager) ListRules() ChargeSequence {
	return m.Sequence
}

func (m *ChargeManager) Contains(model string) bool {
	for _, item := range m.Sequence {
		if item.Contains(model) {
			return true
		}
	}
	return false
}

func (m *ChargeManager) GetRule(id int) *Charge {
	for _, item := range m.Sequence {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (m *ChargeManager) GetRuleByModel(model string) *Charge {
	for _, item := range m.Sequence {
		if item.Contains(model) {
			return item
		}
	}
	return nil
}

func (c *Charge) IsUnsetType() bool {
	return c.Unset
}

func (c *Charge) GetType() string {
	if c.Type == "" {
		return globals.NonBilling
	}
	return c.Type
}

func (c *Charge) GetModels() []string {
	return c.Models
}

func (c *Charge) GetInput() float32 {
	if c.Input <= 0 {
		return 0
	}
	return c.Input
}

func (c *Charge) GetOutput() float32 {
	if c.Output <= 0 {
		return 0
	}
	return c.Output
}

func (c *Charge) SupportAnonymous() bool {
	return c.Anonymous
}

func (c *Charge) IsBilling() bool {
	return c.GetType() != globals.NonBilling
}

func (c *Charge) IsBillingType(t string) bool {
	return c.GetType() == t
}

func (c *Charge) GetLimit() float32 {
	switch c.GetType() {
	case globals.NonBilling:
		return 0
	case globals.TimesBilling:
		return c.GetOutput()
	case globals.TokenBilling:
		// 1k input tokens + 1k output tokens
		return c.GetInput() + c.GetOutput()
	default:
		return 0
	}
}

func (c *Charge) Contains(model string) bool {
	return utils.Contains(model, c.Models)
}

func (c *Charge) New(model string) *Charge {
	return &Charge{
		Type:      c.Type,
		Models:    []string{model},
		Input:     c.Input,
		Output:    c.Output,
		Anonymous: c.Anonymous,
	}
}
