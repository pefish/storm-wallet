package main

import (
	"fmt"
	"github.com/ory/hydra-client-go/client"
	go_core "github.com/pefish/go-core"
	external_service "github.com/pefish/go-core/driver/external-service"
	global_api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy2 "github.com/pefish/go-core/global-api-strategy"
	"os"
	"runtime/debug"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"
	"wallet-storm-wallet/route"
	"wallet-storm-wallet/route/manage"

	httptransport "github.com/go-openapi/runtime/client"
	go_application "github.com/pefish/go-application"
	go_config "github.com/pefish/go-config"
	go_http "github.com/pefish/go-http"
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

	go_http.Http = go_http.NewHttpRequester(go_http.WithTimeout(20 * time.Second))

	// 处理日志
	env := go_config.Config.GetString(`env`)
	go_application.Application.SetEnv(env)
	go_logger.Logger.Init(go_core.Service.GetName(), ``)
	logger.LoggerDriver.Register(go_logger.Logger)

	// 初始化数据库连接
	mysqlConfig := go_config.Config.MustGetMap(`mysql`)
	go_mysql.MysqlHelper.SetLogger(go_logger.Logger)
	go_mysql.MysqlHelper.MustConnectWithMap(mysqlConfig)
	defer go_mysql.MysqlHelper.Close()

	// 初始化redis连接
	redisConfig := go_config.Config.MustGetMap(`redis`)
	go_redis.RedisHelper.SetLogger(go_logger.Logger)
	go_redis.RedisHelper.MustConnectWithMap(redisConfig)
	defer go_redis.RedisHelper.Close()

	go_core.Service.SetName(`storm钱包服务api`)
	go_core.Service.SetPath(`/api/storm`)
	global_api_strategy2.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	global.HydraClientInstance = client.New(
		httptransport.New(
			go_config.Config.GetString(`/authServer/host`),
			go_config.Config.GetString(`/authServer/basePath`),
			[]string{go_config.Config.GetString(`/authServer/scheme`)}),
		nil)
	global.AuthServerUrl = go_config.Config.GetString(`/authServer/scheme`) + `://` + go_config.Config.GetString(`/authServer/host`)

	external_service.ExternalServiceDriver.Register(`deposit_address`, &external_service2.DepositAddressService)

	global_api_strategy.GlobalApiStrategyDriver.Register(global_api_strategy.GlobalStrategyData{
		Strategy: &global_api_strategy2.OpenCensusStrategy,
		Disable: go_application.Application.Env == `local`,
		Param: global_api_strategy2.OpenCensusStrategyParam{
			StackDriverOption: nil,
			EnableTrace:       true,
			EnableStats:       false,
		},
	})

	go_core.Service.SetRoutes(
		route.AddressRoute,
		route.TransactionRoute,
		route.WithdrawRoute,
		route.UserRoute,
		manage.MemberRoute,
	)
	go_core.Service.SetHost(go_config.Config.GetString(`host`))
	go_core.Service.SetPort(go_config.Config.GetUint64(`port`))
	go_core.Service.Run()
}
