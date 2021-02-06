

VS Code Remote SSH配置
==============
点击左侧插件图标，搜索ssh插件(Remote - SSH)，一般第一个就是，点击install按扭，安装即可

安装后，IDE左侧出多出一个远程图标。

点击左侧远程图标，点击配置按扭，在弹出中选择配置文件(/xxx/.ssh/config)，设置
```sh
# Read more about SSH config files: https://linux.die.net/man/5/ssh_config
Host web_dev01
    HostName xxx.xxx.xx.xx
    User 你的username
    port 22

Host web_dev02
    HostName xxx.xxx.xx.xx
    User 你的username
    port 22
```

各参数的含义：
```sh
Host 连接的主机的名称，可自定
Hostname 远程主机的IP地址
User 用于登录远程主机的用户名
Port 用于登录远程主机的端口
IdentityFile 本地的id_rsa的路径

如果需要多个连接，可按照如上配置多个。
```

配好后，在左侧会出现在配置文件里配的Host列表。选项一个其中一个Host，点击右侧的"添加目录"按扭，输入密码，打开对应的目录即可。




安装Go插件
==============
这个必须装，Go语言的支持。
点击左侧插件图标，搜索Go插件，一般第一个就是，点击install按扭，安装即可

安装后，打开了编译器，刚写了个Hello World准备震惊世界时，结果第一行package都还没有写完，VSCode就给我提示各种要安装的包。按需要点击安装。

天真的你一而再再而三的点Install All，放心，不管用的，咱家的墙比天还高，搭梯子也翻不过去。
```sh
Installing golang.org/x/lint/golint FAILED
Installing xxx FAILED
...
```
安装不成功能，请手动go get xxx，然手go install xxx。



如果需要用到代码跳转
-----------
代码最好以go mod方式来管理，这样代码不用放在GOPATH目录下
```sh
go mod init xxx
go mod tidy
go mod vendor
```

快捷键：ctrl+英文逗号 打开settings面板
	搜索 use language server
	Use Language Server 改为选中状态,就可以跳转了

注意：
	如要代码是远程，则setting配置有几个tab项目，分别是 user / ssh / workspace, Use Language Server只勾选后面两个的，user的别勾选。不然跳转无效

或者
设置搜索 Docs Tool，把 Docs Tool改成gogetdoc或者guru试试，我的用guru就可以了,其他的可以尝试一下


安装go依赖包
------------
- 使用快捷键：command+shift+P 或者 F1键，打开vs顶部安装输入框;
- 点击第一项 Go:Install/Update tools，进入包选择界面，钩选你想安装的包，点击Ok按扭。

如果安装失败(因大防火墙的原因)，会提示出哪个包失败，在提示浮层中拷贝出包地址，然后手动执行go get xxx包，然后在vscode中手动执行go install xxx包。


需要安装的包：  

gocode：自动完成上下文
https://github.com/mdempsky/gocode

goLint 代码质量检测
https://github.com/golang/lint

go-outline：文件大纲
https://github.com/lukehoban/go-outline

eg：
```sh
go get -u -v github.com/nsf/gocode
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/golang/lint/golint
go get -u -v github.com/lukehoban/go-find-references
go get -u -v github.com/lukehoban/go-outline
go get -u -v sourcegraph.com/sqs/goreturns
go get -u -v golang.org/x/tools/cmd/gorename
go get -u -v github.com/tpng/gopkgs
go get -u -v github.com/newhook/go-symbols
```







