package main

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-http"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-mysql"
	"github.com/pefish/go-redis"
	"os"
	"runtime/debug"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/route"
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

	go_config.Config.MustLoadYamlConfig(go_config.Configuration{
		ConfigEnvName: `GO_CONFIG`,
		SecretEnvName: `GO_SECRET`,
	})

	go_http.Http = go_http.NewHttpRequester(go_http.WithTimeout(20 * time.Second))


	// 处理日志
	env := go_config.Config.GetString(`env`)
	go_application.Application.Debug = env == `local` || env == `dev`
	go_logger.Logger.Init(service.Service.GetName(), ``)
	logger.LoggerDriver.Register(go_logger.Logger)

	// 初始化数据库连接
	mysqlConfig := go_config.Config.MustGetMap(`mysql`)
	go_mysql.MysqlHelper.SetTagName(`json`)
	go_mysql.MysqlHelper.ConnectWithMap(mysqlConfig)
	defer go_mysql.MysqlHelper.Close()

	// 初始化redis连接
	redisConfig := go_config.Config.MustGetMap(`redis`)
	go_redis.RedisHelper.ConnectWithMap(redisConfig)
	defer go_redis.RedisHelper.Close()

	service.Service.SetName(`storm钱包服务api`)
	service.Service.SetPath(`/api/storm`)
	api_strategy.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	service.Service.SetRoutes(route.AddressRoute, route.TransactionRoute, route.WithdrawRoute, route.UserRoute)
	service.Service.SetHost(go_config.Config.GetString(`host`))
	service.Service.SetPort(go_config.Config.GetUint64(`port`))
	service.Service.Run()
}
