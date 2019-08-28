package service

import (
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/service/route"
)

type WalletSvcClass struct {
	service.BaseServiceClass
}

var WalletSvc = WalletSvcClass{}

func (this *WalletSvcClass) Init(opts ...interface{}) service.InterfaceService {
	this.SetName(`storm钱包服务api`)
	this.SetPath(`/api/wallet-storm`)
	api_strategy.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	this.SetRoutes(route.UserRoute)
	return this
}
