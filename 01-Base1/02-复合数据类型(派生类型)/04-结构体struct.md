

结构体 struct
=============

结构体（Struct）是复合类型，可以封装属性和操作（即字段和方法）。是由零个或多个任意类型的值聚合成的实体，每个值称为结构体的成员。

Go 中的结构体类似面向对象编程中的轻量级类，但 Go 中没有类的概念，所以结构体尤为重要。



示例
-------------
```golang
package main

import "fmt"

// 定义一个结构体 person
type person struct {
    name string
    age  int
}

func main() {
    var p person   // 声明一个 person 类型变量 p
    p.name = "max" // 赋值
    p.age = 12
    fmt.Println(p) // 输出：{max 12}

    p1 := person{name: "mike", age: 10} // 直接初始化一个 person
    fmt.Println(p1.name)                // 输出：mike

    p2 := new(person) // new函数分配一个指针，指向 person 类型数据
    p2.name = `张三`
    p2.age = 15
    fmt.Println(*p2) // 输出：{张三 15}
}
```



创建
-------------
定义结构体的一般语法如下：
```golang
type identifier struct {
    field1 type1
    field2 type2
    ...
}
```
如果不需要 field，可以将其命名为 _
结构体成员遵循 Go 的导出原则，一个结构体可以同时包含导出和未导出的成员


初始化
--------------
定义结构体类型：Person
```golang
type Person struct {
    Name   string
    Gender string
    Age    uint8
}
```
结构体的零值是每个成员都是零值：
```golang
// 零值化
var p1 Person
fmt.Printf("%+v\n", p1)  // {Name: Gender: Age:0}

// 选择器
p1.Name = "Ming"
p1.Age = 32
fmt.Printf("%+v\n", p1)  // {Name:Ming Gender: Age:32}
```
使用字段名：

- 使用字段名时，顺序可以打乱；
- 忽略字段名时，需要严格按照定义的顺序，且数量得对应；


未提供初始值的字段将使用该字段的类型的零值
```golang
p2 := Person{Name: "Qiang", Gender: "Male", Age: 30}
fmt.Printf("%+v\n", p2)    // {Name:Qiang Gender:Male Age:30}

p3 := Person{Gender: "Male", Name: "Hua"}
fmt.Printf("%+v\n", p3)    // {Name:Hua Gender:Male Age:0}

p4 := Person{"Jack", "Male", 30}
fmt.Printf("%+v\n", p4)    // {Name:Jack Gender:Male Age:30}
```

使用字段赋值时和 map 的键值对的写法有点类似，但请注意区别。

使用 new()：使用 new 函数给一个新的结构体变量分配内存，它返回指向 星号 已分配内存的指针
```golang
var p5 *Person
fmt.Println(p5 == nil)  // true

p5 = new(Person)
fmt.Printf("%+v\n", p5)   // &{Name: Gender: Age:0}

var p6 *Person = new(Person)
fmt.Printf("%+v\n", p6)  // &{Name: Gender: Age:0}

p7 := new(Person)
fmt.Printf("%+v\n", p7)  // &{Name: Gender: Age:0}

// 指针也有选择器，也可以使用类似属性的操作
p7.Name = "Ming"
fmt.Printf("%+v\n", p7)  // &{Name:Ming Gender: Age:0}
```

混合字面量语法（composite literal syntax）是一种简写，底层仍然会调用 new ()：
```golang
p8 := &Person{"Dai", "Male", 42}
fmt.Printf("%+v\n", p8)  // &{Name:Dai Gender:Male Age:42}
````
注意，new(Type) 和 &Type{} 是等价的，见内存布局。

上述中的 p1、p2 等通常被称做类型 Person 的一个实例（instance）或对象（object）。

简写
当字段的类型相同时，可以写在同一行：
```golang
type T struct {
    a, b int
}

t := T{1, 2}
fmt.Println(t.a)  // 1
fmt.Println(t.b)  // 2
```


工厂方法
--------------
Go 中没有类的概念，不存在 OOP 中的构造方法，但是我们可以使用工厂方法实现类似的行为。

按照惯例，工厂函数的名字以 new 或 New 开头：
```golang
package main

import (
    "fmt"
    "unsafe"
)

type Person struct {
    name string
    age uint8
}

func NewPerson(name string, age uint8) *Person {
    // 注意，这个返回的是局部变量的地址.
    // 如果在 C++ 中是个典型的错误，Go 中没问题
    return &Person{name, age}
}

func main() {
    p := NewPerson("Jack", 20)
    fmt.Printf("%+v\n", p)  // &{name:Jack age:20}
    fmt.Println(unsafe.Sizeof(Person{}))  // 查看一个实例占用了多少内存
}
```
这类似于面向对象中的实例化一个类。



选择器
--------------
无论一个结构体类型还是一个结构体类型指针，都使用同样的选择器符来引用结构体的字段：
```golang
type Person struct {
    name string
    age  uint8
}
p1 := Person{"Jack", 20}      // 结构体类型变量
p2 := &Person{"Duncan", 42}   // 指向一个结构体类型变量的指针

fmt.Println(p1.name, p1.age)  // Jack 20
fmt.Println(p2.name, p2.age)  // Duncan 42
```


内存布局
-------------
用图说明结构体类型 Point 的实例和指向它的指针的内存布局：
```golang
type Point struct { 
    x, y int 
}

//new(Point) 和 &Point{} 返回的是指针。
````
Go 语言中，结构体和它所包含的数据在内存中是以连续块的形式存在的，即使结构体中嵌套有其他的结构体，这在性能上带来了很大的优势。不像 Java 中的引用类型，一个对象和它里面包含的对象可能会在不同的内存空间中。

举例：
```golang
package main

import (
    "fmt"
    "strings"
)

type Person struct {
    name string
    age uint8
}

// 这是方法，详细请见后续笔记
func update(p *Person) {
    p.name = strings.ToUpper(p.name)
}

func main() {
    p1 := Person{"Duncan", 43}
    update(&p1)
    fmt.Println(p1.name)  // DUNCAN

    p2 := new(Person)
    p2.name = "Curry"     // 指针也能使用选择器，Go自动做了转换，不像C++中那样需要使用 -> 操作符
    update(p2)
    fmt.Println(p2.name)  // CURRY

    (*p2).name = "Tony"   // 也可以取指针的值
    update(p2)
    fmt.Println(p2.name)  // TONY

    p3 := &Person{"Durant", 30}
    update(p3)
    fmt.Println(p3.name)  // DURANT
}
```



匿名字段
------------
匿名字段，即没有名字的字段，声明时只指定类型。

类型必须是命名的类型或指向一个命名的类型的指针
匿名字段的名字默认是其类型的名字，因此不能同时包含两个相同类型的匿名字段，这会导致名字冲突
匿名字段也有可见性的规则约束
```golang
type Person struct {
    name string
    age int
    int  // 匿名字段，不能有两个匿名字段 int
    string
}

p1 := Person{}
fmt.Printf("%+v\n", p1)  // {name: age:0 int:0 string:}
fmt.Println(p1.int)  // 匿名字段的名字是其类型名

p2 := new(Person)
p2.int = 10
p2.string = "abc"
fmt.Printf("%+v\n", p2)  // &{name: age:0 int:10 string:abc}

p := Person{"a", 10, 10}
fmt.Printf("%+v\n", p)  // {Name:a Age:10 int:10}
```


内嵌结构体
--------------
当一个结构体嵌套进另一个结构体在 Go 语言中也是很常见的，举例：
```golang
type Address struct {
    province string
    city string
}

type User struct {
    name string
    age int
    address Address
}

u := &User{
    name: "Ming",
    age: 30,
    address: Address{
        province: "Jiangsu",
        city: "Nanjing",
    },
}
fmt.Printf("%+v\n", u)  // &{name:Ming age:30 address:{province:Jiangsu city:Nanjing}}
````

当结构体作为匿名成员的时候，会有一些特殊的用法：
```golang
type Address struct {
    province string
    city string
}

type User struct {
    name string
    age int
    Address
}

// 方法一：正常直观方式定义
u1 := &User{
    name: "Ming",
    age: 30,
    Address: Address{
        province: "Jiangsu",
        city: "Nanjing",
    },
}
fmt.Printf("%+v\n", u1)  // &{name:Ming age:30 Address:{province:Jiangsu city:Nanjing}}

// 同上
var u2 User
u2.name = "Qiang"
u2.age = 35
u2.Address = Address{province: "Jiangsu", city: "Suzhou"}
fmt.Printf("%+v\n", u2)  // {name:Qiang age:35 Address:{province:Jiangsu city:Suzhou}}

// 方法二：匿名嵌入时可以直接访问叶子属性而不需要给出完整的路径
var u3 User
u3.name = "A"
u3.age = 40
u3.province = "Jiangsu"
u3.city = "Wuxi"
fmt.Printf("%+v\n", u3)  // {name:A age:40 Address:{province:Jiangsu city:Wuxi}}

// 但下面的方式是错误的，编译不能通过
// cannot use promoted field Address.province in struct literal of type User
// cannot use promoted field Address.city in struct literal of type User
u4 := User{
    name: "A",
    age: 29,
    province: "Jiangsu",
    city: "Wuxi",
}
fmt.Printf("%+v\n", u4)
```
内嵌结构体可以用来实现类似 OOP 中的继承



命名冲突
---------------
当两个字段拥有相同的名字（可能是继承来的名字）时：

外层名字会覆盖内层名字（但是两者的内存空间都保留），利用此特性可以实现字段或方法重载
如果相同的名字在同一级别出现了两次，当使用这个名字时将会引发一个错误（不使用没关系）。
```golang
type Group struct {
    name string
    id int
}

type Team struct {
    id int
}

type User struct {
    name string
    Group
    Team
}

var u User

// 当结构体中的字段与匿名成员内的字段冲突时，结构体的字段覆盖匿名成员的字段，但匿名成员的字段仍然存在
u.name = "A"
fmt.Printf("%+v\n", u)  // {name:A Group:{name: id:0} Team:{id:0}}

u.Group.name = "B"
fmt.Printf("%+v\n", u)  // {name:A Group:{name:B id:0} Team:{id:0}}

// 当匿名成员中的字段出现冲突，使用时会引起错误
fmt.Println(u.Group.id)  // 正常
fmt.Println(u.Team.id)   // 正常
fmt.Println(u.id)        // ambiguous selector u.id
```



其它
===========
标签
-----------
tag 是结构体的元信息，是一个附属于字段的字符串，可以是文档或其他的重要标记。

标签的内容不能在一般的编程中使用，只有包 reflect 能获取它。

根据 Go conventions，tag 有一些约定：

tag 字符串是可选的空格分隔的 key-value
key 是非空字符串，由空格(' ')、引号(' ')和冒号(':')以外的非控制字符组成
value 使用双引号 " 和 Go 字符串
检查字段是否为可导出
```golang
import (
    "encoding/json"
    "fmt"
    "reflect"
)

type User struct {
    Name string  `json:"xxx"`
    Sex string   `json:"sex"`
    Age uint8
}

func main() {
    u := User{"Ming", "male", 30}

    // 通过 reflect 包获取 tag
    fmt.Println(reflect.TypeOf(u).Field(0).Tag)  // json:"xxx"
    fmt.Println(reflect.TypeOf(u).Field(1).Tag)  // json:"sex"
    fmt.Println(reflect.TypeOf(u).Field(2).Tag)

    data, _ := json.Marshal(u)
    fmt.Printf("%s\n", string(data))  
    // {"xxx":"Ming","sex":"male","Age":30}
    // 注意字段名，Age 没添加 tag，默认为原来的
    
    // User 结构体的非导出字段不会被 json 访问到
}
```


递归结构体
------------
结构体类型可以引用自身的指针类型，这就是递归结构体。

定义一个二叉树：
```golang
type tree struct {
    value int
    left *tree
    right *tree
}
```


结构体比较
------------
如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的:
```golang
type Point struct {
    x, y int
}

p1 := Point{1, 2}
p2 := Point{1, 2}
p3 := Point{2, 3}
fmt.Println(p1 == p2)  // true
fmt.Println(p1 == p3)  // false
```