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
	//memberModel := model.MemberModel.GetByMemberId(apiSession.UserId)
	//if memberModel == nil {
	//	go_error.Throw(`member not found`, constant.ILLEGAL_MEMBER)
	//}
	//model.MemberModel.MustAddMember(memberModel.TeamId, params.Email)
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

	//memberModel := model.MemberModel.GetByMemberId(params.MemberId)
	//if memberModel == nil {
	//	go_error.Throw(`member not found`, constant.ILLEGAL_MEMBER)
	//}
	//update := map[string]interface{}{}
	//if params.Email != nil {
	//	update[`email`] = params.Email
	//}
	//if params.Role != nil {
	//	update[`role`] = params.Role
	//}
	//if params.Password != nil {
	//	update[`password`] = params.Password
	//}
	//if params.IsBanned != nil {
	//	update[`is_banned`] = params.IsBanned
	//}
	//if len(update) == 0 {
	//	go_error.Throw(`params error`, constant.PARAM_ERROR)
	//}
	//model.MemberModel.UpdateByMap(params.MemberId, update)
	return map[string]string{}
}

func (this *MemberControllerClass) ListMember(apiSession *api_session.ApiSessionClass) interface{} {
	// TODO
	return map[string]string{}
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
