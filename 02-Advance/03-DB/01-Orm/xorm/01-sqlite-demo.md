

```golang
package library

import (
	"sync"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	sqliteEngine *xorm.Engine
	sqliteLock sync.Mutex
)

func NewSqliteEngine() *xorm.Engine {
	if sqliteEngine == nil {
		sqliteLock.Lock()
		defer sqliteLock.Unlock()

		InitSqliteEngine()
	}

	return sqliteEngine
}

func InitSqliteEngine() {
	dbConfig := GetSqliteConfig()

	engine, err := xorm.NewEngine("sqlite3", dbConfig.FilePath)
	if err != nil {
		panic(err)
	}
	engine.SetMaxOpenConns(30)
	engine.SetConnMaxLifetime(10)
	engine.SetMaxIdleConns(5)
	// 是否打印sql.default false
	engine.ShowSQL(true)

	// 日志
	if dbConfig.LogPath != "" {
		f, err := os.Create(dbConfig.LogPath)
		if err != nil {
			panic(err)
		} 

		logger := log.NewSimpleLogger(f)
		logger.ShowSQL(dbConfig.ShowSql)              // 是否打印sql.default false
		logger.SetLevel(log.LogLevel(dbConfig.LogLevel)) // 日志等级
		engine.SetLogger(logger)
	}

	globalSqliteEngine = engine
}
```