package api_strategy

import (
	"wallet-storm-wallet/constant"

	api_channel_builder "github.com/pefish/go-core/api-channel-builder"
	api_session "github.com/pefish/go-core/api-session"
	go_error "github.com/pefish/go-error"
	go_reflect "github.com/pefish/go-reflect"
)

type FetchMemberIdClass struct {
}

var FetchMemberIdStrategy = FetchMemberIdClass{}

func (this *FetchMemberIdClass) GetName() string {
	return `fetch_oauth_userId`
}

func (this *FetchMemberIdClass) GetDescription() string {
	return `提取oauth请求的成员id`
}

func (this *FetchMemberIdClass) GetErrorCode() uint64 {
	return constant.AUTH_ERROR
}

func (this *FetchMemberIdClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	memberId := out.Ctx.GetHeader(`MEMBER_ID`)
	if memberId == `` {
		go_error.ThrowInternal(`auth error. member id not found.`)
	}
	out.UserId = go_reflect.Reflect.MustToUint64(memberId)
}
