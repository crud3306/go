

反射 与 断言
=============



示例
-------------
```golang
// 反射
type User struct {
	Name string
	Age int
}

func interface2Sclice(obj interface{}) {
  var list []User
  if reflect.TypeOf(obj).Kind() == reflect.Slice {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(User))
		}
	} 

	fmt.Println(obj, list)
}

user := User{"zhang", 11}
interface2Sclice(user)


// 断言
// Interface2str 提取string
func Interface2str1(v interface{}) string {
    s := ""

    switch v.(type) {
    case string:
        s, _ = v.(string)
    case float64:
        s = strconv.FormatFloat(v.(float64), 'f', 6, 64)
    case int64:
        s = strconv.FormatInt(v.(int64), 10)
    case int32:
        s = strconv.FormatInt(int64(v.(int32)), 10)
    case int:
        s = strconv.Itoa(v.(int))
    default:
        break
    }

    return s
}

// Interface2str 提取string
func Interface2str2(v interface{}) string {
    s := ""

    switch v.(type) {
    case string:
        s = v.(string)
    case int:
        s = fmt.Sprintf("%d", v)
    case float64:
        s = fmt.Sprintf("%f", v)
    case int64:
        s = fmt.Sprintf("%d", v)
    default:
        break
    }

    return s
}

// 注意，如果断言类型不对用panic错误
ok:=arg.(A)

// 如果要防止该错误，接收时用两个参数
_, ok:=arg.(A)                       
if ok {
    fmt.Println("yes")              
} else {
    fmt.Println("no")               
} 
```



reflect的基本功能TypeOf和ValueOf
===========
既然反射就是用来检测存储在接口变量内部(值value；类型concrete type) pair对的一种机制。那么在Golang的reflect反射包中有什么样的方式可以让我们直接获取到变量内部的信息呢？ 它提供了两种类型（或者说两个方法）让我们可以很容易的访问接口变量内容，分别是reflect.ValueOf() 和 reflect.TypeOf()，看看官方的解释
```sh
// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i.  ValueOf(nil) returns the zero 
func ValueOf(i interface{}) Value {...}

翻译一下：ValueOf用来获取输入参数接口中的数据的值，如果接口为空则返回0


// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {...}

翻译一下：TypeOf用来动态获取输入参数接口中的值的类型，如果接口为空则返回nil
```


reflect.TypeOf()是获取pair中的type，reflect.ValueOf()获取pair中的value，示例如下：
```golang
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var num float64 = 1.2345

    fmt.Println("type: ", reflect.TypeOf(num))
    fmt.Println("value: ", reflect.ValueOf(num))


    // a := []string{"foo", "bar", "baz"}
    // fmt.Println(reflect.TypeOf(a), reflect.ValueOf(a), 
    // 	reflect.ValueOf(a).Kind(),
	   //  reflect.ValueOf(a).Index(0).Kind(),
	   //  reflect.ValueOf(a).Kind().String() == "slice",
	   //  reflect.ValueOf(a).Kind() == reflect.Slice,
	   //  reflect.TypeOf(a).NumField())
}
```
运行结果:
```sh
type:  float64
value:  1.2345
```
说明
- reflect.TypeOf： 直接给到了我们想要的type类型，如float64、int、各种pointer、struct 等等真实的类型
- reflect.ValueOf：直接给到了我们想要的具体的值，如1.2345这个具体数值，或者类似&{1 "Allen.Wu" 25} 这样的结构体struct的值

也就是说明反射可以将“接口类型变量”转换为“反射类型对象”，反射类型指的是reflect.Type和reflect.Value这两种从relfect.Value中获取接口interface的信息

当执行reflect.ValueOf(interface)之后，就得到了一个类型为”relfect.Value”变量，可以通过它本身的Interface()方法获得接口变量的真实内容，然后可以通过类型判断进行转换，转换为原有真实类型。不过，我们可能是已知原有类型，也有可能是未知原有类型，因此，下面分两种情况进行说明。


已知原有类型【进行“强制转换”】
-------------
已知类型后转换为其对应的类型的做法如下，直接通过Interface方法然后强制转换，如下：
```golang
realValue := value.Interface().(已知的类型)
```
示例如下：
```golang
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var num float64 = 1.2345

    pointer := reflect.ValueOf(&num)
    value := reflect.ValueOf(num)

    // 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
    // Golang 对类型要求非常严格，类型一定要完全符合
    // 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
    convertPointer := pointer.Interface().(*float64)
    convertValue := value.Interface().(float64)

    fmt.Println(convertPointer)
    fmt.Println(convertValue)
}
```
运行结果：
```sh
0xc42000e238
1.2345
```
说明
转换的时候，如果转换的类型不完全符合，则直接panic，类型要求非常严格！

也就是说反射可以将“反射类型对象”再重新转换为“接口类型变量”


未知原有类型【遍历探测其Filed】
-------------
很多情况下，我们可能并不知道其具体类型，那么这个时候，该如何做呢？需要我们进行遍历探测其Filed来得知，示例如下:
```golang
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Id   int
    Name string
    Age  int
}

func (u User) ReflectCallFunc() {
    fmt.Println("Allen.Wu ReflectCallFunc")
}

func main() {

    user := User{1, "Allen.Wu", 25}

    DoFiledAndMethod(user)

}

// 通过接口来获取任意参数，然后一一揭晓
func DoFiledAndMethod(input interface{}) {

    getType := reflect.TypeOf(input)
    fmt.Println("get Type is :", getType.Name())

    getValue := reflect.ValueOf(input)
    fmt.Println("get all Fields is:", getValue)

    // 获取方法字段
    // 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
    // 2. 再通过reflect.Type的Field获取其Field
    // 3. 最后通过Field的Interface()得到对应的value
    for i := 0; i < getType.NumField(); i++ {
        field := getType.Field(i)
        value := getValue.Field(i).Interface()
        fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
    }

    // 获取方法
    // 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
    for i := 0; i < getType.NumMethod(); i++ {
        m := getType.Method(i)
        fmt.Printf("%s: %v\n", m.Name, m.Type)
    }
}
```
运行结果：
```sh
get Type is : User
get all Fields is: {1 Allen.Wu 25}
Id: int = 1
Name: string = Allen.Wu
Age: int = 25
ReflectCallFunc: func(main.User)
```
说明  
通过运行结果可以得知获取未知类型的interface的具体变量及其类型的步骤为：
- 先获取interface的reflect.Type，然后通过NumField进行遍历
- 再通过reflect.Type的Field获取其Field
- 最后通过Field的Interface()得到对应的value


通过运行结果可以得知获取未知类型的interface的所属方法（函数）的步骤为：
- 先获取interface的reflect.Type，然后通过NumMethod进行遍历
- 再分别通过reflect.Type的Method获取对应的真实的方法（函数）
- 最后对结果取其Name和Type得知具体的方法名

也就是说反射可以将“反射类型对象”再重新转换为“接口类型变量”。   
struct 或者 struct 的嵌套都是一样的判断处理方式



通过reflect.ValueOf来进行方法的调用
----------------
这算是一个高级用法了，前面我们只说到对类型、变量的几种反射的用法，包括如何获取其值、其类型、如果重新设置新值。但是在工程应用中，另外一个常用并且属于高级的用法，就是通过reflect来进行方法【函数】的调用。比如我们要做框架工程的时候，需要可以随意扩展方法，或者说用户可以自定义方法，那么我们通过什么手段来扩展让用户能够自定义呢？关键点在于用户的自定义方法是未可知的，因此我们可以通过reflect来搞定

示例如下：
```golang
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Id   int
    Name string
    Age  int
}

func (u User) ReflectCallFuncHasArgs(name string, age int) {
    fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User) ReflectCallFuncNoArgs() {
    fmt.Println("ReflectCallFuncNoArgs")
}

// 如何通过反射来进行方法的调用？
// 本来可以用u.ReflectCallFuncXXX直接调用的，但是如果要通过反射，那么首先要将方法注册，也就是MethodByName，然后通过反射调动mv.Call

func main() {
    user := User{1, "Allen.Wu", 25}
    
    // 1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理
    getValue := reflect.ValueOf(user)

    // 一定要指定参数为正确的方法名
    // 2. 先看看带有参数的调用方法
    methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
    args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
    methodValue.Call(args)

    // 一定要指定参数为正确的方法名
    // 3. 再看看无参数的调用方法
    methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
    args = make([]reflect.Value, 0)
    methodValue.Call(args)
}
```

运行结果：
```sh
ReflectCallFuncHasArgs name:  wudebao , age: 30 and origal User.Name: Allen.Wu
ReflectCallFuncNoArgs
```
说明
- 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理
- reflect.Value.MethodByName这.MethodByName，需要指定准确真实的方法名字，如果错误将直接panic，MethodByName返回一个函数值对应的reflect.Value方法的名字。
- []reflect.Value，这个是最终需要调用的方法的参数，可以没有或者一个或者多个，根据实际参数来定。
- reflect.Value的 Call 这个方法，这个方法将最终调用真实的方法，参数务必保持一致，如果reflect.Value'Kind不是一个方法，那么将直接panic。
- 本来可以用u.ReflectCallFuncXXX直接调用的，但是如果要通过反射，那么首先要将方法注册，也就是MethodByName，然后通过反射调用methodValue.Call



通过reflect.Value设置实际变量的值
--------------
reflect.Value是通过reflect.ValueOf(X)获得的，只有当X是指针的时候，才可以通过reflec.Value修改实际变量X的值，即：要修改反射类型的对象就一定要保证其值是“addressable”的。

示例如下：
```golang
package main

import (
    "fmt"
    "reflect"
)

func main() {

    var num float64 = 1.2345
    fmt.Println("old value of pointer:", num)

    // 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
    pointer := reflect.ValueOf(&num)
    newValue := pointer.Elem()

    fmt.Println("type of pointer:", newValue.Type())
    fmt.Println("settability of pointer:", newValue.CanSet())

    // 重新赋值
    newValue.SetFloat(77)
    fmt.Println("new value of pointer:", num)

    ////////////////////
    // 如果reflect.ValueOf的参数不是指针，会如何？
    pointer = reflect.ValueOf(num)
    //newValue = pointer.Elem() // 如果非指针，这里直接panic，“panic: reflect: call of reflect.Value.Elem on float64 Value”
}
```
运行结果：
```sh
old value of pointer: 1.2345
type of pointer: float64
settability of pointer: true
new value of pointer: 77
```
说明
- 需要传入的参数是* float64这个指针，然后可以通过pointer.Elem()去获取所指向的Value，注意一定要是指针。
- 如果传入的参数不是指针，而是变量，那么
- 通过Elem获取原始值对应的对象则直接panic
- 通过CanSet方法查询是否可以设置返回false
- newValue.CantSet()表示是否可以重新设置其值，如果输出的是true则可修改，否则不能修改，修改完之后再进行打印发现真的已经修改了。
- reflect.Value.Elem() 表示获取原始值对应的反射对象，只有原始对象才能修改，当前反射对象是不能修改的
- 也就是说如果要修改反射类型对象，其值必须是“addressable”【对应的要传入的是指针，同时要通过Elem方法获取原始值对应的反射对象】
- struct 或者 struct 的嵌套都是一样的判断处理方式



总结：
=============
- 反射可以大大提高程序的灵活性，使得interface{}有更大的发挥余地
	反射必须结合interface才玩得转  
	变量的type要是concrete type的（也就是interface变量）才有反射一说  

- 反射可以将“接口类型变量”转换为“反射类型对象”
	反射使用 TypeOf 和 ValueOf 函数从接口中获取目标对象信息  

- 反射可以将“反射类型对象”转换为“接口类型变量
	reflect.value.Interface().(已知的类型)  
	遍历reflect.Type的Field获取其Field  
 
- 通过反射可以“动态”调用方法
	因为Golang本身不支持模板，因此在以往需要使用模板的场景下往往就需要使用反射(reflect)来实现  

- 反射可以修改反射类型对象，但是其值必须是“addressable”
	想要利用反射修改对象状态，前提是 interface.data 是 settable,即 pointer-interface  







来源：https://www.jianshu.com/p/b46b1ccd2757
