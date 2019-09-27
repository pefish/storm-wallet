package api_strategy

import (
	"fmt"
	"github.com/pefish/go-reflect"
	"testing"
	"time"
)

func TestApikeyAuthStrategyClass_sign(t *testing.T) {
	timestamp := go_reflect.Reflect.ToString(time.Now().UnixNano() / 1e6)
	sig := ApikeyAuthStrategy.sign(`tw2456245twe2`, timestamp, `POST`, `/api/storm-wallet/v1/withdraw`, map[string]interface{}{
		"currency": "BTC",
		"chain": "Btc",
		"request_id": "5",
		"address": "1HL2bVtN673KVpELnY54uYdCrSJD9Kv1Ta",
		"amount": "89",
		"memo": "63562",
	})
	fmt.Println(timestamp, sig)

	sig1 := ApikeyAuthStrategy.sign(`tw2456245twe2`, timestamp, `GET`, `/api/storm-wallet/v1/balance`, map[string]interface{}{

	})
	fmt.Println(timestamp, sig1)

	sig2 := ApikeyAuthStrategy.sign(`tw2456245twe2`, timestamp, `GET`, `/api/storm-wallet/v1/validate-address`, map[string]interface{}{
		`currency`: `ETH`,
		`chain`: `Eth`,
		`address`: `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a1`,
	})
	fmt.Println(timestamp, sig2)
}
