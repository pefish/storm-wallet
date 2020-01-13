package team_admin

import (
	api_strategy "wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"

	api_channel_builder "github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
)

var MemberRoute = map[string]*api_channel_builder.Route{
	`add_member`: {
		Description: "新增成员",
		Path:        "/add-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.FetchMemberIdStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.AddMemberParam{
			Email: `laijiyong@qq.com`,
		},
		Controller: controller.MemberController.AddMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
	`edit_member`: {
		Description: "编辑成员",
		Path:        "/edit-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.FetchMemberIdStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.EditMemberParam{
			MemberId: 534,
		},
		Controller: controller.MemberController.EditMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
}
