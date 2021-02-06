package library

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	LogLevel string

	MySQL    MySQLConfig
	Redis    RedisConfig
	ConsoleLog    LogConfig
	HttpLog    LogConfig
}

type MySQLConfig struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	IP   string
	Port int
}

type LogConfig struct {
	FilePath          string
	Level             string
	MaxSize           int
	MaxAge            int
	MaxBackups        int
	OutputProbability float32 //日志输出概率(小数)
	ToStdout          bool    //是否将日志输出到标准输出中
	Compress		  bool    //是否压缩
}

var conf *Config
var configLock sync.Mutex

func GetMySQLConfig() *MySQLConfig {
	return &conf.MySQL
}

func GetRedisConfig() *RedisConfig {
	return &conf.Redis
}

func GetHttpLogConfig() *LogConfig {
	return &conf.HttpLog
}

func GetConsoleLogConfig() *LogConfig {
	return &conf.ConsoleLog
}

func InitConfig(configFile string) *Config {
	if conf != nil {
		return conf
	}

	configLock.Lock()
	defer configLock.Unlock()

	conf = new(Config)
	viper.SetConfigFile(configFile) // 指定配置文件路径
	err := viper.ReadInConfig()               // 读取配置信息
	if err != nil {                           // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("夭寿啦~配置文件被人修改啦...", conf)
		if err := viper.Unmarshal(conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}

		fmt.Println("new config", conf)
	})

	return conf
}
