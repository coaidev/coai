# PayTheFly Crypto Payment Integration

PayTheFly enables cryptocurrency payments on BSC (BNB Smart Chain) and TRON networks, allowing users to purchase quota with crypto.

## Configuration

Add the following to your `config.yaml`:

```yaml
paythefly:
  enabled: true
  project_id: "your-project-id"
  project_key: "your-project-key"
  private_key: "your-private-key"
  chain_id: 56
  token_address: "0x..."
```

> **Security**: Store `project_key` and `private_key` as environment variables in production.

## Database Migration

```bash
mysql -u root -p chatnio < migration/paythefly.sql
```

## API Endpoints

### Create Payment Order
`POST /api/paythefly/create` — Creates order, returns serial_no and payment parameters.

### Webhook
`POST /api/paythefly/webhook` — Receives confirmed payment notifications from PayTheFly.

### Check Order Status
`GET /api/paythefly/status?serial_no=xxx` — Polls order completion.

## Payment Flow

1. User creates order → 2. Redirect to PayTheFly → 3. User pays crypto → 4. Webhook confirms → 5. Quota credited

## PayTheFly API Spec

- **EIP-712 Domain**: `{ name: "PayTheFlyPro", version: "1" }`
- **Amount**: Human-readable (e.g., `"0.01"`), NOT raw token units
- **Webhook signature**: `HMAC-SHA256(data + "." + timestamp, projectKey)`
- **Webhook response**: Must contain `"success"` string
- **Payload fields**: `value`, `confirmed`, `serial_no`, `tx_hash`, `wallet`, `tx_type`
- **tx_type**: 1=payment, 2=withdrawal
- **Chains**: BSC (chainId=56, 18 decimals), TRON (chainId=728126428, 6 decimals)

For more info: [https://pro.paythefly.com](https://pro.paythefly.com)
