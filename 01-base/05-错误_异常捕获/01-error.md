


错误 error
================

Go 语言通过内置的错误接口提供了非常简单的错误处理机制。

error类型是go语言的一种内置类型，使用的时候不用特定去import，他本质上是一个接口
```golang
type error interface {
	Error() string //Error()是每一个订制的error对象需要填充的错误消息,可以理解成是一个字段Error
}
```


怎样去理解这个订制呢？
我们知道接口这个东西，必须拥有它的实现块才能调用，放在这里就是说，Error()必须得到填充，才能使用.

比方说下面三种方式：
- 通过errors包New()方法去订制error
- 通过fmt.Errorf()去订制
- 通过自定义的MyError块去订制了


第一种:通过errors包去订制error
---------------
```golang
error := errors.New("hello,error")//使用errors必须import "errors"包
if error != nil {
    fmt.Print(err)
}
```

来解释一下errors包，只是一个为Error()填充的简易封装,整个包的内容，只有一个New方法，可以查看源码如下
```golang
package errors

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```



第二种，通过fmt.Errorf()去订制
---------------
```golang
err := fmt.Errorf("hello error")
if err != nil {
    fmt.Print(err)
}
```
可以说和第一种雷同了.



第三种，就是通过自定义的MyError块去订制了
---------------
//一个包裹了错误类型对象的自定义错误类型
```golang
type MyError struct {
	err error 
}

//订制Error()
func (e MyError) Error() string {
    return e.err.Error()
}

func main() {
	err := MyError{
		errors.New("hello error"),
	}
	fmt.Println(err.Error())
}
```
或者
```golang
type MyError struct {
	s string 
}

//订制Error()
func (e *MyError) Error() string {
    return e.s
}

func main() {
	err := MyError{
		"hello error",
	}
	fmt.Println(err.Error())
}
```


三种方式差异都不大,输出结果都是 hello error
实际上error只是一段错误信息，真正抛出异常并不是单纯靠error，panic和recover的用法以后总结