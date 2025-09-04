package auth

import (
	"chat/channel"
	"chat/globals"
	"chat/utils"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

func disableSubscription() bool {
	return !channel.PlanInstance.IsEnabled()
}

func (u *User) GetSubscription(db *sql.DB) (time.Time, int) {
	if u.Subscription != nil && u.Subscription.Unix() > 0 {
		return *u.Subscription, u.Level
	}

	var expiredAt sql.NullTime
	if err := globals.QueryRowDb(db, "SELECT expired_at, level FROM subscription WHERE user_id = ?", u.GetID(db)).Scan(&expiredAt, &u.Level); err != nil {
		return time.Unix(0, 0), 0
	}

	var t *time.Time
	if expiredAt.Valid {
		t = utils.ToPtr(expiredAt.Time)
	} else {
		t = utils.ToPtr(time.Unix(0, 0))
	}

	u.Subscription = t
	return *u.Subscription, u.Level
}

func (u *User) GetSubscriptionLevel(db *sql.DB) int {
	_, level := u.GetSubscription(db)
	if !u.IsSubscribe(db) {
		return 0
	}
	return level
}

func (u *User) GetPlan(db *sql.DB) channel.Plan {
	return channel.PlanInstance.GetPlan(u.GetSubscriptionLevel(db))
}

func (u *User) GetSubscriptionExpiredAt(db *sql.DB) time.Time {
	stamp, _ := u.GetSubscription(db)
	return stamp
}

func (u *User) GetSubscriptionTime(db *sql.DB) time.Time {
	stamp, _ := u.GetSubscription(db)
	return stamp
}

func (u *User) IsSubscribe(db *sql.DB) bool {
	stamp, level := u.GetSubscription(db)
	return stamp.Unix() > time.Now().Unix() && level > 0
}

func (u *User) IsEnterprise(db *sql.DB) bool {
	if !u.IsSubscribe(db) {
		return false
	}

	var enterprise sql.NullBool
	if err := globals.QueryRowDb(db, "SELECT enterprise FROM subscription WHERE user_id = ?", u.GetID(db)).Scan(&enterprise); err != nil {
		return false
	}

	return enterprise.Valid && enterprise.Bool
}

func (u *User) GetSubscriptionExpiredDay(db *sql.DB) int {
	stamp := u.GetSubscriptionTime(db).Sub(time.Now())
	return int(math.Round(stamp.Hours() / 24))
}

func (u *User) AddSubscription(db *sql.DB, month int, level int) bool {
	current := u.GetSubscriptionTime(db)
	if current.Unix() < time.Now().Unix() {
		current = time.Now()
	}
	expiredAt := current.AddDate(0, month, 0)
	date := utils.ConvertSqlTime(expiredAt)
	_, err := globals.ExecDb(db, `
		INSERT INTO subscription (user_id, expired_at, total_month, level) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE expired_at = ?, total_month = total_month + ?, level = ?
	`, u.GetID(db), date, month, level, date, month, level)
	return err == nil
}

func (u *User) DowngradePlan(db *sql.DB, target int) error {
	expired, current := u.GetSubscription(db)
	if current == 0 || current == target {
		return fmt.Errorf("invalid plan level")
	}

	now := time.Now()
	weight := channel.PlanInstance.GetPlan(current).Price / channel.PlanInstance.GetPlan(target).Price
	stamp := float32(expired.Unix()-now.Unix()) * weight

	// ceil expired time
	expiredAt := now.Add(time.Duration(stamp)*time.Second).AddDate(0, 0, -1)
	date := utils.ConvertSqlTime(expiredAt)
	_, err := globals.ExecDb(db, "UPDATE subscription SET level = ?, expired_at = ? WHERE user_id = ?", target, date, u.GetID(db))

	return err
}

func (u *User) CountUpgradePrice(db *sql.DB, target int) float32 {
	expired := u.GetSubscriptionExpiredAt(db)
	weight := channel.PlanInstance.GetPlan(target).Price - u.GetPlan(db).Price
	if weight < 0 {
		return 0
	}

	days := expired.Sub(time.Now()).Hours() / 24
	return float32(days) * weight / 30
}

func (u *User) SetSubscriptionLevel(db *sql.DB, level int) bool {
	_, err := globals.ExecDb(db, "UPDATE subscription SET level = ? WHERE user_id = ?", level, u.GetID(db))
	return err == nil
}

func CountSubscriptionPrize(level int, month int) float32 {
	plan := channel.PlanInstance.GetPlan(level)
	base := plan.Price * float32(month)
	if month >= 36 {
		return base * 0.7
	} else if month >= 12 {
		return base * 0.8
	} else if month >= 6 {
		return base * 0.9
	}
	return base
}

func BuySubscription(db *sql.DB, cache *redis.Client, user *User, level int, month int) error {
	if disableSubscription() {
		return errors.New("subscription feature does not enable of this site")
	}

	if month < 1 || month > 999 || !channel.IsValidPlan(level) {
		return errors.New("invalid subscription params")
	}

	before := user.GetSubscriptionLevel(db)
	if before == 0 || before == level {
		// buy new subscription or renew subscription
		money := CountSubscriptionPrize(level, month)
		if user.Pay(db, cache, money) {
			// migrate subscription
			user.AddSubscription(db, month, level)

			if before == 0 {
				// new subscription

				plan := user.GetPlan(db)
				for _, usage := range plan.Items {
					// create usage
					usage.CreateUsage(user, cache)
				}
			}

			return nil
		}
	} else if before > level {
		// downgrade subscription
		return user.DowngradePlan(db, level)
	} else {
		// upgrade subscription
		money := user.CountUpgradePrice(db, level)
		if user.Pay(db, cache, money) {
			user.SetSubscriptionLevel(db, level)
			return nil
		}
	}

	return errors.New("not enough money")
}

func HandleSubscriptionUsage(db *sql.DB, cache *redis.Client, user *User, model string) bool {
	if disableSubscription() {
		return false
	}
	plan := user.GetPlan(db)
	return plan.IncreaseUsage(user, cache, model)
}

func RevertSubscriptionUsage(db *sql.DB, cache *redis.Client, user *User, model string) bool {
	if disableSubscription() {
		return false
	}
	plan := user.GetPlan(db)
	return plan.DecreaseUsage(user, cache, model)
}
