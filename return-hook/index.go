package return_hook

import (
	"encoding/json"
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	"github.com/pefish/storm-golang-sdk/signature"
	"time"
	"wallet-storm-wallet/model"
)

func ReturnHook(apiContext *api_session.ApiSessionClass, apiResult *api.ApiResult) (interface{}, *go_error.ErrorInfo) {
	bytes, err := json.Marshal(apiResult)
	if err != nil {
		return nil, &go_error.ErrorInfo{
			ErrorMessage: apiResult.InternalMsg + ` - ` + err.Error(),
			InternalErrorMessage: apiResult.InternalMsg,
			ErrorCode: go_error.INTERNAL_ERROR_CODE,
		}
	}
	timestamp := go_reflect.Reflect.MustToString(time.Now().UnixNano() / 1e6)
	apiContext.Ctx.Header(`STM-RES-TIMESTAMP`, timestamp)
	responseKeyModel := model.ResponseKeyModel.GetByUserId(apiContext.UserId)
	if responseKeyModel == nil {
		return nil, &go_error.ErrorInfo{
			ErrorMessage: `user do not have response keys.`,
			InternalErrorMessage: apiResult.InternalMsg,
			ErrorCode: go_error.INTERNAL_ERROR_CODE,
		}
	}
	apiContext.Ctx.Header(`STM-RES-SIGNATURE`, signature.SignMessage(string(bytes)+`|`+timestamp, responseKeyModel.PrivateKey))
	return apiResult, nil
}
