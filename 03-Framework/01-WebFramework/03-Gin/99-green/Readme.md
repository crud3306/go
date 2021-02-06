

框架目录
```
config //配置文件
	dev/config.toml
	qa/config.toml
	online/config.toml
constants //常量
	system.go
	kafka.go
	redis.go
	error_code.go //错误
doc //文档
	sql
library //公用库
	config.go
	request.go
	response.go
	logger.go
	funcs
		common.go
		file.go
		time.go
	q_errors
	q_exception
	db
		mysql.go
		mongodb.go
	cache
		redis.go
	search
		es.go
	mq
		kafka.go
app
	console
		run.go
	controller //入口
	logic //业务逻辑
	dao //数据
	middleware //中间件
	router //路由
	bootstrap.go //初始
	router.go
main.go
```


1130
-----
路由
log


1201
-----
