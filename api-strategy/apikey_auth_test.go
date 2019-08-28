package api_strategy

import (
	"fmt"
	"testing"
)

func TestApikeyAuthStrategyClass_sign(t *testing.T) {
	sig := ApikeyAuthStrategy.sign(`tw2456245twe2`, `1566987644000`, `POST`, `/api/wallet-storm/v1/new-address`, map[string]interface{}{
		`currency`: "ETH",
	})
	fmt.Println(sig)
}
