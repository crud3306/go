
  
安装gomobile  
----------------  
go get golang.org/x/mobile/cmd/gomobile
gomobile init

如果以上命令被墙,可以自己把gomobile项目clone到$GOPATH/src/golang.org/x下面，然后执行gomobile init  
注：(gomobile项目地址：https://github.com/golang/mobile/)   
  
  
go 与 android 、ios  
----------------
1） 可以直接用go开发原生的android、ios应用  
2） 可以在一个native 应用里使用 go 的包  
  
详细请参考：  
https://www.oschina.net/translate/ios-and-android-programming-with-go?cmp  
https://github.com/golang/go/wiki/Mobile  
  
  
1） 可以直接用go开发原生的android、ios应用  
Android  
构建一个 Android 的 APK 包  
> gomobile build -target=android golang.org/x/mobile/example/basic  
部署到设备上  
> gomobile install golang.org/x/mobile/example/basic  

iOS  
构建一个 iOS 的 IPA 包  
> gomobile build -target=ios golang.org/x/mobile/example/basic  
部署到设备  
跟 Android 不一样，对于 iOS   来说没有一个统一的部署命令，你需要用你熟知的方式把包拷贝到设备或者模拟器上，例如使用 ios-deploy 工具。  
  
  
  
2） 可以在一个native 应用里使用 go 的包  
native中所需的go包生成方式

然后编译Android .arr文件：  
注意：该命令是 gomobile bind，而不是gomobile build。  
> gomobile bind -target=android golang.org/x/mobile/example/bind/hello  

生成iOS用的.framework文件  
> gomobile bind -target=ios golang.org/x/mobile/example/bind/hello  
备注:  
Mac 下.arr和.framework文件会被放在用户根目录下面  
  
native中使用go包的 参考地址：  
https://www.studygolang.com/articles/14074    
  




