package main

import (
	"github.com/pefish/go-core/swagger"
	"wallet-storm-wallet/service"
)

func main() {
	service.WalletSvc.Init()
	swagger.GetSwaggerInstance().SetService(&service.WalletSvc).GeneSwagger(`www.zexchange.xyz`, `swagger.json`, `json`)
}
