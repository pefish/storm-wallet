package model

import (
	go_mysql "github.com/pefish/go-mysql"
)

var MemberModel = Member{}

type Member struct {
	TeamId   uint64 `json:"team_id"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
	UserId   uint64 `json:"user_id"`
	IsBanned int64  `json:"is_banned"`
	BaseModel
}

func (this *Member) GetTableName() string {
	return `member`
}

func (this *Member) MustAddMember(userId uint64, email string) {
	go_mysql.MysqlHelper.MustAffectedInsert(this.GetTableName(), map[string]interface{}{
		`user_id`: userId,
		`email`:   email,
		`role`:    3,
	})
}

func (this *Member) GetByMemberId(memberId uint64) *Member {
	result := Member{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`id`: memberId,
	}); notFound {
		return nil
	}
	return &result
}

func (this *Member) ListByTeamId(results *[]Member, teamId uint64) {
	go_mysql.MysqlHelper.MustSelect(results, this.GetTableName(), `*`, map[string]interface{}{
		`team_id`: teamId,
	})
}

func (this *Member) GetByUserId(userId uint64) *Member {
	result := Member{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`: userId,
	}); notFound {
		return nil
	}
	return &result
}

func (this *Member) GetByEmail(email string) *Member {
	result := Member{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`email`: email,
	}); notFound {
		return nil
	}
	return &result
}

func (this *Member) GetValidByUserId(userId uint64) *Member {
	result := Member{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:   userId,
		`is_banned`: 0,
	}); notFound {
		return nil
	}
	return &result
}

func (this *Member) UpdateByUserIdTeamId(userId uint64, teamId uint64, update map[string]interface{}) {
	go_mysql.MysqlHelper.MustAffectedUpdate(this.GetTableName(), update, map[string]interface{}{
		`user_id`: userId,
		`team_id`: teamId,
	})
}

func (this *Member) UpdateByUserId(userId uint64, update map[string]interface{}) {
	go_mysql.MysqlHelper.MustAffectedUpdate(this.GetTableName(), update, map[string]interface{}{
		`user_id`: userId,
	})
}

func (this *Member) Insert(member map[string]interface{}) {
	go_mysql.MysqlHelper.MustInsert(this.GetTableName(), member)
}
