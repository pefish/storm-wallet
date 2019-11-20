package main

import (
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-core/swagger"
	"wallet-storm-wallet/route"
)

func main() {
	service.Service.SetName(`storm钱包服务api`)
	service.Service.SetPath(`/api/storm-wallet`)
	service.Service.SetRoutes(route.AddressRoute, route.TransactionRoute, route.WithdrawRoute, route.UserRoute)
	swagger.GetSwaggerInstance().GeneSwagger(`www.zgtest.club`, `swagger.json`, `json`)
}
