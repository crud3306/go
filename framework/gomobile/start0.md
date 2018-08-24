
React Native 调用gomobile 编译的SDK  
地址：  
https://www.studygolang.com/articles/14074  
https://www.oschina.net/translate/ios-and-android-programming-with-go?cmp  
  
  
思路  
-----------  
react native写应用界面,go写底层   
大体思路:react native 可以通过NativeModules 调用原生代码,go代码可以通过gomobile框架编译成Android的.arr文件和iOS的.framework文件,通过添加依赖文件到本地即可实现js->java/oc->go的调用.  
  

整个过程:  
-----------
1,搭建go环境;  
2,clone gomobile 框架,书写go代码,通过gomobile编译成.arr&&.framework;  
3,搭建react native环境,初始化一个react native项目;  
4,在react native app中通过Android Studio中添加.arr依赖,在Xcode中添加.framework依赖;  
5,写Android&&iOS 通过NativeModules调用本地原生代码;  
   
    
第1步），搭建go语言环境  
见go-start/install  


第2步），搭建gomobile环境

具体步骤请参考:
https://github.com/golang/go/wiki/Mobile
需要注意的是:
1,go get golang.org/x/mobile/cmd/gomobile命令被墙,可以自己把gomobile项目clone到$GOPATH/src/golang.org/x下面，然后执行gomobile init

这里编译 gomobile 下example目录下的hello.go文件
go代码很简单:
```golang
package hello

import "fmt"

func Greetings(name string) string {
    return fmt.Sprintf("Welcome To Tradove.\nHello, %s!", name)
}
```
  
然后编译Android .arr文件：  
注意：该命令是 gomobile bind，而不是gomobile build。  
> gomobile bind -target=android golang.org/x/mobile/example/bind/hello  
生成iOS用的.framework文件  
> gomobile bind -target=ios golang.org/x/mobile/example/bind/hello  
备注:  
Mac 下.arr和.framework文件会被放在用户根目录下面  
  


第3步），搭建react native 环境:  

这里不多说,请参考官网:  
https://facebook.github.io/react-native/docs/0.55/getting-started.html  
初始化一个react native 项目  
页面很简单,就一个按钮  
```react
type Props = {};
export default class App extends Component<Props> {
    render() {
        return (
            <View style={styles.container}>
                <TouchableOpacity onPress={() => {
                    this._getFromNative();
                }} style={styles.touchView}><Text
                    style={styles.touchText}>{' get values from native '}</Text></TouchableOpacity>
            </View>
        );
    }
    _getFromNative() {
        const rnParam = Platform.select({
            ios: 'str from ios',
            android: 'str from android'
        });
        GoMobileModules.getNativeGo(rnParam, (str) => {
            alert(str);
        });
    }
}
```

  

第4步）,将.arr和.framework文件依赖进react native项目中

Android部分:
使用android studio 打开react native 工程目录/android(注意弹出提示更新时不要更新)
选择
File->new module->import .jar/.arr package

之后在app/build.gradle文件dependencies里添加以下一行
compile project(':hello')
详细可参考:
Android studio2.3导入aar包


iOS部分:
使用Xcode 打开react native 工程目录/ios

下载好所需要的第三方提供的.framework
2）将第三方.framework文件拷贝到工程所处的文件夹中

选中项目名称
4）选中TARGETS

5）选中Build Phases

6）在Link Binary With Libraries中添加



第5步）,书写react native 调用Android&&iOS调用原生的代码

详细内容请参考:

Android 部分:
GoMobileBridge.class
```java
public class GoMobileBridge implements ReactPackage {
    @Override
    public List<NativeModule> createNativeModules(ReactApplicationContext reactContext) {
        List<NativeModule> modules = new ArrayList<>();
        modules.add(new GoMobileModule(reactContext));
        return modules;
    }

    @Override
    public List<ViewManager> createViewManagers(ReactApplicationContext reactContext) {
        return Collections.emptyList();
    }
}
```
GoMobileModule.class
```java
public class GoMobileModule extends ReactContextBaseJavaModule {
    public GoMobileModule(ReactApplicationContext reactContext) {
        super(reactContext);
    }

    @Override
    public String getName() {
        return this.getClass().getSimpleName();
    }

    @Nullable
    @Override
    public Map<String, Object> getConstants() {
        return super.getConstants();
    }

    @ReactMethod
    public void getNativeGo(String rnStr, Callback callback) {
        String greetings = hello.Hello.greetings(rnStr);
        callback.invoke(greetings);
    }
}
```

iOS原生代码:

GoMobileModule.h
```objectc
#ifndef GoMobileModule_h
#define GoMobileModule_h
#import <Foundation/Foundation.h>
#import <React/RCTBridgeModule.h>

@interface GoMobileModule : NSObject<RCTBridgeModule>
@end
#endif /* GoMobileModule_h */
```
GoMobileModule.m
```objectc
#import "GoMobileModule.h"
#import <React/RCTLog.h>
#import "Hello/Hello.h"
@implementation GoMobileModule
RCT_EXPORT_MODULE();
//从go层获取数据
RCT_EXPORT_METHOD(getNativeGo:(NSString *) rnStr :(RCTResponseSenderBlock)callback{
  NSString *goStr=HelloGreetings(rnStr);
  callback(@[goStr]);
});
@end
```

运行:






