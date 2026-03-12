package channel

import (
	"chat/globals"
	"chat/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type PlanManager struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled"`
	Plans   []Plan `json:"plans" mapstructure:"plans"`
}

type Plan struct {
	Level int        `json:"level" mapstructure:"level"`
	Price float32    `json:"price" mapstructure:"price"`
	Items []PlanItem `json:"items" mapstructure:"items"`
}

type PlanItem struct {
	Id     string   `json:"id" mapstructure:"id"`
	Name   string   `json:"name" mapstructure:"name"`
	Icon   string   `json:"icon" mapstructure:"icon"`
	Value  int64    `json:"value" mapstructure:"value"`
	Models []string `json:"models" mapstructure:"models"`
}

type Usage struct {
	Used  int64 `json:"used" mapstructure:"used"`
	Total int64 `json:"total" mapstructure:"total"`
}
type UsageMap map[string]Usage

var planExp int64 = 0

func NewPlanManager() *PlanManager {
	manager := &PlanManager{}
	if err := viper.UnmarshalKey("subscription", manager); err != nil {
		panic(err)
	}

	return manager
}

func (c *PlanManager) SaveConfig() error {
	return utils.SaveConfig("subscription", c)
}

func (c *PlanManager) UpdateConfig(data *PlanManager) error {
	c.Enabled = data.Enabled
	c.Plans = data.Plans
	return c.SaveConfig()
}

func (c *PlanManager) GetPlan(level int) Plan {
	for _, plan := range c.Plans {
		if plan.Level == level {
			return plan
		}
	}
	return Plan{}
}

func (c *PlanManager) GetPlans() []Plan {
	if c.Enabled {
		return c.Plans
	}

	return []Plan{}
}

func (c *PlanManager) GetRawPlans() []Plan {
	return c.Plans
}

func (c *PlanManager) IsEnabled() bool {
	return c.Enabled
}

func getOffsetFormat(offset time.Time, usage int64) string {
	return fmt.Sprintf("%s/%d", offset.Format("2006-01-02:15:04:05"), usage)
}

func GetSubscriptionUsage(cache *redis.Client, user globals.AuthLike, t string) (usage int64, offset time.Time) {
	// example cache value: 2021-09-01:19:00:00/100
	// if date is longer than 1 month, reset usage

	offset = time.Now()

	key := globals.GetSubscriptionLimitFormat(t, user.HitID())
	v, err := utils.GetCache(cache, key)
	if (err != nil && errors.Is(err, redis.Nil)) || len(v) == 0 {
		usage = 0
	}

	seg := strings.Split(v, "/")
	if len(seg) != 2 {
		usage = 0
	} else {
		date, err := time.Parse("2006-01-02:15:04:05", seg[0])
		usage = utils.ParseInt64(seg[1])
		if err != nil {
			usage = 0
		}

		// check if date is longer than current date after 1 month, if true, reset usage

		if date.AddDate(0, 1, 0).Before(time.Now()) {
			// date is longer than 1 month, reset usage
			usage = 0

			// get current date offset (1 month step)
			// example: 2021-09-01:19:00:0/100 -> 2021-10-01:19:00:00/100

			// copy date to offset
			offset = date

			// example:
			// current time: 2021-09-08:14:00:00
			// offset: 2021-07-01:19:00:00
			// expected offset: 2021-09-01:19:00:00
			// offset is not longer than current date, stop adding 1 month

			for offset.AddDate(0, 1, 0).Before(time.Now()) {
				offset = offset.AddDate(0, 1, 0)
			}
		} else {
			// date is not longer than 1 month, use current date value

			offset = date
		}
	}

	// set new cache value
	_ = utils.SetCache(cache, key, getOffsetFormat(offset, usage), planExp)

	return
}

func IncreaseSubscriptionUsage(cache *redis.Client, user globals.AuthLike, t string, limit int64) bool {
	key := globals.GetSubscriptionLimitFormat(t, user.HitID())
	usage, offset := GetSubscriptionUsage(cache, user, t)

	usage += 1
	if usage > limit {
		return false
	}

	// set new cache value
	err := utils.SetCache(cache, key, getOffsetFormat(offset, usage), planExp)
	return err == nil
}

func DecreaseSubscriptionUsage(cache *redis.Client, user globals.AuthLike, t string) bool {
	key := globals.GetSubscriptionLimitFormat(t, user.HitID())
	usage, offset := GetSubscriptionUsage(cache, user, t)

	usage -= 1
	if usage < 0 {
		return true
	}

	// set new cache value
	err := utils.SetCache(cache, key, getOffsetFormat(offset, usage), planExp)
	return err == nil
}

func ReleaseSubscriptionUsage(cache *redis.Client, user globals.AuthLike, t string) bool {
	key := globals.GetSubscriptionLimitFormat(t, user.HitID())
	_, offset := GetSubscriptionUsage(cache, user, t)

	// set new cache value
	err := utils.SetCache(cache, key, getOffsetFormat(offset, 0), planExp)
	return err == nil
}

func (p *Plan) GetUsage(user globals.AuthLike, db *sql.DB, cache *redis.Client) UsageMap {
	return utils.EachObject[PlanItem, Usage](p.Items, func(usage PlanItem) (string, Usage) {
		return usage.Id, usage.GetUsageForm(user, db, cache)
	})
}

func (p *PlanItem) GetUsage(user globals.AuthLike, db *sql.DB, cache *redis.Client) int64 {
	// preflight check
	user.GetID(db)
	usage, _ := GetSubscriptionUsage(cache, user, p.Id)
	return usage
}

func (p *PlanItem) ResetUsage(user globals.AuthLike, cache *redis.Client) bool {
	key := globals.GetSubscriptionLimitFormat(p.Id, user.HitID())
	_, offset := GetSubscriptionUsage(cache, user, p.Id)

	err := utils.SetCache(cache, key, getOffsetFormat(offset, 0), planExp)
	return err == nil
}

func (p *PlanItem) CreateUsage(user globals.AuthLike, cache *redis.Client) bool {
	key := globals.GetSubscriptionLimitFormat(p.Id, user.HitID())

	err := utils.SetCache(cache, key, getOffsetFormat(time.Now(), 0), planExp)
	return err == nil
}

func (p *PlanItem) GetUsageForm(user globals.AuthLike, db *sql.DB, cache *redis.Client) Usage {
	return Usage{
		Used:  p.GetUsage(user, db, cache),
		Total: p.Value,
	}
}

func (p *PlanItem) IsInfinity() bool {
	return p.Value == -1
}

func (p *PlanItem) IsExceeded(user globals.AuthLike, db *sql.DB, cache *redis.Client) bool {
	return p.IsInfinity() || p.GetUsage(user, db, cache) < p.Value
}

func (p *PlanItem) Increase(user globals.AuthLike, cache *redis.Client) bool {
	state := IncreaseSubscriptionUsage(cache, user, p.Id, p.Value)
	return state || p.IsInfinity()
}

func (p *PlanItem) Decrease(user globals.AuthLike, cache *redis.Client) bool {
	if p.Value == -1 {
		return true
	}
	return DecreaseSubscriptionUsage(cache, user, p.Id)
}

func (p *PlanItem) Release(user globals.AuthLike, cache *redis.Client) bool {
	return ReleaseSubscriptionUsage(cache, user, p.Id)
}

func (p *Plan) IncreaseUsage(user globals.AuthLike, cache *redis.Client, model string) bool {
	for _, usage := range p.Items {
		if utils.Contains(model, usage.Models) {
			return usage.Increase(user, cache)
		}
	}

	return false
}

func (p *Plan) DecreaseUsage(user globals.AuthLike, cache *redis.Client, model string) bool {
	for _, usage := range p.Items {
		if utils.Contains(model, usage.Models) {
			return usage.Decrease(user, cache)
		}
	}

	return false
}

func (p *Plan) ReleaseUsage(user globals.AuthLike, cache *redis.Client, model string) bool {
	for _, usage := range p.Items {
		if utils.Contains(model, usage.Models) {
			return usage.Release(user, cache)
		}
	}

	return false
}

func (p *Plan) ReleaseAll(user globals.AuthLike, cache *redis.Client) bool {
	for _, usage := range p.Items {
		if !usage.Release(user, cache) {
			return false
		}
	}

	return true
}

func IsValidPlan(level int) bool {
	return utils.InRange(level, 1, 3)
}
