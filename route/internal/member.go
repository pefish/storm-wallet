package route

import (
	"wallet-storm-wallet/controller"

	api_channel_builder "github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
)

var MemberRoute = map[string]*api_channel_builder.Route{
	`add_member`: {
		Description: "新增成员",
		Path:        "/add-member",
		Method:      "POST",
		Strategies:  []api_channel_builder.StrategyRoute{},
		ParamType:   api_strategy2.ALL_TYPE,
		Params: controller.AddMemberParam{
			Email: `laijiyong@qq.com`,
		},
		Controller: controller.MemberController.AddMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
}
