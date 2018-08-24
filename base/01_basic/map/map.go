package main

import "fmt"

// map[K]V, map[K1]map[K2]V

// m := map[string]string {
// 	"name": "ccmouse",
// 	"course": "golang",
// 	"site":"imooc",
// 	"quality":"notbad",
// }

func main() {
	m := map[string]string {
		"name": "ccmouse",
		"course": "golang",
		"site":"imooc",
		"quality":"notbad",
	}

	m2 := make(map[string]int)	// m2 == empty map

	var m3 map[string]int	// m3 == nil

	fmt.Println(m, m2, m3)


	// 获取元素个数
	fmt.Println(len(m))


	// 可用range来遍历map，但是不保证遍历顺序；如果需顺序，需手动对key排序
	// key 和 value 都要
	for k, v := range m {
		fmt.Println(k, v)
	}

	// 只要key
	// for k := range m {
	// 	fmt.Println(k)
	// }

	// 只要value
	// for _, v := range m {
	// 	fmt.Println(v)
	// }


	// 获取某个key对应的值
	fmt.Println("test get values by key");
	courseName := m["course"]
	fmt.Println(courseName)

	courseName, ok := m["course"]
	fmt.Println(courseName, ok)

	if courseName, ok := m["course1"]; ok {
		fmt.Println(courseName)
	} else {	
		fmt.Println("key does not exist")
	}

	// 测试删除某个key
	fmt.Println("test delete key");
	name, ok := m["name"]
	fmt.Println(name, ok)

	delete(m, "name")
	name, ok := m["name"]
	fmt.Println(name, ok)


}

















