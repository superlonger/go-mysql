package main

import (
	"fmt"
	"github.com/superlonger/go-mysql"
	"strconv"
)

func main() {
	// init database, and test connect.
	err := mysql.Init("testdb_root:test_password@tcp(127.0.0.1:3306)/testdb?charset=utf8")
	if err != nil {
		fmt.Printf("[Init]: %s\n", err)
		return
	}
	defer mysql.Close()

	// insert data to database, return the last insert id
	_, lastId, err := mysql.Exec("insert into userinfo set username=?,password=?", "test_username", "test_password")
	if err != nil {
		fmt.Printf("[Insert]->%s", err)
		return
	}
	fmt.Printf("Last ID: %d\n", lastId)

	// update data to database, return the number of rows affected
	affect, _, err := mysql.Exec("update userinfo set logintime=CURRENT_TIMESTAMP where username=? and password=?", "test_username", "test_password")
	if err != nil {
		fmt.Printf("[Update]->%s", err)
		return
	}
	if affect != 1 {
		fmt.Println("username or password invalid")
	}

	// delete data to database, return the number of rows affected
	affect, _, err = mysql.Exec("delete from userinfo where username=?", "test_username")
	if err != nil {
		fmt.Printf("[Delete]->%s", err)
		return
	}
	fmt.Printf("%d rows affected\n", affect)

	// select datas from  database, return [][]interface{} slice
	datas, err := mysql.Query("select * from userinfo where id>=?", 5)
	if err != nil {
		fmt.Printf("Select error: %s\n", err)
		return
	}
	nRows := len(datas)
	if nRows < 1 {
		fmt.Println("no datas selected.")
		return
	}
	fmt.Printf("%d rows, %d cols\n", nRows, len(datas[0]))
	for i := range datas {
		for j := range datas[i] {
			fmt.Print(datas[i][j])
		}
		fmt.Println()
	}
	// convert data to struct
	// use 'reflect' can convert automatic, but is inefficient, user should handling this by self.
	type UserInfo struct {
		Id         int64
		Username   string
		Password   string
		Createtime string
		Logintime  string
	}
	for _, data := range datas {
		userInfo := UserInfo{}
		userInfo.Id, err = strconv.ParseInt(string(data[0]), 10, 64)
		if err != nil {
			fmt.Println("strconv.ParseInt error")
			return
		}
		userInfo.Username = string(data[1])
		userInfo.Password = string(data[2])
		userInfo.Createtime = string(data[3])
		if data[4] == nil {
			userInfo.Logintime = "NULL"
		} else {
			userInfo.Logintime = string(data[4])
		}
		fmt.Printf("%v\n", userInfo)
	}
}
