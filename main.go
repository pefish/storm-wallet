package main

import (
	"fmt"
	"github.com/pefish/go-config"
	"github.com/pefish/go-http"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-mysql"
	"os"
	"runtime/debug"
	"time"
	"wallet-storm-wallet/service"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
		os.Exit(0)
	}()

	go_config.Config.LoadYamlConfig(go_config.Configuration{})

	go_http.Http.SetTimeout(20 * time.Second)

	service.WalletSvc.Init().SetHealthyCheck(nil)

	// 处理日志
	env := go_config.Config.GetString(`env`)
	if env == `local` || env == `dev` {
		loggerInstance := go_logger.Log4goClass{}
		go_logger.Logger.Init(&loggerInstance, service.WalletSvc.GetName(), `debug`)
	} else {
		loggerInstance := go_logger.LogrusClass{}
		go_logger.Logger.Init(&loggerInstance, service.WalletSvc.GetName(), `info`)
	}

	// 初始化数据库连接
	mysqlConfig := go_config.Config.GetMap(`mysql`)
	go_mysql.MysqlHelper.ConnectWithMap(mysqlConfig)
	defer go_mysql.MysqlHelper.Close()

	service.WalletSvc.SetHost(go_config.Config.GetString(`host`))
	service.WalletSvc.SetPort(go_config.Config.GetUint64(`port`))
	service.WalletSvc.Run()
}
