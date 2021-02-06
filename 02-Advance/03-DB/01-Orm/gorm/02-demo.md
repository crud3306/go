
下载扩展
```sh
go get github.com/go-sql-driver/mysql 
go get github.com/jinzhu/gorm
go get github.com/gin-gonic/gin
```

建表语句
```sql
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `password` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
```

结构
```sh
├──api
│    ├── apis  
│    │    └── user.go                   
│    ├── database 
│    │    └── mysql.go          
│    ├── models 
│    │    └── user.go       
│    └── router 
│         └──  router.go
└──main.go
```

apis/apis/user.go
```golang
package apis

import (
    "github.com/gin-gonic/gin"
    model "api/models"
    "net/http"
    "strconv"
)

//列表数据
func Users(c *gin.Context) {
    var user model.User
    user.Username = c.Request.FormValue("username")
    user.Password = c.Request.FormValue("password")
    result, err := user.Users()

    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "code":    -1,
            "message": "抱歉未找到相关信息",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": 1,
        "data":   result,
    })
}

//添加数据
func Store(c *gin.Context) {
    var user model.User
    user.Username = c.Request.FormValue("username")
    user.Password = c.Request.FormValue("password")
    id, err := user.Insert()

    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "code":    -1,
            "message": "添加失败",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":  1,
        "message": "添加成功",
        "data":    id,
    })
}

//修改数据
func Update(c *gin.Context) {
    var user model.User
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    user.Password = c.Request.FormValue("password")
    result, err := user.Update(id)

    if err != nil || result.ID == 0 {
        c.JSON(http.StatusOK, gin.H{
            "code":    -1,
            "message": "修改失败",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":  1,
        "message": "修改成功",
    })
}

//删除数据
func Destroy(c *gin.Context) {
    var user model.User
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    result, err := user.Destroy(id)

    if err != nil || result.ID == 0 {
        c.JSON(http.StatusOK, gin.H{
            "code":    -1,
            "message": "删除失败",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "code":  1,
        "message": "删除成功",
    })
}
```


database/mysql.go
```golang
package database

import (
    _ "github.com/go-sql-driver/mysql" //加载mysql
    "github.com/jinzhu/gorm"
    "fmt"
)

var Eloquent *gorm.DB

func init() {
    var err error
    Eloquent, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local&timeout=10ms")

    if err != nil {
        fmt.Printf("mysql connect error %v", err)
    }

    if Eloquent.Error != nil {
        fmt.Printf("database error %v", Eloquent.Error)
    }
}
```


models/user.go
```golang
package models

import (
    orm "api/database"
)

type User struct {
    ID       int64  `json:"id"`       // 列名为 `id`
    Username string `json:"username"` // 列名为 `username`
    Password string `json:"password"` // 列名为 `password`
}

var Users []User

//添加
func (user User) Insert() (id int64, err error) {

    //添加数据
    result := orm.Eloquent.Create(&user)
    id = user.ID
    if result.Error != nil {
        err = result.Error
        return
    }

    return
}

//列表
func (user *User) Users() (users []User, err error) {
    if err = orm.Eloquent.Find(&users).Error; err != nil {
        return
    }
    return
}

//修改
func (user *User) Update(id int64) (updateUser User, err error) {

    if err = orm.Eloquent.Select([]string{"id", "username"}).First(&updateUser, id).Error; err != nil {
        return
    }

    //参数1:是要修改的数据
    //参数2:是修改的数据
    if err = orm.Eloquent.Model(&updateUser).Updates(&user).Error; err != nil {
        return
    }
    return
}

//删除数据
func (user *User) Destroy(id int64) (Result User, err error) {

    if err = orm.Eloquent.Select([]string{"id"}).First(&user, id).Error; err != nil {
        return
    }

    if err = orm.Eloquent.Delete(&user).Error; err != nil {
        return
    }
    Result = *user
    return
}
```


router/router.go
```golang
package router

import (
    "github.com/gin-gonic/gin"
    . "api/apis"
)

func InitRouter() *gin.Engine {

    router := gin.Default()

    router.GET("/users", Users)

    router.POST("/user", Store)

    router.PUT("/user/:id", Update)

    router.DELETE("/user/:id", Destroy)

    return router
}
```

main.go
```golang
package main

import (
    "api/router"
    orm "api/database"
)

func main() {
    defer orm.Eloquent.Close()

    router := router.InitRouter()
    router.Run(":8000")
}
````

执行   
> go run main.go

访问地址
```sh
POST localhost:8006/user 添加
GET localhost:8006/users 列表
DELETE localhost:8006/user/id 删除
PUT localhost:8006/user/id 修改
```
