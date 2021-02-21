
这里就详细介绍下如何使用go module导入本地包。


我们分两种情况讨论：

- 在同一个项目下
- 不在同一个项目下



在同一个项目下
=================
注意：在一个项目（project）下我们是可以定义多个包（package）的。


假设我们现在有moduledemo和mypackage两个包，其中moduledemo包中会导入mypackage包并使用它的New方法。

目录结构
-----------------
现在的情况是，我们在moduledemo/main.go中调用了mypackage这个包。
```sh
moduledemo
├── go.mod
├── main.go
└── mypackage
    └── mypackage.go
```

导入包
-----------------
我们需要在moduledemo/go.mod中按如下定义：
```sh
module moduledemo

go 1.14
```


mypackage/mypackage.go内容如下：
```golang
package mypackage

import "fmt"

func New(){
	fmt.Println("mypackage.New")
}
```


在moduledemo/main.go中按如下方式导入mypackage
```golang
package main

import (
	"fmt"
	"moduledemo/mypackage"  // 导入同一项目下的mypackage包
)
func main() {
	mypackage.New()
	fmt.Println("main")
}
```

举个例子
----------------
举一反三，假设我们现在有文件目录结构如下：
```sh
└── bubble
    ├── dao
    │   └── mysql.go
    ├── go.mod
    └── main.go
```

其中bubble/go.mod内容如下：
```sh
module github.com/q1mi/bubble

go 1.14
```

bubble/dao/mysql.go内容如下：
```golang
package dao

import "fmt"

func New(){
	fmt.Println("mypackage.New")
}
```

bubble/main.go内容如下：
```golang
package main

import (
	"fmt"
	"github.com/q1mi/bubble/dao"
)
func main() {
	dao.New()
	fmt.Println("main")
}
```




不在同一个项目下
===================

目录结构
```sh
├── moduledemo
│   ├── go.mod
│   └── main.go
└── mypackage
    ├── go.mod
    └── mypackage.go
```


导入包
--------------------

mypackage也需要进行module初始化，即拥有一个属于自己的go.mod文件，内容如下：
```sh
module mypackage

go 1.14
```

然后我们在moduledemo/main.go中按如下方式导入：
```golang
import (
	"fmt"
	"mypackage"
)
func main() {
	mypackage.New()
	fmt.Println("main")
}
```

因为这两个包不在同一个项目路径下，你想要导入本地包，并且这些包也没有发布到远程的github或其他代码仓库地址。这个时候我们就需要在go.mod文件中使用replace指令。

在调用方也就是moduledemo/go.mod中按如下方式指定使用相对路径来寻找mypackage这个包。
```sh
module moduledemo

go 1.14


require "mypackage" v0.0.0
replace "mypackage" => "../mypackage"
```


举个例子
------------------
最后我们再举个例子巩固下上面的内容。

我们现在有文件目录结构如下：
```sh
├── p1
│   ├── go.mod
│   └── main.go
└── p2
    ├── go.mod
    └── p2.go
```
p1/main.go中想要导入p2.go中定义的函数。


p2/go.mod内容如下：
```sh
module liwenzhou.com/q1mi/p2

go 1.14
```


p1/main.go中按如下方式导入
```golang
import (
	"fmt"
	"liwenzhou.com/q1mi/p2"
)

func main() {
	p2.New()
	fmt.Println("main")
}
```
因为我并没有把liwenzhou.com/q1mi/p2这个包上传到liwenzhou.com这个网站，我们只是想导入本地的包，这个时候就需要用到replace这个指令了。


p1/go.mod内容如下：
```sh
module github.com/q1mi/p1

go 1.14


require "liwenzhou.com/q1mi/p2" v0.0.0
replace "liwenzhou.com/q1mi/p2" => "../p2"
```

此时，我们就可以正常编译p1这个项目了。