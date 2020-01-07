package controller

import (
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
	// TODO
	return map[string]string{}
}
