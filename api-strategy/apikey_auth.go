package api_strategy

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-crypto"
	"github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	"sort"
	"strings"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/model"
)

type ApikeyAuthStrategyClass struct {
}

var ApikeyAuthStrategy = ApikeyAuthStrategyClass{}

func (this *ApikeyAuthStrategyClass) GetName() string {
	return `apikey_auth`
}

func (this *ApikeyAuthStrategyClass) GetDescription() string {
	return `对apikey以及签名进行校验`
}

func (this *ApikeyAuthStrategyClass) GetErrorCode() uint64 {
	return constant.AUTH_ERROR
}

type ApikeyAuthParam struct {
	AllowedType string // 允许的api key类型，逗号隔开
}

func (this *ApikeyAuthStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	var p ApikeyAuthParam

	apiKey := out.Ctx.GetHeader(`BIZ-API-KEY`)
	if apiKey == `` {
		go_error.ThrowInternal(`auth error. api key not found.`)
	}
	util.UpdateCtxValuesErrorMsg(out.Ctx, `apiKey`, apiKey)
	signature := out.Ctx.GetHeader(`BIZ-API-SIGNATURE`)
	if signature == `` {
		go_error.ThrowInternal(`auth error. signature not found.`)
	}
	timestamp := out.Ctx.GetHeader(`BIZ-API-TIMESTAMP`)
	if timestamp == `` {
		go_error.ThrowInternal(`auth error. timestamp not found`)
	}
	nowTimestamp := time.Now().UnixNano() / 1e6
	if nowTimestamp-go_reflect.Reflect.ToInt64(timestamp) > 10*60*1000 {
		go_error.ThrowInternal(`auth expired`)
	}
	apiKeyModel := model.ApiKeyModel.GetByApiKey(apiKey)
	if apiKeyModel == nil {
		go_error.ThrowInternal(`auth key error`)
	}
	out.UserId = apiKeyModel.UserId
	util.UpdateCtxValuesErrorMsg(out.Ctx, `jwtAuth`, apiKeyModel.UserId)
	if param != nil {
		p = param.(ApikeyAuthParam)
		isAllowed := false
		for _, v := range strings.Split(p.AllowedType, `,`) {
			if v == go_reflect.Reflect.ToString(apiKeyModel.Type) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			go_error.ThrowInternal(`auth key has not enough right`)
		}
	}
	// 检查用户是否被禁用
	userModel := model.UserModel.GetByUserIdIsBanned(apiKeyModel.UserId, false)
	if userModel == nil {
		go_error.ThrowInternal(`user is invalid or is baned`)
	}
	realSignature := this.sign(apiKeyModel.ApiSecret, timestamp, out.Ctx.Method(), out.Ctx.Path(), out.Params)
	if realSignature != signature {
		go_error.ThrowInternal(`auth signature error.`)
	}
	if apiKeyModel.Ip == `` {
		return
	}
	apiIp := out.Ctx.RemoteAddr()
	for _, ip := range strings.Split(apiKeyModel.Ip, `,`) {
		if ip == apiIp {
			return
		}
	}
	go_error.ThrowInternal(`ip is baned`)
}

func (this *ApikeyAuthStrategyClass) sign(secret string, timestamp string, method string, apiPath string, params map[string]interface{}) string {
	sortedStr := ``
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sortedStr += k + `=` + go_reflect.Reflect.ToString(params[k]) + `&`
	}
	sortedStr = strings.TrimSuffix(sortedStr, `&`)
	toSignStr := method + `|` + apiPath + `|` + timestamp + `|` + sortedStr
	return go_crypto.Crypto.HmacSha256ToHex(toSignStr, secret)
}
