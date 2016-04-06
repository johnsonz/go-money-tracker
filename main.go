package main

import (
	"database/sql"
	"flag"
	"html/template"
	"net/http"
	"strconv"
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
type Subcategory struct {
	ID          int
	Name        string
	CreatedTime string
	CreatedBy   string
	Category
}

var categorytemplate *template.Template
var subcategorytemplate *template.Template

const (
	dbDrive     = "sqlite3"
	dbName      = "data.db"
	ShortFormat = "2006-01-02"
	LongFormat  = "2006-01-02 15:04:05"
)

func init() {
	flag.Parse()
	categorytemplate = template.Must(template.New("category.gtpl").ParseFiles("./templates/category.gtpl"))
	subcategorytemplate = template.Must(template.New("subcategory.gtpl").ParseFiles("./templates/subcategory.gtpl"))
	glog.Infoln("initial done")
}
func main() {
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	http.HandleFunc("/subcategory", SubcategoryHandler)
	err := http.ListenAndServe(":8888", nil) //设置监听的端口
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
			Title      string
			Categories []Category
		}{
			Title:      "Category",
			Categories: cates,
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
func (subcate Subcategory) GetEntity() []Subcategory {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID, Name, CreatedTime, CreatedBy FROM Subcategory where IsDeleted=0 and CategoryID=?")
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->stmt err: %v\n", err)
	}
	rows, err := stmt.Query(subcate.Category.ID)
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->rows err: %v\n", err)
	}
	var subcates []Subcategory
	for rows.Next() {
		var subcate Subcategory
		rows.Scan(&subcate.ID, &subcate.Name, &subcate.CreatedTime, &subcate.CreatedBy)
		subcates = append(subcates, subcate)
	}
	return subcates
}
func (subcate Subcategory) AddEntity() int64 {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->open db err: %v\n", err)
	}
	stmt, err := db.Prepare("insert into Subcategory(CategoryID,Name,CreatedTime,CreatedBy) values(?,?,?,?)")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(subcate.Category.ID, subcate.Name, subcate.CreatedTime, subcate.CreatedBy)
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func SubcategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var cate Category
		cates := cate.GetEntity()
		var subcate Subcategory
		cateIDFromURL := r.URL.Query().Get("id")
		cateID, err := strconv.Atoi(cateIDFromURL)
		subcate.Category.ID = 0
		if err != nil {
			if len(cates) > 0 {
				subcate.Category.ID = cates[0].ID
			}
			glog.Infof("Subcategory->convert id err: %v", err)
		} else {
			subcate.Category.ID = cateID
		}

		subcates := subcate.GetEntity()
		data := struct {
			Categories    []Category
			Subcategories []Subcategory
		}{
			Categories:    cates,
			Subcategories: subcates,
		}
		subcategorytemplate.Execute(w, data)
	} else if r.Method == "POST" {
		subcateName := r.FormValue("subcateName")
		cateIDForm := r.FormValue("category")
		var subcate Subcategory
		subcate.Name = subcateName
		cateID, err := strconv.Atoi(cateIDForm)
		if err != nil {
			glog.Errorf("SubcategoryHandler->convert id err: %v\n", err)
		}
		subcate.Category.ID = cateID
		subcate.CreatedTime = time.Now().Format(LongFormat)
		subcate.CreatedBy = "johnson"

		lastInsertId := subcate.AddEntity()
		if lastInsertId > -1 {
			//insert successful
		}
		http.Redirect(w, r, "/subcategory?id="+cateIDForm, http.StatusMovedPermanently)
	}
}
