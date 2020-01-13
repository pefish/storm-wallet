package controller

import (
	"github.com/pefish/go-error"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/model"

	api_session "github.com/pefish/go-core/api-session"
)

type MemberControllerClass struct {
}

var MemberController = MemberControllerClass{}

type AddMemberParam struct {
	Email string `json:"email" validate:"required" desc:"邮箱"`
}

func (this *MemberControllerClass) AddMember(apiSession *api_session.ApiSessionClass) interface{} {
	var params AddMemberParam
	apiSession.ScanParams(&params)
	memberModel := model.MemberModel.GetByMemberId(apiSession.UserId)
	if memberModel == nil {
		go_error.Throw(`member not found`, constant.ILLEGAL_MEMBER)
	}
	model.MemberModel.MustAddMember(memberModel.TeamId, params.Email)
	return map[string]string{}
}

type EditMemberParam struct {
	MemberId uint64  `json:"member_id" validate:"required" desc:"id"`
	Email    *string `json:"email,omitempty" validate:"omitempty" desc:"邮箱"`
	Role     *uint64 `json:"role,omitempty" validate:"omitempty" desc:"角色"`
	Password *string `json:"password,omitempty" validate:"omitempty" desc:"密码"`
	IsBanned *uint64 `json:"is_banned,omitempty" validate:"omitempty" desc:"是否禁用"`
}

func (this *MemberControllerClass) EditMember(apiSession *api_session.ApiSessionClass) interface{} {
	var params EditMemberParam
	apiSession.ScanParams(&params)

	memberModel := model.MemberModel.GetByMemberId(params.MemberId)
	if memberModel == nil {
		go_error.Throw(`member not found`, constant.ILLEGAL_MEMBER)
	}
	update := map[string]interface{}{}
	if params.Email != nil {
		update[`email`] = params.Email
	}
	if params.Role != nil {
		update[`role`] = params.Role
	}
	if params.Password != nil {
		update[`password`] = params.Password
	}
	if params.IsBanned != nil {
		update[`is_banned`] = params.IsBanned
	}
	if len(update) == 0 {
		go_error.Throw(`params error`, constant.PARAM_ERROR)
	}
	model.MemberModel.UpdateByMap(params.MemberId, update)
	return map[string]string{}
}
