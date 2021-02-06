

```golang
package xxx

import (
	"sync"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	mysqlEngine *xorm.Engine
	mysqlLock sync.Mutex
)

func NewMySQLEngine() *xorm.Engine {
	if mysqlEngine == nil {
		mysqlLock.Lock()
		defer mysqlLock.Unlock()

		InitMySQLEngine()
	}
	return mysqlEngine
}

func InitMySQLEngine() {
	conf := library.GetDBConfig()

	//dsn := "user:password@(host:port)/dbname?charset=utf8"
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname)

	mysqlEngine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if conf.MaxOpenCons == 0 {
		mysqlEngine.SetMaxOpenConns(DefaultMaxOpenConns)
	} else {
		mysqlEngine.SetMaxOpenConns(conf.MaxOpenCons)
	}

	if conf.MaxIdleCons == 0 {
		mysqlEngine.SetMaxIdleConns(DefaultMaxIdleConns)
	} else {
		mysqlEngine.SetMaxIdleConns(conf.MaxIdleCons)
	}

	if conf.MaxLifeTime == 0 {
		mysqlEngine.SetConnMaxLifetime(DefaultMaxLifeTime)
	} else {
		mysqlEngine.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifeTime))
	}

	// 是否打印sql.default false
	mysqlEngine.ShowSQL(conf.ShowSql)

	// 日志
	if conf.LogPath != "" {
		f, err := os.Create(conf.LogPath)
		if err != nil {
			panic(err)
		} 

		logger := log.NewSimpleLogger(f)
		logger.ShowSQL(conf.ShowSql)              // 是否打印sql.default false
		logger.SetLevel(log.LogLevel(conf.Level)) // 日志等级
		mysqlEngine.SetLogger(logger)
	}
}
```