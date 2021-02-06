
sort库


sort简单使用排序
------------
```golang

import （
	"sort"
	"fmt"
）

/*
说明：对于int、float64、string数组/分片的排序，
     go分别提供sort.Ints()、sort.Float64s()、sort.Strings()函数（默认从小->大排序）
*/
func upSort(){
    intList := [] int {2, 4, 3, 5, 7, 6, 9, 8, 1, 0}
    floatList := [] float64 {4.2, 5.9, 12.3, 10.0, 50.4, 99.9, 31.4, 27.81828, 3.14}
    stringList := [] string {"a", "c", "b", "d", "f", "i", "z", "x", "w", "y"}

    sort.Ints(intList)
    sort.Float64s(floatList)
    sort.Strings(stringList)

    fmt.Printf("%v\n%v\n%v\n",intList,floatList,stringList)

    /*
    打印结果：
    [0 1 2 3 4 5 6 7 8 9]
    [3.14 4.2 5.9 10 12.3 27.81828 31.4 50.4 99.9]
    [a b c d f i w x y z]
    */
}


//降序排序
func downSort(){
    intList := [] int {2, 4, 3, 5, 7, 6, 9, 8, 1, 0}
    floatList := [] float64 {4.2, 5.9, 12.3, 10.0, 50.4, 99.9, 31.4, 27.81828, 3.14}
    stringList := [] string {"a", "c", "b", "d", "f", "i", "z", "x", "w", "y"}

    sort.Sort(sort.Reverse(sort.IntSlice(intList)))
    sort.Sort(sort.Reverse(sort.Float64Slice(floatList)))
    sort.Sort(sort.Reverse(sort.StringSlice(stringList)))

    fmt.Printf("%v\n%v\n%v\n", intList, floatList, stringList)
    
    /*
    打印结果：
    [9 8 7 6 5 4 3 2 1 0]
    [99.9 50.4 31.4 27.81828 12.3 10 5.9 4.2 3.14]
    [z y x w i f d c b a]
    */
}
```


对slice struct排序
------------
```golang
package main

import (
	"fmt"
	"sort"
)

type result struct {
	IP    string  `json:"ip"`
	Value float64 `json:"value"`
}

type result2 struct {
	IP    string  `json:"ip"`
	Value float64 `json:"value"`
	Age   int     `json:"age"`
}

func main() {
	//使用 sort 提供的排序接口sort.Slice

	// 单字段排序，Value
	results := []result{result{"ip1", 1.2}, result{"ip2", 1.1}, result{"ip3", 1.3}}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value < results[j].Value
	})

	fmt.Println(results)


	// 多字段排序，Value、Age
	result2s := []result2{result2{"ip1", 1.2, 3}, result2{"ip2", 1.1, 2}, result2{"ip2-1", 1.1, 1}, result2{"ip3", 1.3, 1}, result2{"ip2-2", 1.1, 3}}

	sort.Slice(result2s, func(i, j int) bool {
		if result2s[i].Value < result2s[j].Value {
			return true
		} else if result2s[i].Value == result2s[j].Value {
			if result2s[i].Age < result2s[j].Age {
				return true
			}
		}

		return false
	})

	fmt.Println(result2s)
}

//输出
//[{ip2 1.1} {ip1 1.2} {ip3 1.3}]
//[{ip2-1 1.1 1} {ip2 1.1 2} {ip2-2 1.1 3} {ip1 1.2 3} {ip3 1.3 1}]
```



封装
==============

sort接口封装01
--------------
```golang
/*
(1)模拟IntSlice排序
缺点：
根据 Age 排序需要重新定义 PersonSlice 方法，绑定 Len 、 Less 和 Swap 方法，
如果需要根据 Name 排序， 又需要重新写三个函数； 如果结构体有 4 个字段，有四种类型的排序，那么就要写 3 × 4 = 12 个方法，
即使有一些完全是多余的;
根据不同的标准 Age 或是 Name，真正不同的体现在 Less 方法上，所以可以将 Less 抽象出来，
每种排序的 Less 让其变成动态的.见（2）
*/
type Person struct {
    Name string
    Age int
}

//按照Person.Age从大-》小排序（PersonSlice是person[]的模版）
type PersonSlice [] Person

//重写len()方法
func(a PersonSlice) Len() int{
    return len(a)
}
//重写Swap()方法
func (a PersonSlice) Swap(i,j int){
    a[i],a[j]=a[j],a[i]
}
//重写Less()方法
func (a PersonSlice) Less(i,j int ) bool{
    return a[j].Age < a[i].Age
}

func IntSliceSort(){
    people:=[]Person{
        {"zhang san", 12},
        {"li si", 30},
        {"wang wu", 52},
        {"zhao liu", 26},
    }
    fmt.Println(people) //[{zhang san 12} {li si 30} {wang wu 52} {zhao liu 26}]

    sort.Sort(PersonSlice(people)) //按照 Age 的逆序排序
    fmt.Println(people) //[{wang wu 52} {li si 30} {zhao liu 26} {zhang san 12}]

    sort.Sort(sort.Reverse(PersonSlice(people))) //按照 Age 的升序排序
    fmt.Println(people)//[{zhang san 12} {zhao liu 26} {li si 30} {wang wu 52}]
}
```


sort接口封装02
--------------
```golang
type Person3 struct {
    Name string
    Age  int
}

type PersonWrapper3 struct {
    people [] Person3
    by func(p, q * Person3) bool
}

type SortBy func(p, q *Person3) bool


func (pw PersonWrapper3) Len() int {         // 重写 Len() 方法
    return len(pw.people)
}
func (pw PersonWrapper3) Swap(i, j int){     // 重写 Swap() 方法
    pw.people[i], pw.people[j] = pw.people[j], pw.people[i]
}
func (pw PersonWrapper3) Less(i, j int) bool {    // 重写 Less() 方法
    return pw.by(&pw.people[i], &pw.people[j])
}


// 封装成 SortPerson 方法
func SortPerson(people [] Person3, by SortBy){
    sort.Sort(PersonWrapper3{people, by})
}


func wrapperSorts(){
    people := [] Person3{
        {"zhang san", 12},
        {"li si", 30},
        {"wang wu", 52},
        {"zhao liu", 26},
    }
    fmt.Println(people)

    // 方式1
    sort.Sort(PersonWrapper3{people, func (p, q *Person3) bool {
        return q.Age < p.Age    // Age 递减排序
    }})
    fmt.Println(people)

    // 方式2
    SortPerson(people, func (p, q *Person3) bool {
        return p.Name < q.Name    // Name 递增排序
    })
    fmt.Println(people)

    /*
    运行结果：
    [{zhang san 12} {li si 30} {wang wu 52} {zhao liu 26}]
    [{wang wu 52} {li si 30} {zhao liu 26} {zhang san 12}]
    [{li si 30} {wang wu 52} {zhang san 12} {zhao liu 26}]
    */
}
```


第三方库
==============
github.com/patrickmn/sortutil

```golang
package main

import (
    "fmt"

    "github.com/patrickmn/sortutil"
)

type result struct {
    IP    string  `json:"ip"`
    Value float64 `json:"value"`
    Age   int     `json:"age"`
}

func main() {
    // 这个包好像只能对单字段排序，如:Value
    results := []result{result{"ip1", 1.2, 1}, result{"ip2", 1.1, 3}, result{"ip3", 1.3, 2}}
    sortutil.AscByField(results, "Value")

    fmt.Println(results)
}
```
