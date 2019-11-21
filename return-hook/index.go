package return_hook

import (
	"encoding/json"
	"github.com/kataras/iris/core/errors"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-reflect"
	"github.com/pefish/storm-golang-sdk/signature"
	"time"
	"wallet-storm-wallet/model"
)

func ReturnHook(apiContext *api_session.ApiSessionClass, apiResult *api_channel_builder.ApiResult) (interface{}, error) {
	bytes, err := json.Marshal(apiResult)
	if err != nil {
		return nil, errors.New(apiResult.InternalMsg + ` - ` + err.Error())
	}
	timestamp := go_reflect.Reflect.MustToString(time.Now().UnixNano() / 1e6)
	apiContext.Ctx.Header(`STM-RES-TIMESTAMP`, timestamp)
	responseKeyModel := model.ResponseKeyModel.GetByUserId(apiContext.UserId)
	if responseKeyModel == nil {
		return nil, errors.New(apiResult.InternalMsg + ` - user do not have response keys.`)
	}
	apiContext.Ctx.Header(`STM-RES-SIGNATURE`, signature.SignMessage(string(bytes)+`|`+timestamp, responseKeyModel.PrivateKey))
	return apiResult, nil
}
