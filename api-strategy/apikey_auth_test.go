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
		"currency": "ETH",
		"chain": "Eth",
		"request_id": "1",
		"address": "0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a",
		"amount": "4",
		"memo": "63562",
	})
	fmt.Println(timestamp, sig)
}
