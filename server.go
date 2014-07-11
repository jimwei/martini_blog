package main

import (
	"bytes"
	_ "code.google.com/p/odbc"
	_ "code.google.com/p/odbc/api"
	"database/sql"
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	"log"
)

var (
	mssrv    = flag.String("mssrv", "xa-lsr-jimweiw7", "server name")
	msdb     = flag.String("msdb", "MartiniBlog", "database")
	msuser   = flag.String("user", "sa", "user name")
	mspwd    = flag.String("pwd", "xA123456", "password")
	msdriver = flag.String("msdriver", "sql server", "msdriver")
)

type Blog struct {
	Id          int
	Title       string
	Date        string
	Description template.HTML
	Author      int
}

var db *sql.DB = nil

//创建数据库连接字符串
func CreateConnection() (db *sql.DB, err error) {
	sqlConnString := getSqlConnectionString()
	db, err = sql.Open("odbc", sqlConnString)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db, nil
}

//构造连接字符串
func getSqlConnectionString() string {
	buf := bytes.NewBufferString("")
	buf.WriteString("Driver=%s;")
	buf.WriteString("server=%s;")
	buf.WriteString("database=%s;")
	buf.WriteString("user id=%s;")
	buf.WriteString("pwd=%s;")

	connString := fmt.Sprintf(buf.String(), *msdriver, *mssrv, *msdb, *msuser, *mspwd)

	return connString
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	db, _ = CreateConnection()
	defer db.Close()
	//获取全部的blog
	m.Get("/", getHandle)
	//获取指定的blog
	m.Get("/blog/:id", GetBlogByID)

	m.Run()
}

/*
*获取全部的blog
 */
func getHandle(r render.Render) {
	var blogs []Blog
	rows, err := db.Query("select * from blogs")
	if err != nil {
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer rows.Close()

	for rows.Next() {
		var blog Blog
		var description string
		err := rows.Scan(&blog.Id, &blog.Title, &blog.Date, &description, &blog.Author)
		if err != nil {
			log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		blog.Description = template.HTML(description)
		blogs = append(blogs, blog)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	r.HTML(200, "home", blogs)
}

/*
*获取指定id的blog
 */
func GetBlogByID(params martini.Params, r render.Render) {
	var blog Blog
	fmt.Println("the params id is", params["id"])
	rows, err := db.Query("select * from blogs where id = ?", params["id"])
	if err != nil {
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer rows.Close()

	for rows.Next() {
		var description string
		err := rows.Scan(&blog.Id, &blog.Title, &blog.Date, &description, &blog.Author)
		if err != nil {
			log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		blog.Description = template.HTML(description)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	r.HTML(200, "blog", blog)
}
