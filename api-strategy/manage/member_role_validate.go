package manage

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-error"
	"strings"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/model"
)

type MemberRoleValidateStrategyClass struct {
}

var MemberRoleValidateStrategy = MemberRoleValidateStrategyClass{}

func (this *MemberRoleValidateStrategyClass) GetName() string {
	return `member_role_validate`
}

func (this *MemberRoleValidateStrategyClass) GetDescription() string {
	return `校验用户的角色`
}

func (this *MemberRoleValidateStrategyClass) GetErrorCode() uint64 {
	return constant.ROLE_AUTH_ERROR
}

type MemberRoleValidateParam struct {
	RequiredRole string
}

func (this *MemberRoleValidateStrategyClass) InitAsync(param interface{}, onAppTerminated chan interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s InitAsync`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s InitAsync defer`, this.GetName())
}

func (this *MemberRoleValidateStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

func (this *MemberRoleValidateStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	memberModel := model.MemberModel.GetValidByUserId(out.UserId)
	if memberModel == nil {
		go_error.Throw(`user not found or banned`, constant.MEMBER_NOT_FOUND)
	}
	if param != nil {
		newParam := param.(MemberRoleValidateParam)
		if !strings.Contains(memberModel.Roles, newParam.RequiredRole) {
			go_error.ThrowInternal(`required scope: ` + newParam.RequiredRole)
		}
	}
	out.Datas[`memberModel`] = memberModel
}
