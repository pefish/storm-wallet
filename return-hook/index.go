package return_hook

import (
	"encoding/json"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-crypto"
	"github.com/pefish/go-reflect"
	"time"
)

func ReturnHook(apiContext *api_session.ApiSessionClass, apiResult *api_channel_builder.ApiResult) interface{} {
	bytes, err := json.Marshal(apiResult)
	if err != nil {
		panic(err)
	}
	timestamp := go_reflect.Reflect.MustToString(time.Now().UnixNano() / 1e6)
	apiContext.Ctx.Header(`BIZ_TIMESTAMP`, timestamp)
	apiKey := apiContext.Datas[`apiKey`].(string)
	apiSecret := apiContext.Datas[`apiSecret`].(string)
	apiContext.Ctx.Header(`BIZ_RESP_SIGNATURE`, go_crypto.Crypto.HmacSha256ToHex(string(bytes) + `|` + timestamp + `|` + apiKey, apiSecret))
	return apiResult
}
