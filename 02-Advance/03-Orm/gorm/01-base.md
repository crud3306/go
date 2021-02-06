


使用Go自带的"database/sql"数据库连接api，"github.com/go-sql-driver/mysql"MYSQL驱动，通过比较原生的写法去写sql和处理事务。目前开源界也有很多封装好的orm操作框架，帮我们简省一些重复的操作，提高代码可读性。


安装
===============
gorm框架只是简单封装了数据库的驱动包，在安装时仍需要下载原始的驱动包
```
# 由于在项目中使用mysql比较多，这里使用mysql进行数据存储$ 
go get -u github.com/jinzhu/gorm$ 
go get -u github.com/go-sql-driver/mysql
```


CRUD 使用
===============
下面我们使用一张User表来就CRUD做一些操作示例：

表结构如下：
```sql
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL DEFAULT '',
  `age` int(3) NOT NULL DEFAULT '0',
  `sex` tinyint(3) NOT NULL DEFAULT '0',
  `phone` varchar(40) NOT NULL DEFAULT '',
  `email` varchar(40) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;
```

首先初始化数据库连接：
```golang
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"time"
)

var db *gorm.DB

type User struct {
	Id int
	Name string
	Age int
	Sex byte
	Phone string
	Email string
	Birthday *time.Time    
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//设置全局表名禁用复数
	db.SingularTable(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

// 一些基础操作
func baseSql() {

	// 新建
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Create(&user)
	fmt.Println(user.Id) // 自增id在这里取

	db.Create(&User{Name:"bgbiao", Age:18, Email:"bgbiao@bgbiao.top"})


	// ========================
	// 查询
	var u []User
	// Select指定查询字段
	// Where查询条件，可Struct、map或单独字段条件
	// Order排序
	// Limit限制条数、Offset从哪条开始
	db.Select("name,age").Where(&User{Age:12,Sex:1}).Order("age desc").Limit(10).Offset(300).Find(&u)
	db.Select("name,age").Where(map[string]interface{}{"age":12,"sex":1}).Order("age desc").Limit(10).Offset(300).Find(&u)
	db.Where("age > ?",12).Or("sex = ?",1).Order("age desc").Limit(10).Offset(300).Find(&u)

	//count(*)
	var count int
	db.Table("user").Where("age > ?", 0).Count(&count)

	// group by
	type Result struct {  
		Date  time.Time  
		Total int64
	}
	db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)



    // ========================
	// 更新
	user := User{Id: 1,Name:"xiaoming"}
	db.Model(&user).Update(user)
	db.Model(&User{}).Where("sex = ?",1).Update("name","xiaohong")
	db.Model(&User{}).Where("sex = ?",1).Update(&User{Id: 1,Name:"xioaming",Age:12})

	// gorm默认不更新结构体的空值，如果你想手动将某个字段set为空值, 可以使用单独选定某些字段的方式来更新：
	user := User{Id: 1}
	db.Model(&user).Select("name").Update(map[string]interface{}{"name":"","age":0})



    // ========================
	// 删除
	//delete from user where id=1;
	user := User{Id: 1}
	db.Delete(&user)

	db.Where("id = ?", 20).Delete(&User{})
}
```




下面所有的操作都是在上面的初始化连接上执行的操作。



检查表是否存在
-----------
```
// 检查模型是否存在
db.HasTable(&User{})

// 检查表是否存在
db.HasTable("users")
```


插入
-----------
```golang
//插入数据
func (user *User) Insert()  {
	//这里使用了Table()函数，如果你没有指定全局表名禁用复数，或者是表名跟结构体名不一样的时候
	//你可以自己在sql中指定表名。这里是示例，本例中这个函数可以去除。
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	if err := db.Table("user").Create(&user).Error; err != nil {
		return err
	}

	return user.id
}
```



更新
-----------
```golang
//注意，Model方法必须要和Update方法一起使用
//使用效果相当于Model中设置更新的主键key（如果没有where指定，那么默认更新的key为id），Update中设置更新的值
//如果Model中没有指定id值，且也没有指定where条件，那么将更新全表
//相当于：update user set name='xiaoming' where id=1;
user := User{Id: 1,Name:"xiaoming"}
db.Model(&user).Update(user)

//注意到上面Update中使用了一个Struct，你也可以使用map对象。
//需要注意的是：使用Struct的时候，只会更新Struct中这些非空的字段。
//对于string类型字段的""，int类型字段0，bool类型字段的false都被认为是空白值，不会去更新表

//下面这个更新操作只使用了where条件没有在Model中指定id
//update user set name='xiaohong' wehre sex=1
db.Model(&User{}).Where("sex = ?",1).Update("name","xiaohong")


//如果你想手动将某个字段set为空值, 可以使用单独选定某些字段的方式来更新：这个验证后不可用，需用map代替struct才行
user := User{Id: 1}
db.Model(&user).Select("name").Update(map[string]interface{}{"name":"","age":0})
// 提示： 通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体。

//忽略掉某些字段：

//当你的更新的参数为结构体，而结构体中某些字段你又不想去更新，那么可以使用Omit方法过滤掉这些不想update到库的字段：
user := User{Id: 1,Name:"xioaming",Age:12}
db.Model(&user).Omit("name").Update(&user)
```



删除
-----------
```golang
//delete from user where id=1;
user := User{Id: 1}
db.Delete(&user)

//delete from user where id > 11;
db.Delete(&User{},"id > ?",11)
```



事务
-----------
```golang
func CreateAnimals(db *gorm.DB) err {
  tx := db.Begin()
  // 注意，一旦你在一个事务中，使用tx作为数据库句柄

  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
     tx.Rollback()
     return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
     tx.Rollback()
     return err
  }

  tx.Commit()
  return nil
}
```



查询：
-----------
```golang
func (user *User) query() (u []User) {
	//查询所有记录
	db.Find(&u)

	//Find方法可以带 where 参数
	db.Find(&u,"id > ? and age > ?",2,12)

	//带where 子句的查询，注意where要在find前面
	db.Where("id > ?", 2).Find(&u)

	// where name in ("xiaoming","xiaohong")
	db.Where("name in (?)", []string{"xiaoming","xiaohong"}).Find(&u)

	//获取第一条记录，按照主键顺序排序
	db.First(&u)

	//First方法可以带where 条件
	db.First(&u,"where sex = ?",1)

	// LIKE
	db.Where("name LIKE ?", "%jin%").Find(&users)


	// BETWEEN
	db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)

	//获取最后一条记录，按照主键顺序排序
	//同样 last方法也可以带where条件
	db.Last(&u)

	return u
}

//注意：方法中带的&u表示是返回值用u这个对象来接收。
```


上面的查询都将返回表中所有的字段，如果你想指定查询某些字段该怎么做呢？


指定查询字段-Select
```golang
//指定查询字段
db.Select("name,age").Where(map[string]interface{}{"age":12,"sex":1}).Find(&u)
使用Struct和map作为查询条件
//使用Struct，相当于：select * from user where age =12 and sex=1
db.Where(&User{Age:12,Sex:1}).Find(&u)

//等同上一句
db.Where(map[string]interface{}{"age":12,"sex":1}).Find(&u)
```

not 条件的使用
```golang
//where name not in ("xiaoming","xiaohong")
db.Not("name","xiaoming","xiaohong").Find(&u)

//同上
db.Not("name",[]string{"xiaoming","xiaohong"}).Find(&u)
```


or 的使用
```golang
//where age > 12 or sex = 1
db.Where("age > ?",12).Or("sex = ?",1).Find(&u)
```

order by 的使用
```golang
//order by age desc
db.Where("age > ?",12).Or("sex = ?",1).Order("age desc").Find(&u)
```

limit 的使用
```golang
//limit 10
db.Not("name",[]string{"xiaoming","xiaohong"}).Limit(10).Find(&u)
```

offset 的使用
```golang
//limit 300,10
db.Not("name",[]string{"xiaoming","xiaohong"}).Limit(10).Offset(300).Find(&u)
```

count
```golang
//count(*)
var count int
db.Table("user").Where("age > ?", 0).Count(&count)
//注意：这里你在指定表名的情况下sql为：select count(*) from user where age > 0;


//如上代码如果改为：
var count int
var user []User
db.Where("age > ?", 0).Find(&user).Count(&count)
//相当于你先查出来[]User，然后统计这个list的长度。跟你预期的sql不相符。
```

group having
```golang
rows, _ := db.Table("user").Select("count(*),sex").Group("sex").Having("age > ?", 10).Rows()
for rows.Next() {
    fmt.Print(rows.Columns())
}

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
for rows.Next() {  
	...
}

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
for rows.Next() {
  ...
}


type Result struct {  
	Date  time.Time  
	Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

join
```golang
db.Table("user u").Select("u.name,u.age").Joins("left join user_ext ue on u.user_id = ue.user_id").Row()
//如果有多个连接，用多个Join方法即可。

rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
  ...
}

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
// 多连接及参数db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```



原生函数
```golang
db.Exec("DROP TABLE user;")
db.Exec("UPDATE user SET name=? WHERE id IN (?)", "xiaoming", []int{11,22,33})
db.Exec("select * from user where id > ?",10).Scan(&user)
```

一些函数
FirstOrInit 和 FirstOrCreate

获取第一个匹配的记录，若没有，则根据条件初始化一个新的记录：
```golang
//注意：where条件只能使用Struct或者map。如果这条记录不存在，那么会新增一条name=xiaoming的记录
db.FirstOrInit(&u,User{Name:"xiaoming"})
//同上
db.FirstOrCreate(&u,User{Name:"xiaoming"})
```


Attrs

如果没有找到记录，则使用Attrs中的数据来初始化一条记录：
```golang
//使用attrs来初始化参数，如果未找到数据则使用attrs中的数据来初始化一条
//注意：attrs 必须 要和FirstOrInit 或者 FirstOrCreate 连用
db.Where(User{Name:"xiaoming"}).Attrs(User{Name:"xiaoming",Age:12}).FirstOrInit(&u)
```


Assign
```golang
//不管是否找的到，最终返回结构中都将带上Assign指定的参数
db.Where("age > 12").Assign(User{Name:"xiaoming"}).FirstOrInit(&u)
```


Pluck：查询指定的单列

如果user表中你只想查询age这一列，该怎么返回呢，gorm提供了Pluck函数用于查询单列，返回数组：
```golang
var ages []int64
db.Find(&users).Pluck("age", &ages)

var names []string
db.Model(&User{}).Pluck("name", &names)
db.Table("deleted_users").Pluck("name", &names)
```


Scan

Scan函数可以将结果转存储到另一个结构体中。
```golang
type SubUser struct{
    Name string
    Age int
}

db.Table("user").Select("name,age").Where("name = ?", "Antonio").Scan(&SubUser)

// 原生 SQL
db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
```

sql.Row & sql.Rows  
row和rows用户获取查询结果。
```golang
//查询一行
row := db.Table("user").Where("name = ?", "xiaoming").Select("name, age").Row() // (*sql.Row)
//获取一行的结果后，调用Scan方法来将返回结果赋值给对象或者结构体
row.Scan(&name, &age)

//查询多行
rows, err := db.Model(&User{}).Where("sex = ?",1).Select("name, age, phone").Rows() // (*sql.Rows, error)
defer rows.Close()
for rows.Next() {
    ...
    rows.Scan(&name, &age, &email)
    ...
}
```




日志#
-----------
Gorm有内置的日志记录器支持，默认情况下，它会打印发生的错误。
```golang
// 启用Logger，显示详细日志
db.LogMode(true)

// 禁用日志记录器，不显示任何日志
db.LogMode(false)

// 调试单个操作，显示此操作的详细日志
db.Debug().Where("name = ?", "xiaoming").First(&User{})
```















