package api_strategy

import (
	"errors"
	"github.com/pefish/go-application"
	_type "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/driver/logger"
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

func (this *ApikeyAuthStrategyClass) InitAsync(param interface{}, onAppTerminated chan interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s InitAsync`, this.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s InitAsync defer`, this.GetName())
}

func (this *ApikeyAuthStrategyClass) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

func (this *ApikeyAuthStrategyClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	var p ApikeyAuthParam

	reqPubKey := out.Header(`STM-REQ-KEY`)
	if reqPubKey == `` {
		return go_error.Wrap(errors.New(`auth error. api key not found.`))
	}
	util.UpdateSessionErrorMsg(out, `reqPubKey`, reqPubKey)
	signature := out.Header(`STM-REQ-SIGNATURE`)
	if signature == `` {
		return go_error.Wrap(errors.New(`auth error. signature not found.`))
	}
	timestamp := out.Header(`STM-REQ-TIMESTAMP`)
	if timestamp == `` {
		return go_error.Wrap(errors.New(`auth error. timestamp not found`))
	}
	if !go_application.Application.Debug {
		nowTimestamp := time.Now().UnixNano() / 1e6
		if nowTimestamp-go_reflect.Reflect.MustToInt64(timestamp) > 10*60*1000 {
			return go_error.Wrap(errors.New(`auth expired`))
		}
	}
	requestKeyModel := model.RequestKeyModel.GetByPubKey(reqPubKey)
	if requestKeyModel == nil {
		return go_error.Wrap(errors.New(`auth key error`))
	}
	out.SetUserId(requestKeyModel.UserId)
	util.UpdateSessionErrorMsg(out, `jwtAuth`, requestKeyModel.UserId)
	if param != nil {
		p = param.(ApikeyAuthParam)
		isAllowed := false
		for _, v := range strings.Split(p.AllowedType, `,`) {
			if v == go_reflect.Reflect.ToString(requestKeyModel.Type) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return go_error.Wrap(errors.New(`auth key has not enough right`))
		}
	}
	// 检查用户是否被禁用
	userModel := model.TeamModel.GetByUserIdIsBanned(requestKeyModel.UserId, false)
	if userModel == nil {
		return go_error.Wrap(errors.New(`user is invalid or is baned`))
	}

	if !signature2.VerifySignature(this.structContent(timestamp, out.Method(), out.Path(), out.OriginalParams()), signature, reqPubKey) {
		return go_error.Wrap(errors.New(`auth signature error.`))
	}
	if requestKeyModel.Ip == `` || requestKeyModel.Ip == `*` {
		return nil
	}
	apiIp := out.RemoteAddress()
	for _, ip := range strings.Split(requestKeyModel.Ip, `,`) {
		if ip == apiIp {
			return nil
		}
	}
	return go_error.Wrap(errors.New(`ip is baned`))
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
		sortedStr += k + `=` + go_reflect.Reflect.ToString(params[k]) + `&`
	}
	sortedStr = strings.TrimSuffix(sortedStr, `&`)
	return method + `|` + apiPath + `|` + timestamp + `|` + sortedStr
}
