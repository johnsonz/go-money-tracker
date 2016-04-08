package main

import (
	"database/sql"
	"encoding/base64"
	"flag"
	"html/template"
	"io/ioutil"
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
	Selected    bool
}
type Subcategory struct {
	ID          int
	Name        string
	CreatedTime string
	CreatedBy   string
	Selected    bool
	Category
}
type Item struct {
	ID            int
	Store         string
	Address       string
	PurchasedDate string
	Receipt       string
	Remark        string
	CreatedTime   string
	CreatedBy     string
	Subcategory
}

var categorytemplate *template.Template
var subcategorytemplate *template.Template
var itemtemplate *template.Template

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
	itemtemplate = template.Must(template.New("item.gtpl").ParseFiles("./templates/item.gtpl"))
	glog.Infoln("initial done")
}
func main() {
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	http.HandleFunc("/subcategory", SubcategoryHandler)
	http.HandleFunc("/item", ItemHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
		for i, _ := range cates {
			if cates[i].ID == subcate.Category.ID {
				cates[i].Selected = true
			}
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
func (item Item) GetEntity() []Item {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select ID,Store,Address,PurchasedDate,ReceiptImage,Remark,CreatedTime,CreatedBy from Item where IsDeleted=0")
	if err != nil {
		glog.Errorf("Item->GetEntity->stmt err: %v\n", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		glog.Errorf("Item->GetEntity->query err: %v\n", err)
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var item Item
		var receiptimage []byte
		rows.Scan(&item.ID, &item.Store, &item.Address, &item.PurchasedDate, &receiptimage, &item.Remark, &item.CreatedTime, &item.CreatedBy)
		item.Receipt = base64.StdEncoding.EncodeToString(receiptimage)
		items = append(items, item)
	}
	return items
}
func (item Item) AddEntity() int64 {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->AddEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into Item(Store,Address,PurchasedDate,ReceiptImage,Remark,CreatedTime,CreatedBy) values(?,?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("Item->AddEntity->stmt err: %v\n", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(item.Store, item.Address, item.PurchasedDate, item.Receipt, item.Remark, item.CreatedTime, item.CreatedBy)
	if err != nil {
		glog.Errorf("Item->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Item->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var item Item
		var cate Category
		var subcate Subcategory

		items := item.GetEntity()
		cates := cate.GetEntity()
		subcates := subcate.GetEntity()
		data := struct {
			Items         []Item
			Categories    []Category
			Subcategories []Subcategory
		}{
			Items:         items,
			Categories:    cates,
			Subcategories: subcates,
		}
		itemtemplate.Execute(w, data)
	} else if r.Method == "POST" {
		purchasedDate := r.FormValue("purchaseddate")
		store := r.FormValue("store")
		address := r.FormValue("address")
		remark := r.FormValue("remark")
		file, _, err := r.FormFile("receiptimage")
		var receiptData []byte
		switch err {
		case nil:
			receiptData, err = ioutil.ReadAll(file)
			if err != nil {
				glog.Errorf("read file err: %v\n", err)
			}
			// receiptData, err = base64.StdEncoding.DecodeString(string(receiptData))
			// if err != nil {
			// 	glog.Errorf("convert file to base64 err: %v\n", err)
			// }
		case http.ErrMissingFile:
			glog.Infof("no file uploaded \n")
		default:
			glog.Errorf("upload file err: %v\n", err)
		}
		_, err = time.Parse(ShortFormat, purchasedDate)
		if err != nil {
			glog.Errorf("parse purchased date %s err: %v\n", purchasedDate, err)
		}
		var item Item
		item.PurchasedDate = purchasedDate
		item.Receipt = string(receiptData)
		item.Store = store
		item.Address = address
		item.Remark = remark
		item.CreatedTime = time.Now().Format(LongFormat)
		item.CreatedBy = "johnson"

		lastInsertId := item.AddEntity()
		if lastInsertId > -1 {
			//insert successful
		}
		http.Redirect(w, r, "/item", http.StatusMovedPermanently)
	}
}
