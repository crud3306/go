

```golang
package library

import (
	"sync"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	pgEngine *xorm.Engine
	pgLock sync.Mutex
)

func NewPgEngine() *xorm.Engine {
	if pgEngine == nil {
		pgLock.Lock()
		defer pgLock.Unlock()

		InitPgEngine()
	}

	return pgEngine
}

func InitPgEngine() {
	// GetDBConfig
	conf := library.GetDBConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	pgEngine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if conf.MaxOpenCons == 0 {
		pgEngine.SetMaxOpenConns(DefaultMaxOpenConns)
	} else {
		pgEngine.SetMaxOpenConns(conf.MaxOpenCons)
	}

	if conf.MaxIdleCons == 0 {
		pgEngine.SetMaxIdleConns(DefaultMaxIdleConns)
	} else {
		pgEngine.SetMaxIdleConns(conf.MaxIdleCons)
	}

	if conf.MaxLifeTime == 0 {
		pgEngine.SetConnMaxLifetime(DefaultMaxLifeTime)
	} else {
		pgEngine.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifeTime))
	}

	// 是否打印sql.default false
	pgEngine.ShowSQL(conf.ShowSql)

	// 日志
	if conf.LogPath != "" {
		f, err := os.Create(conf.LogPath)
		if err != nil {
			println(err)
		} 
		
		logger := log.NewSimpleLogger(f)
		logger.ShowSQL(conf.ShowSql)              // 是否打印sql.default false
		logger.SetLevel(log.LogLevel(conf.Level)) // 日志等级
		pgEngine.SetLogger(logger)
	}
}
```