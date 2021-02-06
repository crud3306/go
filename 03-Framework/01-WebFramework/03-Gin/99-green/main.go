package main

import (
	"fmt"
	"green/app"
	"green/library"
)

func main() {
	configFile := "./config/dev/config.toml"
	library.InitConfig(configFile)

	mysqlConfig := library.GetMySQLConfig()
	fmt.Println("mysqlConfig", mysqlConfig, mysqlConfig.IP, mysqlConfig.Port)

	redisConfig := library.GetRedisConfig()
	fmt.Println("redisConfig", redisConfig, redisConfig.IP, redisConfig.Port)

	httpLogConfig := library.GetHttpLogConfig()
	fmt.Println("httpLogConfig", httpLogConfig, httpLogConfig.FilePath, httpLogConfig.Level)

	library.GetHttpLoggerInstance().Info("aaa==========")
	library.GetHttpLoggerInstance().Error("bbb==========")

	app.Run()
}
