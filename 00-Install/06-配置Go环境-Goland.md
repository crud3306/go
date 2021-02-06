



Goland中 gomod模式启用情况下，代码无法跳转问题
==============
```sh
去设置中调整  
	file - settings - Go - GoModules - 勾选中(Enable Go Modules integration)
```


Golang中 deployment代码至远程服务器
==============
```sh
打开 Tools - Deployment  
	在弹出框中的左上角，点击"+"，添加新的sftp，设个名称。
	然后为该sftp配置相关信息
		右侧的-Connection标签 
			ssh configuration: ip\port\user\password 
			Root path: 远程服务器的要同步的目录的绝对地址
			Web server URL: 默认http://即可

		右侧的-Mapping标签 
			Local path: 即当前项目在本机的目录
			Deployment path: /，注意这里一定要配，不能为空，不然deployment用不了。
			Web path: /

去设置中同步的配置为，保存文件时自动同步
	file - settings - Build,... - Deployment - Options 
		upload changed files automatically ... 后的选项，选择 on explicit save action(Ctrl+S) 
```



