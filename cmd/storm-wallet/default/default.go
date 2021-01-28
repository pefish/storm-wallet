package _default

import (
	"flag"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy2 "github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	go_redis "github.com/pefish/go-redis"
	"log"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"
	"wallet-storm-wallet/route"
	"wallet-storm-wallet/route/manage"
)

type DefaultCommand struct {
}

func (d DefaultCommand) DecorateFlagSet(flagSet *flag.FlagSet) error {
	return nil
}

func (d DefaultCommand) Start(data *commander.StartData) error {
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
	global_api_strategy2.ParamValidateStrategyInstance.SetErrorCode(constant.PARAM_ERROR)

	logger.LoggerDriverInstance.Register(go_logger.Logger)

	service.Service.SetRoutes(
		route.AddressRoute,
		route.TransactionRoute,
		route.WithdrawRoute,
		route.UserRoute,
		manage.MemberRoute,
	)
	service.Service.SetHost(go_config.Config.MustGetString(`host`))
	service.Service.SetPort(go_config.Config.MustGetUint64Default(`port`, 8000))

	global.Global.Init()

	err := service.Service.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (d DefaultCommand) OnExited(data *commander.StartData) error {
	return nil
}
