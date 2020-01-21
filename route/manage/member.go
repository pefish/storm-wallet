package manage

import (
	"wallet-storm-wallet/api-strategy/manage"
	manage2 "wallet-storm-wallet/controller/manage"

	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
)

var MemberRoute = map[string]*api_channel_builder.Route{
	`add_member`: {
		Description: "新增成员(必须先同步用户且不在任何团队中)",
		Path:        "/v1/add-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &manage.OauthJwtValidateStrategy,
				Param: manage.OauthJwtValidateParam{
					RequiredScopes: []string{
						`storm_partner`,
					},
				},
			},
			{
				Strategy: &manage.MemberRoleValidateStrategy,
				Param: manage.MemberRoleValidateParam{
					RequiredRole: `team_admin`,
				},
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: manage2.AddMemberParam{
			Email: `laijiyong@qq.com`,
		},
		Controller: manage2.MemberController.AddMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
	`remove_member`: {
		Description: "从团队中移除成员",
		Path:        "/v1/remove-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &manage.OauthJwtValidateStrategy,
				Param: manage.OauthJwtValidateParam{
					RequiredScopes: []string{
						`storm_partner`,
					},
				},
			},
			{
				Strategy: &manage.MemberRoleValidateStrategy,
				Param: manage.MemberRoleValidateParam{
					RequiredRole: `team_admin`,
				},
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: manage2.RemoveMemberParam{
			UserId: 1,
		},
		Controller: manage2.MemberController.RemoveMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
	`edit_member`: {
		Description: "编辑成员",
		Path:        "/v1/edit-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &manage.OauthJwtValidateStrategy,
				Param: manage.OauthJwtValidateParam{
					RequiredScopes: []string{
						`storm_partner`,
					},
				},
			},
			{
				Strategy: &manage.MemberRoleValidateStrategy,
				Param: manage.MemberRoleValidateParam{
					RequiredRole: `team_admin`,
				},
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: manage2.EditMemberParam{
			UserId: 534,
		},
		Controller: manage2.MemberController.EditMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
	`list_member`: {
		Description: "列出成员",
		Path:        "/v1/list-member",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &manage.OauthJwtValidateStrategy,
				Param: manage.OauthJwtValidateParam{
					RequiredScopes: []string{
						`storm_partner`,
					},
				},
			},
			{
				Strategy: &manage.MemberRoleValidateStrategy,
				Param: manage.MemberRoleValidateParam{
					RequiredRole: `team_admin`,
				},
			},
		},
		ParamType:  api_strategy2.ALL_TYPE,
		Controller: manage2.MemberController.ListMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
	`sync_member`: {
		Description: "从授权服务器同步成员信息(前端控制注册或登陆成功之后才调一次)",
		Path:        "/v1/sync-member",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &manage.OauthJwtValidateStrategy,
			},
		},
		ParamType:  api_strategy2.ALL_TYPE,
		Controller: manage2.MemberController.SyncMember,
		Return: api_channel_builder.ApiResult{
			Data: map[string]interface{}{},
		},
	},
}
