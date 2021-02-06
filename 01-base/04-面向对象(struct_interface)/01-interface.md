

接口(interfaces)
=============

一般而言，接口定义了一组方法的集合，接口不能被实例化，一个类型可以实现多个接口。

举一个简单的例子，定义一个接口 Person和对应的方法 getName()
```golang
type Person interface {
    getName() string
}

// Student
type Student struct {
    name string
    age  int
}
func (stu *Student) getName() string {
    return stu.name
}

// Worker
type Worker struct {
    name   string
    gender string
}
func (w *Worker) getName() string {
    return w.name
}


func main() {
    // 实例化 Student后，强制类型转换为接口类型 Person。
    var p Person = &Student{
        name: "Tom",
        age:  18,
    }

    fmt.Println(p.getName()) // Tom
}
```

在上面的例子中，我们在 main 函数中尝试将 Student 实例类型转换为 Person，如果 Student 没有完全实现 Person 的方法，比如我们将 (星号Student).getName() 删掉，编译时会出现如下报错信息。
```sh
*Student does not implement Person (missing getName method)
```


但是删除 (星号Worker).getName() 程序并不会报错，因为我们并没有在 main 函数中使用。这种情况下我们如何确保某个类型实现了某个接口的所有方法呢？一般可以使用下面的方法进行检测，如果实现不完整，编译期将会报错。
```sh
var _ Person = (*Student)(nil)
var _ Person = (*Worker)(nil)
```
将空值 nil 转换为 星号Student 类型，再转换为 Person 接口，如果转换失败，说明 Student 并没有实现 Person 接口的所有方法。
Worker 同上。


实例可以强制类型转换为接口，接口也可以强制类型转换为实例。
```golang
func main() {
    var p Person = &Student{
        name: "Tom",
        age:  18,
    }

    stu := p.(*Student) // 接口转为实例
    fmt.Println(stu.getAge())
}
```




空接口
-------------
如果定义了一个没有任何方法的空接口，那么这个接口可以表示任意类型。例如
```golang
func main() {
    m := make(map[string]interface{})
    m["name"] = "Tom"
    m["age"] = 18
    m["scores"] = [3]int{98, 99, 85}
    fmt.Println(m) // map[age:18 name:Tom scores:[98 99 85]]
}
```




继承、多态
-------------
go没有 implements, extends 关键字，所以习惯于 OOP 编程，或许一开始会有点无所适从的感觉。 但go作为一种优雅的语言， 给我们提供了另一种解决方案， 那就是鸭子类型：看起来像鸭子， 那么它就是鸭子.

那么什么是鸭子类型， 如何去实现呢 ？

接下来我会以一个简单的例子来讲述这种实现方案。

首先我们需要一个超类：
```golang
type Animal interface {
    Sleep()
    Age() int
    Type() string
}
```

必然我们需要真正去实现这些的子类:
```golang
type Cat struct {
    MaxAge int
}

func (this *Cat) Sleep() {
    fmt.Println("Cat need sleep")
}
func (this *Cat) Age() int {
    return this.MaxAge
}
func (this *Cat) Type() string {
    return "Cat"
}
type Dog struct {
    MaxAge int
}

func (this *Dog) Sleep() {
    fmt.Println("Dog need sleep")
}
func (this *Dog) Age() int {
    return this.MaxAge
}
func (this *Dog) Type() string {
    return "Dog"
}
```
我们有两个具体实现类 Cat, Dog, 但是Animal如何知道Cat, Dog已经实现了它呢？ 原因在于： Cat, Dog实现了Animal中的全部方法， 那么它就认为这就是我的子类。

那么如何去使用这种关系呢？  

我们使用具体工厂类来构造具体的实现类， 在调用时你知道有这些方法， 但是并不清楚具体的实现， 每一种类型的改变都不会影响到其它的类型。
```golang
package main

import (
    "animals"
    "fmt"
)

func Factory(name string) Animal {
    switch name {
    case "dog":
        return &Dog{MaxAge: 20}
    case "cat":
        return &Cat{MaxAge: 10}
    default:
        panic("No such animal")
    }
}

func main() {
    animal := animals.Factory("dog")
    animal.Sleep()
    fmt.Printf("%s max age is: %d", animal.Type(), animal.Age())
}
```

来看看我们的输出会是什么吧
```sh
> Output:
animals
command-line-arguments
Dog need sleep
Dog max age is: 20
> Elapsed: 0.366s
> Result: Success
```
这就是go中的多态， 是不是比 implements/extends 显示的表明关系更优雅呢。