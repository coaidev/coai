package auth

import (
	"chat/globals"
	"chat/utils"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

// PayTheFly configuration keys
const (
	payTheFlyProjectID  = "paythefly.project_id"
	payTheFlyProjectKey = "paythefly.project_key"
	payTheFlyChainID    = "paythefly.chain_id"
	payTheFlyToken      = "paythefly.token_address"
	payTheFlyEnabled    = "paythefly.enabled"
)

// PayTheFlyWebhookBody represents the incoming webhook payload envelope
type PayTheFlyWebhookBody struct {
	Data      string `json:"data"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
}

// PayTheFlyWebhookPayload represents the parsed webhook data content
type PayTheFlyWebhookPayload struct {
	Value     string `json:"value"`
	Confirmed bool   `json:"confirmed"`
	SerialNo  string `json:"serial_no"`
	TxHash    string `json:"tx_hash"`
	Wallet    string `json:"wallet"`
	TxType    int    `json:"tx_type"` // 1=payment, 2=withdrawal
}

// PayTheFlyCreateOrderForm represents the request to create a payment order
type PayTheFlyCreateOrderForm struct {
	Quota int `json:"quota" binding:"required"`
}

// usePayTheFly checks if PayTheFly payment is enabled
func usePayTheFly() bool {
	return viper.GetBool(payTheFlyEnabled)
}

// getPayTheFlyProjectID returns the configured project ID
func getPayTheFlyProjectID() string {
	return viper.GetString(payTheFlyProjectID)
}

// getPayTheFlyProjectKey returns the configured project key for webhook HMAC verification
func getPayTheFlyProjectKey() string {
	return viper.GetString(payTheFlyProjectKey)
}

// getPayTheFlyChainID returns the configured blockchain chain ID (default: 56 for BSC)
func getPayTheFlyChainID() int {
	chainID := viper.GetInt(payTheFlyChainID)
	if chainID == 0 {
		return 56 // default to BSC
	}
	return chainID
}

// getPayTheFlyTokenAddress returns the configured payment token contract address
func getPayTheFlyTokenAddress() string {
	return viper.GetString(payTheFlyToken)
}

// generatePayTheFlySerialNo generates a unique serial number for the payment order
func generatePayTheFlySerialNo() string {
	return utils.Sha2Encrypt(utils.GenerateChar(32))
}

// BuildPayTheFlyPaymentURL constructs the PayTheFly payment redirect URL.
// amount is human-readable (e.g., "0.01"), not raw token units.
func BuildPayTheFlyPaymentURL(serialNo string, amount string, deadline int64, signature string) string {
	return fmt.Sprintf(
		"https://pro.paythefly.com/pay?chainId=%d&projectId=%s&amount=%s&serialNo=%s&deadline=%d&signature=%s&token=%s",
		getPayTheFlyChainID(),
		getPayTheFlyProjectID(),
		amount,
		serialNo,
		deadline,
		signature,
		getPayTheFlyTokenAddress(),
	)
}

// verifyPayTheFlyWebhookSignature verifies the HMAC-SHA256 signature using timing-safe comparison.
// Signature = HMAC-SHA256(data + "." + timestamp, projectKey)
func verifyPayTheFlyWebhookSignature(data string, timestamp int64, signature string) bool {
	projectKey := getPayTheFlyProjectKey()
	if projectKey == "" {
		globals.Warn("[PayTheFly] project key is not configured")
		return false
	}

	message := fmt.Sprintf("%s.%d", data, timestamp)
	mac := hmac.New(sha256.New, []byte(projectKey))
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// timing-safe comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(expectedMAC), []byte(signature)) == 1
}

// PayTheFlyCreateOrderAPI creates a PayTheFly crypto payment order
func PayTheFlyCreateOrderAPI(c *gin.Context) {
	if !usePayTheFly() {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "PayTheFly payment is not enabled",
		})
		return
	}

	user := GetUserByCtx(c)
	if user == nil {
		return
	}

	var form PayTheFlyCreateOrderForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	if form.Quota <= 0 || form.Quota > 99999 {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "invalid quota range (1 ~ 99999)",
		})
		return
	}

	// Convert quota to payment amount (same ratio as existing: quota * 0.1)
	amount := float64(form.Quota) * 0.1
	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)

	serialNo := generatePayTheFlySerialNo()
	deadline := time.Now().Add(30 * time.Minute).Unix()

	// Store pending order in database
	db := utils.GetDBFromContext(c)
	_, err := globals.ExecDb(db, `
		INSERT INTO paythefly_orders (serial_no, user_id, quota, amount, deadline, status, created_at)
		VALUES (?, ?, ?, ?, ?, 'pending', NOW())
	`, serialNo, user.GetID(db), form.Quota, amountStr, deadline)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "failed to create payment order",
		})
		return
	}

	// NOTE: In production, the signature should be computed server-side using
	// EIP-712 typed data signing with the configured private key.
	// Domain: { name: "PayTheFlyPro", version: "1" }
	// The private key is read from environment/config (paythefly.private_key).
	// This endpoint returns the payment URL for the frontend to redirect.

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"serial_no":  serialNo,
		"amount":     amountStr,
		"deadline":   deadline,
		"chain_id":   getPayTheFlyChainID(),
		"project_id": getPayTheFlyProjectID(),
		"token":      getPayTheFlyTokenAddress(),
	})
}

// PayTheFlyWebhookAPI handles webhook callbacks from PayTheFly.
// Webhook body: { "data": "<json string>", "sign": "<hmac hex>", "timestamp": <unix> }
// Response must contain "success" string.
func PayTheFlyWebhookAPI(c *gin.Context) {
	if !usePayTheFly() {
		c.String(http.StatusOK, "success")
		return
	}

	var body PayTheFlyWebhookBody
	if err := c.ShouldBindJSON(&body); err != nil {
		globals.Warn(fmt.Sprintf("[PayTheFly] invalid webhook body: %s", err.Error()))
		c.String(http.StatusBadRequest, "invalid request")
		return
	}

	// Verify HMAC-SHA256 signature (timing-safe)
	if !verifyPayTheFlyWebhookSignature(body.Data, body.Timestamp, body.Sign) {
		globals.Warn("[PayTheFly] webhook signature verification failed")
		c.String(http.StatusForbidden, "invalid signature")
		return
	}

	// Reject stale webhooks (> 5 minutes old)
	if time.Now().Unix()-body.Timestamp > 300 {
		globals.Warn("[PayTheFly] webhook timestamp too old")
		c.String(http.StatusBadRequest, "timestamp expired")
		return
	}

	// Parse the inner data payload
	var payload PayTheFlyWebhookPayload
	if err := json.Unmarshal([]byte(body.Data), &payload); err != nil {
		globals.Warn(fmt.Sprintf("[PayTheFly] failed to parse webhook data: %s", err.Error()))
		c.String(http.StatusBadRequest, "invalid data")
		return
	}

	// Only process confirmed payment transactions (tx_type=1)
	if payload.TxType != 1 {
		globals.Info(fmt.Sprintf("[PayTheFly] ignoring non-payment transaction type: %d", payload.TxType))
		c.String(http.StatusOK, "success")
		return
	}

	if !payload.Confirmed {
		globals.Info(fmt.Sprintf("[PayTheFly] payment not yet confirmed: %s", payload.SerialNo))
		c.String(http.StatusOK, "success")
		return
	}

	// Process the confirmed payment
	db := utils.GetDBFromContext(c)
	if err := processPayTheFlyPayment(db, &payload); err != nil {
		globals.Warn(fmt.Sprintf("[PayTheFly] failed to process payment: %s", err.Error()))
		c.String(http.StatusOK, "success")
		return
	}

	globals.Info(fmt.Sprintf("[PayTheFly] payment processed: serial_no=%s tx_hash=%s value=%s",
		payload.SerialNo, payload.TxHash, payload.Value))

	// Response MUST contain "success" string per PayTheFly API spec
	c.String(http.StatusOK, "success")
}

// processPayTheFlyPayment processes a confirmed payment and credits the user's quota
func processPayTheFlyPayment(db *sql.DB, payload *PayTheFlyWebhookPayload) error {
	var userID int64
	var quota int
	var status string

	err := globals.QueryRowDb(db,
		"SELECT user_id, quota, status FROM paythefly_orders WHERE serial_no = ?",
		payload.SerialNo,
	).Scan(&userID, &quota, &status)

	if err != nil {
		return fmt.Errorf("order not found: %s", payload.SerialNo)
	}

	if strings.ToLower(status) != "pending" {
		return fmt.Errorf("order already processed: %s (status: %s)", payload.SerialNo, status)
	}

	// Update order status
	_, err = globals.ExecDb(db, `
		UPDATE paythefly_orders
		SET status = 'completed', tx_hash = ?, wallet = ?, paid_value = ?, updated_at = NOW()
		WHERE serial_no = ? AND status = 'pending'
	`, payload.TxHash, payload.Wallet, payload.Value, payload.SerialNo)

	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Credit the user's quota
	user := GetUserById(db, userID)
	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	if !user.IncreaseQuota(db, float32(quota)) {
		return fmt.Errorf("failed to increase quota for user: %d", userID)
	}

	return nil
}

// PayTheFlyOrderStatusAPI checks the status of a PayTheFly payment order
func PayTheFlyOrderStatusAPI(c *gin.Context) {
	if !usePayTheFly() {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "PayTheFly payment is not enabled",
		})
		return
	}

	user := GetUserByCtx(c)
	if user == nil {
		return
	}

	serialNo := strings.TrimSpace(c.Query("serial_no"))
	if serialNo == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "serial_no is required",
		})
		return
	}

	db := utils.GetDBFromContext(c)
	var orderStatus string
	var txHash sql.NullString

	err := globals.QueryRowDb(db,
		"SELECT status, tx_hash FROM paythefly_orders WHERE serial_no = ? AND user_id = ?",
		serialNo, user.GetID(db),
	).Scan(&orderStatus, &txHash)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"order_status": orderStatus,
		"tx_hash":      txHash.String,
		"completed":    strings.ToLower(orderStatus) == "completed",
	})
}
