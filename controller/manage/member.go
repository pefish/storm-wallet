package manage

import (
	"errors"
	"fmt"
	"github.com/pefish/go-core/api"
	_type "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"github.com/pefish/go-redis"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"
	"wallet-storm-wallet/model"

	uuid "github.com/satori/go.uuid"
)

type MemberControllerClass struct {
}

var MemberController = MemberControllerClass{}

type AddMemberParam struct {
	Name string `json:"name" validate:"required" desc:"名称"`
	Email string `json:"email" validate:"required" desc:"邮箱"`
}

func (this *MemberControllerClass) AddMember(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	var params AddMemberParam
	apiSession.ScanParams(&params)

	memberModel := apiSession.Data(`memberModel`).(*model.Member)
	if memberModel.TeamId == 0 {
		return nil, go_error.WrapWithCode(errors.New(`no team`), constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByEmail(params.Email)
	if targetMember == nil {
		return nil, go_error.WrapWithCode(errors.New(`member not found`), constant.MEMBER_NOT_FOUND)
	}
	if targetMember.UserId == 0 {
		return nil, go_error.WrapWithCode(errors.New(`member not sync`), constant.MEMBER_NOT_SYNC)
	}
	if targetMember.TeamId != 0 {
		return nil, go_error.WrapWithCode(errors.New(`target member is team member already`), constant.ING_TEAM)
	}
	model.MemberModel.UpdateByUserId(targetMember.UserId, map[string]interface{}{
		`team_id`: memberModel.TeamId,
		`name`: params.Name,
	})
	return map[string]string{}, nil
}

type RemoveMemberParam struct {
	UserId   uint64  `json:"user_id" validate:"required" desc:"user id"`
}

func (this *MemberControllerClass) RemoveMember(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	var params RemoveMemberParam
	apiSession.ScanParams(&params)

	memberModel := apiSession.Data(`memberModel`).(*model.Member)
	if memberModel.TeamId == 0 {
		return nil, go_error.WrapWithCode(errors.New(`no team`), constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByUserId(params.UserId)
	if targetMember == nil {
		return nil, go_error.WrapWithCode(errors.New(`member not found`), constant.MEMBER_NOT_FOUND)
	}
	if targetMember.TeamId != memberModel.TeamId {
		return nil, go_error.WrapWithCode(errors.New(`member not in your team`), constant.USER_NOT_IN_MY_TEAM)
	}
	model.MemberModel.UpdateByUserId(targetMember.UserId, map[string]interface{}{
		`team_id`: 0,
	})
	return map[string]string{}, nil
}

type EditMemberParam struct {
	UserId   uint64  `json:"user_id" validate:"required" desc:"user id"`
	Roles    *string `json:"roles,omitempty" validate:"omitempty" desc:"角色"`
	IsBanned *uint64 `json:"is_banned,omitempty" validate:"omitempty" desc:"是否禁用"`
}

func (this *MemberControllerClass) EditMember(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	var params EditMemberParam
	apiSession.ScanParams(&params)

	memberModel := apiSession.Data(`memberModel`).(*model.Member)
	if memberModel.TeamId == 0 {
		return nil, go_error.WrapWithCode(errors.New(`no team`), constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByUserId(params.UserId)
	if targetMember == nil {
		return nil, go_error.WrapWithCode(errors.New(`member not found`), constant.MEMBER_NOT_FOUND)
	}

	if targetMember.TeamId != memberModel.TeamId {
		return nil, go_error.WrapWithCode(errors.New(`member not in your team`), constant.USER_NOT_IN_MY_TEAM)
	}

	update := map[string]interface{}{}
	if params.Roles != nil {
		update[`roles`] = params.Roles
	}
	if params.IsBanned != nil {
		update[`is_banned`] = params.IsBanned
	}
	if len(update) == 0 {
		return nil, go_error.WrapWithCode(errors.New(`one to edit at least`), constant.PARAM_ERROR)
	}
	model.MemberModel.UpdateByUserId(params.UserId, update)
	return map[string]string{}, nil
}

func (this *MemberControllerClass) ListMember(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	memberModel := apiSession.Data(`memberModel`).(*model.Member)
	if memberModel.TeamId == 0 {
		return nil, go_error.WrapWithCode(errors.New(`no team`), constant.NO_TEAM_ERROR)
	}
	var results []model.Member
	model.MemberModel.ListByTeamId(&results, memberModel.TeamId)
	return results, nil
}

func (this *MemberControllerClass) SyncMember(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	lockKey := fmt.Sprintf(`storm-syncmember-lock_%d`, apiSession.UserId)
	uniqueId := uuid.NewV1().String()
	if !go_redis.RedisHelper.MustGetLock(lockKey, uniqueId, 4*time.Second) {
		return nil, go_error.WrapWithCode(errors.New(`rate limit`), constant.API_RATELIMIT)
	}
	defer go_redis.RedisHelper.MustReleaseLock(lockKey, uniqueId)

	// 检查是否已经同步
	memberModel := model.MemberModel.GetByUserId(apiSession.UserId())
	if memberModel != nil {
		return map[string]string{}, nil
	}
	// 同步
	httpUtil := go_http.NewHttpRequester(go_http.WithTimeout(5 * time.Second))
	var apiResult api.ApiResult
	_, err := httpUtil.GetForStruct(go_http.RequestParam{
		Url: global.AuthServerUrl + `/api/v1/user`,
		Headers: map[string]interface{}{
			`JSON-WEB-TOKEN`: apiSession.Header(`JSON-WEB-TOKEN`),
		},
	}, &apiResult)
	if err != nil {
		return nil, go_error.Wrap(err)
	}
	if apiResult.Code != 0 {
		return nil, go_error.WrapWithCode(errors.New(apiResult.Msg), apiResult.Code)
	}
	data := apiResult.Data.(map[string]interface{})

	model.MemberModel.Insert(map[string]interface{}{
		`email`:   data[`email`],
		`user_id`: apiSession.UserId,
	})

	return map[string]string{}, nil
}
