package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"go-money-tracker/mtacrypto"
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
	CreatedBy   int
	Selected    bool
}
type CategoryEncrypted struct {
	ID          int
	Name        []byte
	CreatedTime []byte
	CreatedBy   int
	Selected    bool
}
type Subcategory struct {
	ID          int
	Name        string
	CreatedTime string
	CreatedBy   int
	Selected    bool
	Category
}
type SubcategoryEncrypted struct {
	ID          int
	Name        []byte
	CreatedTime []byte
	CreatedBy   int
	Selected    bool
	CategoryEncrypted
}
type Item struct {
	ID            int
	Store         string
	Address       string
	PurchasedDate string
	Receipt       string
	Remark        string
	CreatedTime   string
	CreatedBy     int
	Subcategory
}
type ItemEncrypted struct {
	ID            int
	Store         []byte
	Address       []byte
	PurchasedDate []byte
	Receipt       []byte
	Remark        []byte
	CreatedTime   []byte
	CreatedBy     int
	SubcategoryEncrypted
}
type Detail struct {
	ID          int
	Name        string
	Price       float64
	Quantity    int
	LabelOne    string
	LabelTwo    string
	CreatedTime string
	CreatedBy   string
	Item
}

var categorytemplate *template.Template
var subcategorytemplate *template.Template
var itemtemplate *template.Template
var detailtemplate *template.Template

const (
	dbDrive     = "sqlite3"
	dbName      = "data.db"
	ShortFormat = "2006-01-02"
	LongFormat  = "2006-01-02 15:04:05"
	key         = "abcdefghijklmnopqrstuvwxyz012345"
)

func init() {
	flag.Parse()
	categorytemplate = template.Must(template.New("category.gtpl").
		ParseFiles("./templates/category.gtpl"))
	subcategorytemplate = template.Must(template.New("subcategory.gtpl").
		ParseFiles("./templates/subcategory.gtpl"))
	itemtemplate = template.Must(template.New("item.gtpl").
		ParseFiles("./templates/item.gtpl"))
	detailtemplate = template.Must(template.New("detail.gtpl").
		Funcs(template.FuncMap{"getamount": GetAmount}).
		ParseFiles("./templates/detail.gtpl"))
	glog.Infoln("initial done")
}
func main() {
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	http.HandleFunc("/subcategory", SubcategoryHandler)
	http.HandleFunc("/getsubcategory", GetSubcategoryHandler)
	http.HandleFunc("/item", ItemHandler)
	http.HandleFunc("/detail", DetailHandler)
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
		var ecate CategoryEncrypted
		rows.Scan(&ecate.ID, &ecate.Name, &ecate.CreatedTime, &ecate.CreatedBy)
		cates = append(cates, ecate.Decrypt())
	}
	return cates
}

//AddEntity insert data into db
func (cate Category) AddEntity() int64 {
	ecate := cate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into Category(Name,CreatedTime,CreatedBy) values(?,?,?)")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(ecate.Name, ecate.CreatedTime, ecate.CreatedBy)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return id
}
func (cate Category) Encrypt() CategoryEncrypted {
	name, err := mtcrypto.AESEncrypt(key, cate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", cate.Name, err)
	}
	ctime, err := mtcrypto.AESEncrypt(key, cate.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", cate.CreatedTime, err)
	}
	return CategoryEncrypted{
		ID:          cate.ID,
		Name:        name,
		CreatedTime: ctime,
		CreatedBy:   cate.CreatedBy,
	}
}
func (ecate CategoryEncrypted) Decrypt() Category {
	name, err := mtcrypto.AESDecrypt(key, ecate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", ecate.Name, err)
	}
	ctime, err := mtcrypto.AESDecrypt(key, ecate.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", ecate.CreatedTime, err)
	}
	return Category{
		ID:          ecate.ID,
		Name:        string(name),
		CreatedTime: string(ctime),
		CreatedBy:   ecate.CreatedBy,
	}
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
		cate.CreatedBy = 0

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
		var esubcate SubcategoryEncrypted
		rows.Scan(&esubcate.ID, &esubcate.Name, &esubcate.CreatedTime, &esubcate.CreatedBy)
		subcates = append(subcates, esubcate.Decrypt())
	}
	return subcates
}
func (subcate Subcategory) AddEntity() int64 {
	esubcate := subcate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->open db err: %v\n", err)
	}
	stmt, err := db.Prepare("insert into Subcategory(CategoryID,Name,CreatedTime,CreatedBy) values(?,?,?,?)")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(esubcate.CategoryEncrypted.ID, esubcate.Name, esubcate.CreatedTime, esubcate.CreatedBy)
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func (subcate Subcategory) Encrypt() SubcategoryEncrypted {
	name, err := mtcrypto.AESEncrypt(key, subcate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", subcate.Name, err)
	}
	ctime, err := mtcrypto.AESEncrypt(key, subcate.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", subcate.CreatedTime, err)
	}
	return SubcategoryEncrypted{
		ID:                subcate.ID,
		Name:              name,
		CreatedTime:       ctime,
		CreatedBy:         subcate.CreatedBy,
		CategoryEncrypted: subcate.Category.Encrypt(),
	}
}
func (esubcate SubcategoryEncrypted) Decrypt() Subcategory {
	name, err := mtcrypto.AESDecrypt(key, esubcate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", esubcate.Name, err)
	}
	ctime, err := mtcrypto.AESDecrypt(key, esubcate.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", esubcate.CreatedTime, err)
	}
	return Subcategory{
		ID:          esubcate.ID,
		Name:        string(name),
		CreatedTime: string(ctime),
		CreatedBy:   esubcate.CreatedBy,
		Category:    esubcate.CategoryEncrypted.Decrypt(),
	}
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
		subcate.CreatedBy = 0

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
	stmt, err := db.Prepare("select ID,Store,Address,PurchasedDate,ReceiptImage,Remark,CreatedTime,CreatedBy,SubcategoryID,SubcategoryName,CategoryID,CategoryName from vw_Item where IsDeleted=0")
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
		rows.Scan(&item.ID, &item.Store, &item.Address, &item.PurchasedDate, &receiptimage, &item.Remark, &item.CreatedTime, &item.CreatedBy,
			&item.Subcategory.ID, &item.Subcategory.Name, &item.Subcategory.Category.ID, &item.Subcategory.Category.Name)
		item.Receipt = base64.StdEncoding.EncodeToString(receiptimage)
		items = append(items, item)
	}
	return items
}
func (item Item) AddEntity() int64 {
	eitem := item.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->AddEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into Item(Store,Address,PurchasedDate,ReceiptImage,Remark,CreatedTime,CreatedBy,SubcategoryID) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("Item->AddEntity->stmt err: %v\n", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(eitem.Store, eitem.Address, eitem.PurchasedDate, eitem.Receipt, eitem.Remark, eitem.CreatedTime, eitem.CreatedBy, eitem.SubcategoryEncrypted.ID)
	if err != nil {
		glog.Errorf("Item->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Item->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func (item Item) Encrypt() ItemEncrypted {
	store, err := mtcrypto.AESEncrypt(key, item.Store)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.Store, err)
	}
	address, err := mtcrypto.AESEncrypt(key, item.Address)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.Address, err)
	}
	purdate, err := mtcrypto.AESEncrypt(key, item.PurchasedDate)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.PurchasedDate, err)
	}
	receipt, err := mtcrypto.AESEncrypt(key, item.Receipt)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.Receipt, err)
	}
	remark, err := mtcrypto.AESEncrypt(key, item.Remark)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.Remark, err)
	}
	ctime, err := mtcrypto.AESEncrypt(key, item.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.CreatedTime, err)
	}

	return ItemEncrypted{
		ID:                   item.ID,
		Store:                store,
		Address:              address,
		PurchasedDate:        purdate,
		Receipt:              receipt,
		Remark:               remark,
		CreatedTime:          ctime,
		CreatedBy:            item.CreatedBy,
		SubcategoryEncrypted: item.Subcategory.Encrypt(),
	}
}
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var item Item
		var cate Category
		var subcate Subcategory

		items := item.GetEntity()
		cates := cate.GetEntity()

		// itemID := r.URL.Query().Get("id")
		cateID := r.URL.Query().Get("cid")
		subcateID := r.URL.Query().Get("sid")

		cid, err := strconv.Atoi(cateID)
		if err != nil {
			cid = 0
			if len(cates) > 0 {
				cid = cates[0].ID
			}
			glog.Infof("convert cid to int err: %v\n", err)
		}
		subcate.Category.ID = cid
		subcates := subcate.GetEntity()
		sid, err := strconv.Atoi(subcateID)
		if err != nil {
			sid = 0
			if len(subcates) > 0 {
				sid = subcates[0].ID
			}
			glog.Infof("convert cid to int err: %v\n", err)
		}
		for i, _ := range cates {
			if cates[i].ID == cid {
				cates[i].Selected = true
				break
			}
		}
		for i, _ := range subcates {
			if subcates[i].ID == sid {
				subcates[i].Selected = true
				break
			}
		}
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
		cateID := r.FormValue("category")
		subcateID := r.FormValue("subcategory")
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
		sid, err := strconv.Atoi(subcateID)
		if err != nil {
			glog.Errorf("convert sid to int err； %v \n", err)
		}
		_, err = time.Parse(ShortFormat, purchasedDate)
		if err != nil {
			glog.Errorf("parse purchased date %s err: %v\n", purchasedDate, err)
		}
		var item Item
		item.Subcategory.ID = sid
		item.PurchasedDate = purchasedDate
		item.Receipt = string(receiptData)
		item.Store = store
		item.Address = address
		item.Remark = remark
		item.CreatedTime = time.Now().Format(LongFormat)
		item.CreatedBy = 0

		lastInsertId := item.AddEntity()
		if lastInsertId > -1 {
			//insert successful
		}
		http.Redirect(w, r, "/item?id="+strconv.Itoa(int(lastInsertId))+"&sid="+strconv.Itoa(sid)+"&cid="+cateID, http.StatusMovedPermanently)
	}
}
func (detail Detail) GetEntity() []Detail {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("select ID,Price,Quantity,LabelOne,LabelTwo,CreatedTime,CreatedBy from Detail where IsDeleted=0 and ItemID=?")
	if err != nil {
		glog.Errorf("db prepare err: %v\n", err)
	}
	rows, err := stmt.Query(detail.Item.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	var details []Detail
	for rows.Next() {
		var detail Detail
		var labelone, labeltwo []byte
		rows.Scan(&detail.ID, &detail.Price, &detail.Quantity, &labelone,
			&labeltwo, &detail.CreatedTime, &detail.CreatedBy)
		detail.LabelOne = base64.StdEncoding.EncodeToString(labelone)
		detail.LabelTwo = base64.StdEncoding.EncodeToString(labeltwo)
		details = append(details, detail)
	}
	return details
}

func (detail Detail) AddEntity() int64 {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("insert into Detail(ItemID,Name,Price,Quantity,LabelOne,LabelTwo,Remark,CreatedTime,CreatedBy) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(detail.Item.ID, detail.Name, detail.Price,
		detail.Quantity, detail.LabelOne, detail.LabelTwo, detail.Remark,
		detail.CreatedTime, detail.CreatedBy)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("get LastInsertId err: %v\n", err)
	}
	return id
}
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		itemID := r.URL.Query().Get("id")
		var detail Detail
		iid, err := strconv.Atoi(itemID)
		if err != nil {
			detail.Item.ID = 0
			glog.Errorf("get detail by item id err: %v", err)
		} else {
			detail.Item.ID = iid
		}
		details := detail.GetEntity()
		data := struct {
			ItemID  int
			Details []Detail
		}{
			ItemID:  iid,
			Details: details,
		}
		detailtemplate.Execute(w, data)
	} else if r.Method == "POST" {
		var detail Detail
		itemid := r.FormValue("itemid")
		iid, err := strconv.Atoi(itemid)
		if err != nil {
			glog.Fatalf("get item id %s err: %v", itemid, err)
		}
		detail.Item.ID = iid
		name := r.FormValue("name")
		detail.Name = name
		price := r.FormValue("price")
		pri, err := strconv.ParseFloat(price, 64)
		if err != nil {
			detail.Price = 0.0
			glog.Errorf("parse float %s err: %v", price, err)
		} else {
			detail.Price = pri
		}
		quantity := r.FormValue("quantity")
		quan, err := strconv.Atoi(quantity)
		if err != nil {
			detail.Quantity = 1
			glog.Errorf("parse float %s err: %v", price, err)
		} else {
			detail.Quantity = quan
		}
		var labeloneData, labeltwoData []byte
		labelone, _, err := r.FormFile("labelone")
		switch err {
		case nil:
			labeloneData, err = ioutil.ReadAll(labelone)
			if err != nil {
				glog.Errorf("read file err: %v\n", err)
			}
		case http.ErrMissingFile:
			glog.Infof("no file uploaded \n")
		default:
			glog.Errorf("upload file err: %v\n", err)
		}
		labeltwo, _, err := r.FormFile("labeltwo")
		switch err {
		case nil:
			labeltwoData, err = ioutil.ReadAll(labeltwo)
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
		detail.LabelOne = string(labeloneData)
		detail.LabelTwo = string(labeltwoData)
		remark := r.FormValue("remark")
		detail.Remark = remark
		detail.CreatedTime = time.Now().Format(LongFormat)
		detail.CreatedBy = "johnson"
		lastInsertId := detail.AddEntity()
		if lastInsertId > 0 {
			//insert successful
		}
		http.Redirect(w, r, "/detail?id="+itemid, http.StatusMovedPermanently)
	}
}
func GetSubcategoryHandler(w http.ResponseWriter, r *http.Request) {

	var subcate Subcategory
	cateIDForm := r.URL.Query().Get("id")
	subcate.Category.ID = 0
	cateID, err := strconv.Atoi(cateIDForm)
	if err != nil {
		glog.Infof("convert cateID %s to int err: %v", cateIDForm, err)
	} else {
		subcate.Category.ID = cateID
	}
	subcates := subcate.GetEntity()
	data, err := json.Marshal(subcates)
	if err != nil {
		glog.Errorf("convert %T to json err: %v", subcates, err)
	}
	fmt.Fprint(w, string(data))
}
func GetAmount(price float64, quantity int) float64 {
	return price * float64(quantity)
}
