package main

import (
	"database/sql"
	"flag"
	"html/template"
	"net/http"
	"time"

	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

// Category
type Category struct {
	ID          int
	Name        string
	CreatedTime string
	CreatedBy   string
}

var categorytemplate *template.Template

const (
	dbDrive     = "sqlite3"
	dbName      = "data.db"
	ShortFormat = "2006-01-02"
	LongFormat  = "2006-01-02 15:04:05"
)

func init() {
	flag.Parse()
	categorytemplate = template.Must(template.New("category.gtpl").ParseFiles("./templates/category.gtpl"))
	glog.Infoln("initial done")
}
func main() {
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	err := http.ListenAndServe(":8888", nil)      //设置监听的端口
	if err != nil {
		glog.Errorf("main->ListenAndServe err: %v\n", err)
	}
}

//GetEntity get data from db
func (cate Category) GetEntity() []Category {
	db, err := sql.Open(dbDrive, "./data.db")

	if err != nil {
		glog.Errorf("Category->GetEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT ID, Name,CreatedTime, CreatedBy FROM Category where IsDeleted=0")
	if err != nil {
		glog.Errorf("Category->GetEntity->query err: %v\n", err)
	}
	defer rows.Close()
	var cates []Category
	for rows.Next() {
		var cate Category
		rows.Scan(&cate.ID, &cate.Name, &cate.CreatedTime, &cate.CreatedBy)
		cates = append(cates, cate)
	}
	return cates
}

//AddEntity insert data into db
func (cate Category) AddEntity() int64 {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into Category(Name,CreatedTime,CreatedBy) values(?,?,?)")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(cate.Name, cate.CreatedTime, cate.CreatedBy)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return id
}

//CategoryHandler handler
func CategoryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		var cate Category
		cates := cate.GetEntity()

		data := struct {
			Title     string
			Categorys []Category
		}{
			Title:     "Category",
			Categorys: cates,
		}
		categorytemplate.Execute(w, data)
	} else if r.Method == "POST" {
		cateName := r.FormValue("cateName")
		var cate Category
		cate.Name = cateName
		cate.CreatedTime = time.Now().Format(LongFormat)
		cate.CreatedBy = "johnson"

		lastInsertId := cate.AddEntity()
		if lastInsertId > -1 {
			//insert successful
		}
		http.Redirect(w, r, "/category", http.StatusMovedPermanently)
	}
}
