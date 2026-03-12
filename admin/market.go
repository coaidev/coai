package admin

import (
	"chat/globals"
	"chat/utils"
	"fmt"

	"github.com/spf13/viper"
)

type ModelTag []string
type MarketModel struct {
	Id          string   `json:"id" mapstructure:"id" required:"true"`
	Name        string   `json:"name" mapstructure:"name" required:"true"`
	Description string   `json:"description" mapstructure:"description"`
	Default     bool     `json:"default" mapstructure:"default"`
	HighContext bool     `json:"high_context" mapstructure:"highcontext"`
	Avatar      string   `json:"avatar" mapstructure:"avatar"`
	Tag         ModelTag `json:"tag" mapstructure:"tag"`
}
type MarketModelList []MarketModel

type Market struct {
	Models MarketModelList `json:"models" mapstructure:"models"`
}

func NewMarket() *Market {
	var models MarketModelList
	if err := viper.UnmarshalKey("market", &models); err != nil {
		globals.Warn(fmt.Sprintf("[market] read config error: %s, use default config", err.Error()))
		models = MarketModelList{}
	}

	return &Market{
		Models: models,
	}
}

func (m *Market) GetModels() MarketModelList {
	return m.Models
}

func (m *Market) GetModel(id string) *MarketModel {
	for _, model := range m.Models {
		if model.Id == id {
			return &model
		}
	}
	return nil
}

func (m *Market) SaveConfig() error {
	return utils.SaveConfig("market", m.Models)
}

func (m *Market) SetModels(models MarketModelList) error {
	m.Models = models
	return m.SaveConfig()
}
