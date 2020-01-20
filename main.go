package main

import (
	"fmt"
	"github.com/ory/hydra-client-go/client"
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
	api_strategy "github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/service"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	go_redis "github.com/pefish/go-redis"
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
	go_mysql.MysqlHelper.SetLogger(go_logger.Logger)
	go_mysql.MysqlHelper.MustConnectWithMap(mysqlConfig)
	defer go_mysql.MysqlHelper.Close()

	// 初始化redis连接
	redisConfig := go_config.Config.MustGetMap(`redis`)
	go_redis.RedisHelper.SetLogger(go_logger.Logger)
	go_redis.RedisHelper.MustConnectWithMap(redisConfig)
	defer go_redis.RedisHelper.Close()

	service.Service.SetName(`storm钱包服务api`)
	service.Service.SetPath(`/api/storm`)
	api_strategy.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	global.HydraClientInstance = client.New(
		httptransport.New(
			go_config.Config.GetString(`/authServer/host`),
			go_config.Config.GetString(`/authServer/basePath`),
			[]string{go_config.Config.GetString(`/authServer/scheme`)}),
		nil)
	global.AuthServerUrl = go_config.Config.GetString(`/authServer/scheme`) + `://` + go_config.Config.GetString(`/authServer/host`)

	service.Service.SetRoutes(
		route.AddressRoute,
		route.TransactionRoute,
		route.WithdrawRoute,
		route.UserRoute,
		manage.MemberRoute,
	)
	service.Service.SetHost(go_config.Config.GetString(`host`))
	service.Service.SetPort(go_config.Config.GetUint64(`port`))
	service.Service.Run()
}
