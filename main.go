package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Category struct
type Category struct {
	ID          int
	Name        string
	CreatedTime string
	CreatedBy   string
}

var categorytemplate *template.Template

const (
	dbDrive = "sqlite3"
	dbName  = "data.db"
)

func init() {
	categorytemplate = template.Must(template.New("category").ParseFiles("./templates/category.gtpl"))
}
func main() {
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	err := http.ListenAndServe(":8888", nil)      //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe err: ", err)
	}
}

//GetEntity get data from db
func (cate Category) GetEntity() []Category {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		log.Fatal("open sqlites err: ", err)
	}
	rows, err := db.Query("SELECT ID, Name, CreatedTime, CreatedBy FROM Category where IsDeleted=0")
	if err != nil {
		log.Fatal("open sqlites err: ", err)
	}
	var cates []Category
	for rows.Next() {
		var cate Category
		rows.Scan(&cate.ID, &cate.Name, &Cate.CreatedTime, &cate.CreatedBy)
		cates = append(cates, cate)
	}
	return cates
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
	} else {

	}
}
