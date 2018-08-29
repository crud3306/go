快速入门地址：
----------
http://www.runoob.com/go/go-tutorial.html  
  
  
快速开始，hello.go
----------
```go
package main

import "fmt"

func main() {
   /* 这是我的第一个简单的程序 */
   fmt.Println("Hello, World!")
}
```
执行
> go run hello.go  
  

注释  
----------
// 单行注释  
/*  
 Author by 菜鸟教程  
 我是多行注释  
 */  
  
  
数据类型  
----------
1	布尔型  
布尔型的值只可以是常量 true 或者 false。  
一个简单的例子：var b bool = true。  
  
2	数字类型  
整型int 和浮点型float32、float64，Go语言支持整型和浮点型数字，并且原生支持复数，其中位的运算采用补码。  

3	字符串类型:  
字符串就是一串固定长度的字符连接起来的字符序列。Go的字符串是由单个字节连接起来的。Go语言的字符串的字节使用UTF-8编码标识Unicode文本。  
  
4	派生类型:  
(a) 指针类型（Pointer）  
(b) 数组类型  
(c) 结构化类型(struct)  
(d) Channel 类型  
(e) 函数类型  
(f) 切片类型  
(g) 接口类型（interface）  
(h) Map 类型  
   
  
  
变量声明、赋值  
----------
python声明变量时可以指定变量类型，也可以不指定

第一种，指定变量类型，声明后若不赋值，使用默认值。   
> var v_name v_type  
> v_name = value  

第二种，根据值自行判定变量类型。  
> var v_name = value  
  
第三种，不用var，用 :=  ，注意这种方式只能在函数体里面使用
> v_name := value  


例：
```go
var a int = 10
var b = 10
c := 10
```
```go
package main

var x, y int
var (  // 这种因式分解关键字的写法一般用于声明全局变量
    a int
    b bool
)

var c, d int = 1, 2
var e, f = 123, "hello"

// 这种不带声明格式的只能在函数体中出现
// g, h := 123, "hello"

func main(){
    g, h := 123, "hello"
    fmt.Println(x, y, a, b, c, d, e, f, g, h)
}

// 以上实例执行结果为：
0 0 0 false 1 2 123 hello 123 hello
```
  
  
运算符  
----------
Go 语言内置的运算符有：  
  
算术运算符   
假定 A 值为 10，B 值为 20。  
```
+	相加	A + B 输出结果 30
-	相减	A - B 输出结果 -10
*	相乘	A * B 输出结果 200
/	相除	B / A 输出结果 2
%	求余	B % A 输出结果 0
++	自增	A++ 输出结果 11
--	自减	A-- 输出结果 9
```
  
关系运算符  
假定 A 值为 10，B 值为 20。  
```
==	检查两个值是否相等，如果相等返回 True 否则返回 False。	(A == B) 为 False
!=	检查两个值是否不相等，如果不相等返回 True 否则返回 False。	(A != B) 为 True
>	检查左边值是否大于右边值，如果是返回 True 否则返回 False。	(A > B) 为 False
<	检查左边值是否小于右边值，如果是返回 True 否则返回 False。	(A < B) 为 True
>=	检查左边值是否大于等于右边值，如果是返回 True 否则返回 False。	(A >= B) 为 False
<=	检查左边值是否小于等于右边值，如果是返回 True 否则返回 False。	(A <= B) 为 True
```
  
  
逻辑运算符  
假定 A 值为 True，B 值为 False。  
```
&&	逻辑 AND 运算符。 如果两边的操作数都是 True，则条件 True，否则为 False。	(A && B) 为 False
||	逻辑 OR 运算符。 如果两边的操作数有一个 True，则条件 True，否则为 False。	(A || B) 为 True
!	逻辑 NOT 运算符。 如果条件为 True，则逻辑 NOT 条件 False，否则为 True。	!(A && B) 为 True
```
  
  
位运算符  
位运算符对整数在内存中的二进制位进行操作。  
下表列出了位运算符 &, |, 和 ^ 的计算：  
```
p	q	p & q	p | q	p ^ q
0	0	0	    0	    0
0	1	0     	1	    1
1	1	1	    1	    0
1	0	0	    1	    1


假定 A 为60，B 为13：

&	按位与运算符"&"是双目运算符。 其功能是参与运算的两数各对应的二进位相与。	(A & B) 结果为 12, 二进制为 0000 1100

|	按位或运算符"|"是双目运算符。 其功能是参与运算的两数各对应的二进位相或	(A | B) 结果为 61, 二进制为 0011 1101

^	按位异或运算符"^"是双目运算符。 其功能是参与运算的两数各对应的二进位相异或，当两对应的二进位相异时，结果为1。	(A ^ B) 结果为 49, 二进制为 0011 0001

<<	左移运算符"<<"是双目运算符。左移n位就是乘以2的n次方。 其功能把"<<"左边的运算数的各二进位全部左移若干位，由"<<"右边的数指定移动的位数，高位丢弃，低位补0。	A << 2 结果为 240 ，二进制为 1111 0000

>>	右移运算符">>"是双目运算符。右移n位就是除以2的n次方。 其功能是把">>"左边的运算数的各二进位全部右移若干位，">>"右边的数指定移动的位数。	A >> 2 结果为 15 ，二进制为 0000 1111  
```
  
  
赋值运算符  
```
=	简单的赋值运算符，将一个表达式的值赋给一个左值	C = A + B 将 A + B 表达式结果赋值给 C
+=	相加后再赋值	C += A 等于 C = C + A
-=	相减后再赋值	C -= A 等于 C = C - A
*=	相乘后再赋值	C *= A 等于 C = C * A
/=	相除后再赋值	C /= A 等于 C = C / A
%=	求余后再赋值	C %= A 等于 C = C % A
<<=	左移后赋值	C <<= 2 等于 C = C << 2
>>=	右移后赋值	C >>= 2 等于 C = C >> 2
&=	按位与后赋值	C &= 2 等于 C = C & 2
^=	按位异或后赋值	C ^= 2 等于 C = C ^ 2
|=	按位或后赋值	C |= 2 等于 C = C | 2
```
    
  
其他运算符   
```
&	返回变量存储地址	&a 将给出变量的实际地址。
*	指针变量 	    *a 是一个指针变量
```
   
  
条件语句
-----------
if 语句  
if 语句 由一个布尔表达式后紧跟一个或多个语句组成。  
  
if...else 语句	  
if 语句 后可以使用可选的 else 语句, else 语句中的表达式在布尔表达式为 false 时执行。  
  
if 嵌套语句  
你可以在 if 或 else if 语句中嵌入一个或多个 if 或 else if 语句。   
  
switch 语句  
switch 语句用于基于不同条件执行不同动作。配合case使用。go的case语句默认会break，所以不用手动加break(当然加了也不为错)。    
  
select 语句  
select 语句类似于switch语句，但是select会随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。select需配合case channel使用。  
  
if 语句
```go
var a int = 100;
 
/* 判断布尔表达式 */
if a < 20 {
   /* 如果条件为 true 则执行以下语句 */
   fmt.Printf("a 小于 20\n" );
} else {
   /* 如果条件为 false 则执行以下语句 */
   fmt.Printf("a 不小于 20\n" );
}
fmt.Printf("a 的值为 : %d\n", a);
``` 
  
  
switch语句
```go
var grade string = "B"
var marks int = 90

switch marks {
  case 90: grade = "A"
  case 80: grade = "B"
  case 50,60,70 : grade = "C"
  default: grade = "D"  
}

switch {
  case grade == "A" :
     fmt.Printf("优秀!\n" )     
  case grade == "B", grade == "C" :
     fmt.Printf("良好\n" )      
  case grade == "D" :
     fmt.Printf("及格\n" )      
  case grade == "F":
     fmt.Printf("不及格\n" )
  default:
     fmt.Printf("差\n" );
}
fmt.Printf("你的等级是 %s\n", grade ); 
```
   
select语句 （select语句需配合channel使用）  
```go
var c1, c2, c3 chan int
var i1, i2 int
select {
  case i1 = <-c1:
     fmt.Printf("received ", i1, " from c1\n")
  case c2 <- i2:
     fmt.Printf("sent ", i2, " to c2\n")
  case i3, ok := (<-c3):  // same as: i3, ok := <-c3
     if ok {
        fmt.Printf("received ", i3, " from c3\n")
     } else {
        fmt.Printf("c3 is closed\n")
     }
  default:
     fmt.Printf("no communication\n")
}   
```
  

循环语句
------------
go 没有while循环，只有for循环

```go
var b int = 15
var a int

numbers := [6]int{1, 2, 3, 5} 

/* for 循环 */
for a := 0; a < 10; a++ {
  fmt.Printf("a 的值为: %d\n", a)
}

for a < b {
  a++
  fmt.Printf("a 的值为: %d\n", a)
}

for i,x:= range numbers {
  fmt.Printf("第 %d 位 x 的值 = %d\n", i,x)
}   

```

```
// 下面的无限循环与其它语言的while(true)效果一样
for true  {
    fmt.Printf("这是无限循环。\n");
}
```


















