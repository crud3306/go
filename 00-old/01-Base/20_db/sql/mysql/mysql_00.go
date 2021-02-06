
// 安装
// ============
// go get github.com/go-sql-driver/mysql


// 使用
// ============
package main;
 
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)
 
func main() {
    //打开数据库
    //DSN数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
    db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/test?charset=utf8");
    if err != nil {
        fmt.Println(err);
    }

    //关闭数据库，db会被多个goroutine共享，可以不调用
    defer db.Close();

    err = db.Ping()
	if err != nil {
	    // do something here
	}


    //查询数据，指定字段名，返回sql.Rows结果集
    rows, _ := db.Query("select id,name from test");
    id := 0;
    name := "";
    for rows.Next() {
        rows.Scan(&id, &name);
        fmt.Println(id, name);
    }
 

    //查询数据，取所有字段
    rows2, _ := db.Query("select * from test");
    //返回所有列
    cols, _ := rows2.Columns();
    //这里表示一行所有列的值，用[]byte表示
    vals := make([][]byte, len(cols));
    //这里表示一行填充数据
    scans := make([]interface{}, len(cols));
    //这里scans引用vals，把数据填充到[]byte里
    for k, _ := range vals {
        scans[k] = &vals[k];
    }
 
    i := 0;
    result := make(map[int]map[string]string);
    for rows2.Next() {
        //填充数据
        rows2.Scan(scans...);
        //每行数据
        row := make(map[string]string);
        //把vals中的数据复制到row中
        for k, v := range vals {
            key := cols[k];
            //这里把[]byte数据转成string
            row[key] = string(v);
        }
        //放入结果集
        result[i] = row;
        i++;
    }
    fmt.Println(result);
 

    //查询一行数据
    rows3 := db.QueryRow("select id,name from test where id = ?", 1);
    rows3.Scan(&id, &name);
    fmt.Println(id, name);
 

    //插入一行数据
    ret, _ := db.Exec("insert into test(id,name) values(null, '444')");
    //获取插入ID
    ins_id, _ := ret.LastInsertId();
    fmt.Println(ins_id);
 

    //更新数据
    ret2, _ := db.Exec("update test set name = '000' where id > ?", 2);
    //获取影响行数
    aff_nums, _ := ret2.RowsAffected();
    fmt.Println(aff_nums);
 

    //删除数据
    ret3, _ := db.Exec("delete from test where id = ?", 3);
    //获取影响行数
    del_nums, _ := ret3.RowsAffected();
    fmt.Println(del_nums);
 

    //预处理语句
    stmt, _ := db.Prepare("select id,name from test where id = ?");
    rows4, _ := stmt.Query(3);
    //注意这里需要Next()下，不然下面取不到值
    rows4.Next();
    rows4.Scan(&id, &name);
    fmt.Println(id, name);
 
    stmt2, _ := db.Prepare("insert into test values(null, ?, ?)");
    rows5, _ := stmt2.Exec("666", 66);
    fmt.Println(rows5.RowsAffected());
 

    //事务处理
    tx, _ := db.Begin();
 
    ret4, _ := tx.Exec("update test set price = price + 100 where id = ?", 1);
    ret5, _ := tx.Exec("update test set price = price - 100 where id = ?", 2);
    upd_nums1, _ := ret4.RowsAffected();
    upd_nums2, _ := ret5.RowsAffected();
 
    if upd_nums1 > 0 && upd_nums2 > 0 {
        //只有两条更新同时成功，那么才提交
        tx.Commit();
    } else {
        //否则回滚
        tx.Rollback();
    }
    
}



