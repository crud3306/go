// 地址
// ===========
// 使用说明
// https://studygolang.com/articles/8930  (较详细)
// https://www.jianshu.com/p/5a1712e6141f (简单入门)


// 安装
// ===========
// go get gopkg.in/mgo.v2


// 使用
// ===========
// 基本用法
package main 

import ( "fmt" 
	"labix.org/v2/mgo" 
	"labix.org/v2/mgo/bson" 
) 

type Person struct 
{ 
	Name string  `bson:"name"`
	Phone string `bson:"phone"`
} 

func main() 
{ 
	// 连接服务器  ////////////////////
	// 如果在本机
	session, err := mgo.Dial("") 
	// 或者session, err := mgo.Dial("localhost:27017")

	// 如果不在本机或端口不同，传入相应的地址即可。如：
	//mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	if err != nil { panic(err) } 
	defer session.Close() 

	//  ////////////////////
	session.SetMode(mgo.Monotonic, true)

	// 切换数据库 + 切换集合
	c := session.DB("test").C("people") 

	// 对集合进行操作
	// 添加
	err = c.Insert(&Person{"Ale", "111111"}, &Person{"Cla", "222222222"}) 
	if err != nil { 
		panic(err) 
	} 

	// 查找单条，注意 Find + One
	result := Person{} 
	err = c.Find(bson.M{"name": "Ale"}).One(&result) 
	if err != nil { 
		panic(err) 
	} 
	fmt.Println("Phone:", result.Phone) 

	// 查找多条，注意 Find + All
	var persons []Person
	// err = c.Find(bson.M{"name": "Ale"}).One(&persons) 
	err = c.Find(nil).One(&persons) 
	if err != nil { 
		panic(err) 
	} 
	fmt.Println(persons) 


	// 更新
	err = c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$set": bson.M{ "name": "Jimmy Gu", "age": 34, }})


	// 删除
	err = c.Remove(bson.M{"name": "Jimmy Kuu"})

	// 加索引
	index := mgo.Index{
        Key:        []string{"name"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }
    err := c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }

}
// 更多操作方式、查询条件，见地址：https://studygolang.com/articles/8930

// 封装
// =================
// 可以把mgo的连接库，选择集合封装一下
 func ConnecMongo(cName string, cDb string) *mgo.Collection {
    session, err := mgo.Dial("127.0.0.1:27017")
    if err != nil {
        panic(err)
    }
    //defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    c := session.DB(cDb).C(cName)
    return c
}

// 使用封装的
c := ConnecMongo('test', 'book');
err := c.Update(bson.M{"email": "12832984@qq.com"}, bson.M{"$set": bson.M{"name": "haha", "phone": "37848"}})
if err != nil {
    log.Fatal(err)
}
















