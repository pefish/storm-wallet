package return_hook

import (
	"encoding/json"
	"errors"
	"github.com/pefish/go-core/api"
	_type "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	"github.com/pefish/storm-golang-sdk/signature"
	"time"
	"wallet-storm-wallet/model"
)

func ReturnHook(apiContext _type.IApiSession, apiResult *api.ApiResult) (interface{}, *go_error.ErrorInfo) {
	bytes, err := json.Marshal(apiResult)
	if err != nil {
		return nil, go_error.Wrap(err)
	}
	timestamp := go_reflect.Reflect.ToString(time.Now().UnixNano() / 1e6)
	apiContext.SetHeader(`STM-RES-TIMESTAMP`, timestamp)
	responseKeyModel := model.ResponseKeyModel.GetByUserId(apiContext.UserId())
	if responseKeyModel == nil {
		return nil, go_error.Wrap(errors.New(`user do not have response keys.`))
	}
	apiContext.SetHeader(`STM-RES-SIGNATURE`, signature.SignMessage(string(bytes)+`|`+timestamp, responseKeyModel.PrivateKey))
	return apiResult, nil
}
