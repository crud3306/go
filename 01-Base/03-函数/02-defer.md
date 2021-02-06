

defer
=========

defer后边会接一个函数，但该函数不会立刻被执行，而是等到包含它的程序返回时(包含它的函数执行了return语句、运行到函数结尾自动返回、对应的goroutine panic）defer函数才会被执行。通常用于资源释放、打印日志、异常捕获等
```golang
func main() {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    /**
     * 这里defer要写在err判断的后边而不是os.Open后边
     * 如果资源没有获取成功，就没有必要对资源执行释放操作
     * 如果err不为nil而执行资源执行释放操作，有可能导致panic
     */
    defer f.Close()
}
```

如果有多个defer函数，调用顺序类似于栈，越后面的defer函数越先被执行(后进先出)
```golang
func main() {
    defer fmt.Println(1)
    defer fmt.Println(2)
    defer fmt.Println(3)
    defer fmt.Println(4)
}
```
结果：
```sh
4
3
2
1
```

如果包含defer函数的外层函数有返回值，而defer函数中可能会修改该返回值，最终导致外层函数实际的返回值可能与你想象的不一致，这里很容易踩坑，来几个
例1
```golang
func f() (result int) {
    defer func() {
        result++
    }()
    return 0
}
```

例2
```golang
func f() (r int) {
    t := 5
    defer func() {
        t = t + 5
    }()
    return t
}
```

例3
```golang
func f() (r int) {
    defer func(r int) {
        r = r + 5
    }(r)
    return 1
}
```

可能你会认为：例1的结果是0，例2的结果是10，例3的结果是6，那么很遗憾的告诉你，这三个结果都错了。  
为什么呢，最重要的一点就是要明白，return xxx这一条语句并不是一条原子指令。

含有defer函数的外层函数，返回的过程是这样的：先给返回值赋值，然后调用defer函数，最后才是返回到更上一级调用函数中，可以用一个简单的转换规则将return xxx改写成
```sh
返回值 = xxx
调用defer函数(这里可能会有修改返回值的操作)
return 返回值
```

例1可以改写成这样
```golang
func f() (result int) {
    result = 0
    //在return之前，执行defer函数
    func() {
        result++
    }()
    return
}
```
所以例1的返回值是1

例2可以改写成这样
```golang
func f() (r int) {
    t := 5
    //赋值
    r = t
    //在return之前，执行defer函数，defer函数没有对返回值r进行修改，只是修改了变量t
    func() {
        t = t + 5
    }
    return
}
```
所以例2的结果是5

例3可以改写成这样
```golang
func f() (r int) {
    //给返回值赋值
    r = 1
    /**
     * 这里修改的r是函数形参的值，是外部传进来的
     * func(r int){}里边r的作用域只该func内，修改该值不会改变func外的r值
     */
    func(r int) {
        r = r + 5
    }(r)
    return
}
```
所以例3的结果是1

defer函数的参数值，是在申明defer时确定下来的  
在defer函数申明时，对外部变量的引用是有两种方式：
- 作为函数参数
作为函数参数，在defer申明时就把值传递给defer，并将值缓存起来，调用defer的时候使用缓存的值进行计算（如上边的例3）

- 作为闭包引用
而作为闭包引用，在defer函数执行时根据整个上下文确定当前的值

看个例子
```golang
func main() {
    i := 0
    defer fmt.Println("a:", i)
    //闭包调用，将外部i传到闭包中进行计算，不会改变i的值，如上边的例3
    defer func(i int) {
        fmt.Println("b:", i)
    }(i)
    //闭包调用，捕获同作用域下的i进行计算
    defer func() {
        fmt.Println("c:", i)
    }()
    i++
}
```
结果：
```sh
c: 1
b: 0
a: 0
```

