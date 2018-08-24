// 安装
// ============
// go get github.com/go-sql-driver/mysql


// 使用
// 用法非常简单
// ============
package main;
 
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "time"
    _ "reflect"
)

/*
先建一张表
CREATE TABLE IF NOT EXISTS user(
    id int(10) UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    username VARCHAR(24) NOT NULL UNIQUE,
    password VARCHAR(8) NOT NULL DEFAULT '',
    age tinyint(3) UNSIGNED NOT NULL DEFAULT 0,
    create_at INT(10) UNSIGNED NOT NULL DEFAULT 0
)ENGINE=InnoDB DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
*/

// var (
//     db *sql.DB
//     err error
// )

var cccc string

func init() {
    cccc = "123456"
}
 
func main() {
    
    fmt.Println(cccc)

    // 通过sql.Open 打开数据库，此处并不验证是否成功连接数据库，当调用db的方法时，才开始验证，比如下面的db.Ping()
    // DSN数据源支持两种格式
    // 1）用户名:密码@协议(地址:端口)/数据库?charset=utf8
    // 2）用户名@unix(/path/to/socket)/数据库?charset=utf8
    db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/qiantest?charset=utf8");
    if err != nil {
        fmt.Println(err);
    }
    // fmt.Println(reflect.TypeOf(db))

    // 关闭数据库，db会被多个goroutine共享，可以不调用
    defer db.Close();

    err = db.Ping()
    if err != nil {
        // do something here
    }

    // 添加
    // insertSql(db, "张三")
    // insertSql(db, "李四")


    // 修改
    // updateSql(db, "张三101", 1)
    // updateSql(db, "张三202", 1)


    // 查询
    selectSql(db)


    // 查单条
    // selectOne(db, 2)


    // 删除
    // deleteSql(db, 2)
    // selectSql(db)
    
}


// 插入数据
// =================
func insertSql(db *sql.DB, name string) {
    // fmt.Println("abc", time.Now().Unix())

    stmt, _ := db.Prepare("insert into user values(null, ?, ?, ?, ?)")
    res, _ := stmt.Exec(name, 123, 27, time.Now().Unix())
    // 注意LastInsertId()是在表中有自增id时，才会有值
    id, err := res.LastInsertId()
    if err != nil {
        panic(err)
    }
    fmt.Println(id)
    // fmt.Println(res.RowsAffected())

    stmt.Close()
}

// 修改数据
// =================
func updateSql(db *sql.DB, name string, id int) {
    // fmt.Println("abc", time.Now().Unix())

    stmt, _ := db.Prepare("update user set username=? where id=?")
    res, _ := stmt.Exec(name, 1)
    // 注意LastInsertId()是在表中有自增id时，才会有值
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }
    fmt.Println(rowsAffected)

    stmt.Close()
}

// 查询列表数据
// =================
func selectSql(db *sql.DB) {
    stmt, _ := db.Prepare("select id, username from user")
    // stmt, _ := db.Prepare("select * from user") // 如果用select *，最终数据Scan时，需要处理所有字段
    rows, _ := stmt.Query()
    // 如果没有参数，也可以不用prepare返回stmt，而直接用db.Query()。即上面两句可用下面一句代替
    // rows, _ := db.Query("select id, username from user")

    for rows.Next() {
        var id int
        var username string
        err := rows.Scan(&id, &username)
        if err != nil {
            panic(err)
        }
        fmt.Println(id, username)    
    }
    
    stmt.Close()
}

// 查询单条数据
// =================
func selectOne(db *sql.DB, id int) {
    stmt, _ := db.Prepare("select username from user where id = ?")
    var username string
    stmt.QueryRow(id).Scan(&username)
    stmt.Close()
}

// 删除数据
// =================
func deleteSql(db *sql.DB, id int) {
    // fmt.Println("abc", time.Now().Unix())

    stmt, _ := db.Prepare("delete from user where id=?")
    res, _ := stmt.Exec(id)
    // 注意LastInsertId()是在表中有自增id时，才会有值
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }
    fmt.Println(rowsAffected)

    stmt.Close()
}

















