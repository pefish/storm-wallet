package main

import (
	"fmt"
	external_service "github.com/pefish/go-core/driver/external-service"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy2 "github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
	"os"
	"runtime/debug"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/route"
	"wallet-storm-wallet/route/manage"

	go_application "github.com/pefish/go-application"
	go_config "github.com/pefish/go-config"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	go_redis "github.com/pefish/go-redis"
	external_service2 "wallet-storm-wallet/external-service"
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

	// 处理日志
	env := go_config.Config.MustGetString(`env`)
	go_application.Application.SetEnv(env)
	go_logger.Logger = go_logger.NewLogger(go_config.Config.MustGetStringDefault(`logLevel`, `debug`), go_logger.WithPrefix(service.Service.GetName()))
	logger.LoggerDriver.Register(go_logger.Logger)

	// 初始化数据库连接
	mysqlConfig := go_config.Config.MustGetMap(`mysql`)
	go_mysql.MysqlInstance.SetLogger(go_logger.Logger)
	go_mysql.MysqlInstance.MustConnectWithMap(mysqlConfig)
	defer go_mysql.MysqlInstance.Close()

	// 初始化redis连接
	redisConfig := go_config.Config.MustGetMap(`redis`)
	go_redis.RedisHelper.SetLogger(go_logger.Logger)
	go_redis.RedisHelper.MustConnectWithMap(redisConfig)
	defer go_redis.RedisHelper.Close()

	service.Service.SetName(`storm钱包服务api`)
	service.Service.SetPath(`/api/storm`)
	global_api_strategy2.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	//global.HydraClientInstance = client.New(
	//	httptransport.New(
	//		go_config.Config.MustGetString(`/authServer/host`),
	//		go_config.Config.MustGetString(`/authServer/basePath`),
	//		[]string{go_config.Config.MustGetString(`/authServer/scheme`)}),
	//	nil)
	//global.AuthServerUrl = go_config.Config.MustGetString(`/authServer/scheme`) + `://` + go_config.Config.MustGetString(`/authServer/host`)

	external_service.ExternalServiceDriver.Register(`deposit_address`, &external_service2.DepositAddressService)

	//global_api_strategy.GlobalApiStrategyDriver.Register(global_api_strategy.GlobalStrategyData{
	//	Strategy: &global_api_strategy2.OpenCensusStrategy,
	//	Disable: go_application.Application.Env == `local`,
	//	Param: global_api_strategy2.OpenCensusStrategyParam{
	//		StackDriverOption: nil,
	//		EnableTrace:       true,
	//		EnableStats:       false,
	//	},
	//})

	service.Service.SetRoutes(
		route.AddressRoute,
		route.TransactionRoute,
		route.WithdrawRoute,
		route.UserRoute,
		manage.MemberRoute,
	)
	service.Service.SetHost(go_config.Config.MustGetString(`host`))
	service.Service.SetPort(go_config.Config.MustGetUint64Default(`port`, 8000))
	service.Service.Run()
}
