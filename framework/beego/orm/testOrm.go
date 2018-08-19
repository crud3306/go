package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
)

/**
开发时，可以不用整个beego框架，直接用beego中的某项功能，比如这里只用了orm
*/

// type User 操作的表是user
// type UserInfo 操作的表是 user_info
type User struct {
	Id int64
	Username string
	Password string
	Age int
	CreateAt int64
}

func main() {
	// 打开调试，可看到执行的相过过程，含sql语句
	orm.Debug = true

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/qiantest?charset=utf8")

	orm.RegisterModel(new(User))

	o := orm.NewOrm()


	// 新增
	// ========================
	// user := User{Username:"zhang", Password:"haha123", Age:16, CreateAt:1234567}
	// id, err := o.Insert(&user)
	// fmt.Printf("id: %d; err: %v", id, err)


	// 更新
	// ========================
	// user := User{Username:"lishi", Password:"haha123", Age:16, CreateAt:1234567}
	// user.Id  = 1
	// user.Username = "lishi2"
	// num, err := o.Update(&user)
	// fmt.Printf("num: %d; err: %v", num, err)	


	// 删除
	// ========================
	// user.Id  = 1
	// num, err := o.Delete(&user)
	// fmt.Printf("num: %d; err: %v", num, err)	



	// 查询单条
	// ========================
	// user := User{}
	// user.Id  = 2
	// err := o.Read(&user)
	// fmt.Printf("user: %v, err: %v \n", user, err)	


	// 原生sql查询
	// ========================
	// var maps []orm.Params
	// _, err1 := o.Raw("SELECT * FROM user").Values(&maps)
	// fmt.Println(err1)
	// for _, term := range maps{
	// 	fmt.Println(term["id"], ":", term["username"])
	// }


	// 事务
	// ========================
	// o.Begin()

	// user := User{Username:"haha"}
	// _, err := o.Insert(&user)
	// if err == nil {
	// 	o.Commit()
	// } else {
	// 	o.Rollback()
	// }


	// 采用queryBuilder方式进行读取
	// ========================
	// 参考地址：https://beego.me/docs/mvc/model/querybuilder.md  
	var users []User
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("password").From("user").Where("age > ?").Limit(1)

	sql := qb.String()
	o.Raw(sql, 1).QueryRows(&users)
	fmt.Printf("user: %v", users)

}


/*
type QueryBuilder interface {
    Select(fields ...string) QueryBuilder
    From(tables ...string) QueryBuilder
    InnerJoin(table string) QueryBuilder
    LeftJoin(table string) QueryBuilder
    RightJoin(table string) QueryBuilder
    On(cond string) QueryBuilder
    Where(cond string) QueryBuilder
    And(cond string) QueryBuilder
    Or(cond string) QueryBuilder
    In(vals ...string) QueryBuilder
    OrderBy(fields ...string) QueryBuilder
    Asc() QueryBuilder
    Desc() QueryBuilder
    Limit(limit int) QueryBuilder
    Offset(offset int) QueryBuilder
    GroupBy(fields ...string) QueryBuilder
    Having(cond string) QueryBuilder
    Subquery(sub string, alias string) string
    String() string
}
*/





























