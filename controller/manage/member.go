package manage

import (
	"fmt"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"github.com/pefish/go-redis"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"
	"wallet-storm-wallet/model"

	api_session "github.com/pefish/go-core/api-session"
	uuid "github.com/satori/go.uuid"
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

	memberModel := apiSession.Datas[`memberModel`].(*model.Member)
	if memberModel.TeamId == 0 {
		go_error.Throw(`no team`, constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByEmail(params.Email)
	if targetMember == nil {
		go_error.Throw(`member not found`, constant.MEMBER_NOT_FOUND)
	}
	if targetMember.UserId == 0 {
		go_error.Throw(`member not sync`, constant.MEMBER_NOT_SYNC)
	}
	if targetMember.TeamId != 0 {
		go_error.Throw(`target member is team member already`, constant.ING_TEAM)
	}
	model.MemberModel.UpdateByUserId(targetMember.UserId, map[string]interface{}{
		`team_id`: memberModel.TeamId,
	})
	return map[string]string{}
}

type RemoveMemberParam struct {
	UserId   uint64  `json:"user_id" validate:"required" desc:"user id"`
}

func (this *MemberControllerClass) RemoveMember(apiSession *api_session.ApiSessionClass) interface{} {
	var params RemoveMemberParam
	apiSession.ScanParams(&params)

	memberModel := apiSession.Datas[`memberModel`].(*model.Member)
	if memberModel.TeamId == 0 {
		go_error.Throw(`no team`, constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByUserId(params.UserId)
	if targetMember == nil {
		go_error.Throw(`member not found`, constant.MEMBER_NOT_FOUND)
	}
	if targetMember.TeamId != memberModel.TeamId {
		go_error.Throw(`member not in your team`, constant.USER_NOT_IN_MY_TEAM)
	}
	model.MemberModel.UpdateByUserId(targetMember.UserId, map[string]interface{}{
		`team_id`: 0,
	})
	return map[string]string{}
}

type EditMemberParam struct {
	UserId   uint64  `json:"user_id" validate:"required" desc:"user id"`
	Roles    *string `json:"roles,omitempty" validate:"omitempty" desc:"角色"`
	IsBanned *uint64 `json:"is_banned,omitempty" validate:"omitempty" desc:"是否禁用"`
}

func (this *MemberControllerClass) EditMember(apiSession *api_session.ApiSessionClass) interface{} {
	var params EditMemberParam
	apiSession.ScanParams(&params)

	memberModel := apiSession.Datas[`memberModel`].(*model.Member)
	if memberModel.TeamId == 0 {
		go_error.Throw(`no team`, constant.NO_TEAM_ERROR)
	}

	targetMember := model.MemberModel.GetByUserId(params.UserId)
	if targetMember == nil {
		go_error.Throw(`member not found`, constant.MEMBER_NOT_FOUND)
	}

	if targetMember.TeamId != memberModel.TeamId {
		go_error.Throw(`member not in your team`, constant.USER_NOT_IN_MY_TEAM)
	}

	update := map[string]interface{}{}
	if params.Roles != nil {
		update[`roles`] = params.Roles
	}
	if params.IsBanned != nil {
		update[`is_banned`] = params.IsBanned
	}
	if len(update) == 0 {
		go_error.Throw(`one to edit at least`, constant.PARAM_ERROR)
	}
	model.MemberModel.UpdateByUserId(params.UserId, update)
	return map[string]string{}
}

func (this *MemberControllerClass) ListMember(apiSession *api_session.ApiSessionClass) interface{} {
	memberModel := apiSession.Datas[`memberModel`].(*model.Member)
	if memberModel.TeamId == 0 {
		go_error.Throw(`no team`, constant.NO_TEAM_ERROR)
	}
	var results []model.Member
	model.MemberModel.ListByTeamId(&results, memberModel.TeamId)
	return results
}

func (this *MemberControllerClass) SyncMember(apiSession *api_session.ApiSessionClass) interface{} {
	lockKey := fmt.Sprintf(`storm-syncmember-lock_%d`, apiSession.UserId)
	uniqueId := uuid.NewV1().String()
	if !go_redis.RedisHelper.MustGetLock(lockKey, uniqueId, 4*time.Second) {
		go_error.Throw(`rate limit`, constant.API_RATELIMIT)
	}
	defer go_redis.RedisHelper.MustReleaseLock(lockKey, uniqueId)

	// 检查是否已经同步
	memberModel := model.MemberModel.GetByUserId(apiSession.UserId)
	if memberModel != nil {
		return map[string]string{}
	}
	// 同步
	httpUtil := go_http.NewHttpRequester(go_http.WithTimeout(5 * time.Second))
	resultMap := httpUtil.MustGetForMap(go_http.RequestParam{
		Url: global.AuthServerUrl + `/api/v1/user`,
		Headers: map[string]interface{}{
			`JSON-WEB-TOKEN`: apiSession.Ctx.GetHeader(`JSON-WEB-TOKEN`),
		},
	})
	code := resultMap[`code`].(float64)
	if code != 0 {
		go_error.Throw(resultMap[`msg`].(string), uint64(code))
	}
	data := resultMap[`data`].(map[string]interface{})

	model.MemberModel.Insert(map[string]interface{}{
		`email`:   data[`email`],
		`user_id`: apiSession.UserId,
	})

	return map[string]string{}
}
