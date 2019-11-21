package api_strategy

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	signature2 "github.com/pefish/storm-golang-sdk/signature"
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
	return `request_auth`
}

func (this *ApikeyAuthStrategyClass) GetDescription() string {
	return `对请求的签名进行校验`
}

func (this *ApikeyAuthStrategyClass) GetErrorCode() uint64 {
	return constant.AUTH_ERROR
}

type ApikeyAuthParam struct {
	AllowedType string // 允许的key类型，逗号隔开
}

func (this *ApikeyAuthStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	var p ApikeyAuthParam

	reqPubKey := out.Ctx.GetHeader(`STM-REQ-KEY`)
	if reqPubKey == `` {
		go_error.ThrowInternal(`auth error. api key not found.`)
	}
	util.UpdateCtxValuesErrorMsg(out.Ctx, `reqPubKey`, reqPubKey)
	signature := out.Ctx.GetHeader(`STM-REQ-SIGNATURE`)
	if signature == `` {
		go_error.ThrowInternal(`auth error. signature not found.`)
	}
	timestamp := out.Ctx.GetHeader(`STM-REQ-TIMESTAMP`)
	if timestamp == `` {
		go_error.ThrowInternal(`auth error. timestamp not found`)
	}
	if !go_application.Application.Debug {
		nowTimestamp := time.Now().UnixNano() / 1e6
		if nowTimestamp-go_reflect.Reflect.MustToInt64(timestamp) > 10*60*1000 {
			go_error.ThrowInternal(`auth expired`)
		}
	}
	requestKeyModel := model.RequestKeyModel.GetByPubKey(reqPubKey)
	if requestKeyModel == nil {
		go_error.ThrowInternal(`auth key error`)
	}
	out.UserId = requestKeyModel.UserId
	util.UpdateCtxValuesErrorMsg(out.Ctx, `jwtAuth`, requestKeyModel.UserId)
	if param != nil {
		p = param.(ApikeyAuthParam)
		isAllowed := false
		for _, v := range strings.Split(p.AllowedType, `,`) {
			if v == go_reflect.Reflect.MustToString(requestKeyModel.Type) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			go_error.ThrowInternalWithInternalMsg(`auth key has not enough right`, fmt.Sprintf(`AllowedType: %s`, p.AllowedType))
		}
	}
	// 检查用户是否被禁用
	userModel := model.UserModel.GetByUserIdIsBanned(requestKeyModel.UserId, false)
	if userModel == nil {
		go_error.ThrowInternal(`user is invalid or is baned`)
	}

	if !signature2.VerifySignature(this.structContent(timestamp, out.Ctx.Method(), out.Ctx.Path(), out.Params), signature, reqPubKey) {
		go_error.ThrowInternalWithInternalMsg(`auth signature error.`, fmt.Sprintf(`signature: %s`, signature))
	}
	if requestKeyModel.Ip == `` || requestKeyModel.Ip == `*` {
		return
	}
	apiIp := out.Ctx.RemoteAddr()
	for _, ip := range strings.Split(requestKeyModel.Ip, `,`) {
		if ip == apiIp {
			return
		}
	}
	go_error.ThrowInternalWithInternalMsg(`ip is baned`, fmt.Sprintf(`ip: %s, expected ip: %s`, apiIp, requestKeyModel.Ip))
}

func (this *ApikeyAuthStrategyClass) structContent(timestamp string, method string, apiPath string, params map[string]interface{}) string {
	sortedStr := ``
	var keys []string
	for k, v := range params {
		if v != nil { // nil参数不参与签名
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		sortedStr += k + `=` + go_reflect.Reflect.MustToString(params[k]) + `&`
	}
	sortedStr = strings.TrimSuffix(sortedStr, `&`)
	return method + `|` + apiPath + `|` + timestamp + `|` + sortedStr
}
