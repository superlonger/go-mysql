# go-mysql
Most simple mysql package by go lang.  
最简单的golang mysql包
### Installation

go get github.com/superlonger/go-mysql

### Usage
All you need is: Import, Init, Select or  Update(including Insert, Delete, Update)
```go
import "github.com/superlonger/go-mysql"

// init database, and test connect.
err := mysql.Init("testdb_root:test_password@tcp(127.0.0.1:3306)/testdb?charset=utf8")
// close DB when you exit
defer mysql.Close()

// insert data to database, return the last insert id
_, lastId, err := mysql.Exec("insert into userinfo set username=?,password=?", "test_username", "test_password")

// update data to database, return the number of rows affected
affect, _, err := mysql.Exec("update userinfo set logintime=CURRENT_TIMESTAMP where username=? and password=?", "test_username", "test_password")

// delete data to database, return the number of rows affected
affect, _, err = mysql.Exec("delete from userinfo where username=?", "test_username")

// select datas from  database, return [][]interface{} slice
datas, err := mysql.Query("select * from userinfo where id>=?", 5)
```
There is a simple example 'test.go' in the dir 'examples'

### Tips
This package is build base on [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
