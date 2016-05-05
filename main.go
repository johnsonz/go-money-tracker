package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/johnsonz/go-money-tracker/mtacrypto"
	"github.com/johnsonz/go-money-tracker/mtconverter"
	_ "github.com/mattn/go-sqlite3"
)

// Category
type Category struct {
	ID       int
	Name     string
	Selected bool
	Operation
}
type CategoryEncrypted struct {
	ID       int
	Name     []byte
	Selected bool
	OperationEncrypted
}
type Subcategory struct {
	ID       int
	Name     string
	Selected bool
	Category
	Operation
}
type SubcategoryEncrypted struct {
	ID       int
	Name     []byte
	Selected bool
	CategoryEncrypted
	OperationEncrypted
}
type Item struct {
	ID            int
	Store         string
	Address       string
	PurchasedDate string
	Receipt       string
	Amount        float64
	Remark        string
	Subcategory
	Operation
}
type ItemEncrypted struct {
	ID            int
	Store         []byte
	Address       []byte
	PurchasedDate []byte
	Receipt       []byte
	Amount        []byte
	Remark        []byte
	SubcategoryEncrypted
	OperationEncrypted
}
type Detail struct {
	ID       int
	Name     string
	Price    float64
	Quantity int64
	LabelOne string
	LabelTwo string
	Remark   string
	Item
	Operation
}
type DetailEncrypted struct {
	ID       int
	Name     []byte
	Price    []byte
	Quantity []byte
	LabelOne []byte
	LabelTwo []byte
	Remark   []byte
	ItemEncrypted
	OperationEncrypted
}
type User struct {
	ID            int
	Username      string
	Password      string
	Nick          string
	LastLoginTime string
	LastLoginIP   string
	Hostname      string
	Operation
}
type UserEncrypted struct {
	ID            int
	Username      []byte
	Password      []byte
	Nick          string
	LastLoginTime []byte
	LastLoginIP   []byte
	Hostname      []byte
	OperationEncrypted
}
type Operation struct {
	CreatedTime string
	CreatedBy   int
	UpdatedTime string
	UpdatedBy   int
	DeletedTime string
	DeletedBy   int
}
type OperationEncrypted struct {
	CreatedTime []byte
	CreatedBy   int
	UpdatedTime []byte
	UpdatedBy   int
	DeletedTime []byte
	DeletedBy   int
}
type Pagination struct {
	Count    int
	Index    int
	Size     int
	Previous int
	Next     int
}

var store *sessions.CookieStore
var templates *template.Template

const (
	dbDrive     = "sqlite3"
	dbName      = "data.db"
	ShortFormat = "2006-01-02"
	LongFormat  = "2006-01-02 15:04:05"
	key         = "abcdefghijklmnopqrstuvwxyz012345"
	sessionsKey = "johnson"
	sessionName = "mt"
	pageSize    = 3
	pageNavSize = 5
	delAction   = "del"
	updAction   = "upd"
	addAction   = "add"
)

func init() {
	flag.Parse()
	templates = template.Must(template.New("templates").
		Funcs(template.FuncMap{"getamount": GetAmount, "plus": func(m, n int) int { return m + n }, "minus": func(m, n int) int { return m - n }}).
		ParseGlob("./templates/*.gtpl"))
	store = sessions.NewCookieStore([]byte(sessionsKey))
	glog.Infoln("initial done")
}
func main() {
	//http.HandleFunc("/", LoginHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/category", CategoryHandler) //设置访问的路由
	http.HandleFunc("/category/del", CategoryDelHandler)
	http.HandleFunc("/subcategory", SubcategoryHandler)
	http.HandleFunc("/subcategory/del", SubcategoryDelHandler)
	http.HandleFunc("/getsubcategory", GetSubcategoryHandler)
	http.HandleFunc("/item", ItemHandler)
	http.HandleFunc("/item/del", ItemDelHandler)
	http.HandleFunc("/detail", DetailHandler)
	http.HandleFunc("/detail/del", DetailDelHandler)
	http.HandleFunc("/user", UserHandler)
	http.HandleFunc("/user/del", UserDelHandler)
	http.HandleFunc("/rmrept", RemoveReceiptHandler)
	http.HandleFunc("/rmlabel", RemoveLabelHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8888", context.ClearHandler(http.DefaultServeMux)) //设置监听的端口
	if err != nil {
		glog.Errorf("main->ListenAndServe err: %v\n", err)
	}
}
func CheckSessions(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		glog.Errorf("get session err: %v", err)
	}
	if flashes := session.Flashes(); len(flashes) > 0 {
		// Use the flash values.
		glog.Infof("get session flashes: %v", flashes)
	} else {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}
func (user User) GetEntity(pagination Pagination) []User {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("select ID, Username, Password,Nick, Hostname,LastLoginTime,LastLoginIP,CreatedTime,CreatedBy from User where IsDeleted=0 limit ? offset ?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	rows, err := stmt.Query(pagination.Size, pagination.Size*(pagination.Index-1))
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	var users []User
	for rows.Next() {
		var euser UserEncrypted
		rows.Scan(&euser.ID, &euser.Username, &euser.Password, &euser.Nick, &euser.Hostname,
			&euser.LastLoginTime, &euser.LastLoginIP, &euser.OperationEncrypted.CreatedTime,
			&euser.OperationEncrypted.CreatedBy)
		users = append(users, euser.Decrypt())
	}
	return users
}
func (user User) GetAllEntity() []User {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("select ID, Username, Password,Nick, Hostname,LastLoginTime,LastLoginIP,CreatedTime,CreatedBy from User where IsDeleted=0")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	rows, err := stmt.Query()
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	var users []User
	for rows.Next() {
		var euser UserEncrypted
		rows.Scan(&euser.ID, &euser.Username, &euser.Password, &euser.Nick, &euser.Hostname,
			&euser.LastLoginTime, &euser.LastLoginIP, &euser.OperationEncrypted.CreatedTime,
			&euser.OperationEncrypted.CreatedBy)
		users = append(users, euser.Decrypt())
	}
	return users
}
func (user User) AddEntity() int64 {
	user.Username = fmt.Sprintf("%X", mtcrypto.MD5(user.Username))
	user.Password = fmt.Sprintf("%X", mtcrypto.MD5(user.Password))
	euser := user.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("insert into User(Username, Password,Nick, Hostname, CreatedTime, CreatedBy) values(?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(euser.Username, euser.Password, euser.Nick, euser.Hostname,
		euser.OperationEncrypted.CreatedTime, euser.OperationEncrypted.CreatedBy)
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}

	return lastInsertId
}
func (user User) UpdEntity() int64 {
	user.Username = fmt.Sprintf("%X", mtcrypto.MD5(user.Username))
	user.Password = fmt.Sprintf("%X", mtcrypto.MD5(user.Password))
	euser := user.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("update User set Username=?, Password=?,Nick=?, Hostname=?,UpdatedTime=?, UpdatedBy=? where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(euser.Username, euser.Password, euser.Nick, euser.Hostname,
		euser.OperationEncrypted.UpdatedTime, euser.OperationEncrypted.UpdatedBy, euser.ID)
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}

	return rowsAffected
}
func (user User) UpdLoginInfo() int64 {
	euser := user.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("update User set LastLoginTime=?,LastLoginIP=?,UpdatedTime=?, UpdatedBy=? where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(euser.LastLoginTime, euser.LastLoginIP,
		euser.OperationEncrypted.UpdatedTime, euser.OperationEncrypted.UpdatedBy, euser.ID)
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}

	return rowsAffected
}
func (user User) DelEntity() int64 {
	euser := user.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("update User set IsDeleted=1,DeletedTime=?,DeletedBy=? where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(euser.OperationEncrypted.DeletedTime, euser.OperationEncrypted.DeletedBy,
		euser.ID)
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("query err: %v\n", err)
	}

	return rowsAffected
}
func (user User) Encrypt() UserEncrypted {
	username, err := mtcrypto.AESEncrypt(key, user.Username)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", user.Username, err)
	}
	password, err := mtcrypto.AESEncrypt(key, user.Password)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", user.Password, err)
	}
	hostname, err := mtcrypto.AESEncrypt(key, user.Hostname)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", user.Hostname, err)
	}
	lastlogintime, err := mtcrypto.AESEncrypt(key, user.LastLoginTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", user.LastLoginTime, err)
	}
	lastloginip, err := mtcrypto.AESEncrypt(key, user.LastLoginIP)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", user.LastLoginIP, err)
	}
	return UserEncrypted{
		ID:                 user.ID,
		Username:           username,
		Password:           password,
		Nick:               user.Nick,
		Hostname:           hostname,
		LastLoginTime:      lastlogintime,
		LastLoginIP:        lastloginip,
		OperationEncrypted: user.Operation.Encryt(),
	}
}
func (euser UserEncrypted) Decrypt() User {
	username, err := mtcrypto.AESDecrypt(key, euser.Username)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", euser.Username, err)
	}
	password, err := mtcrypto.AESDecrypt(key, euser.Password)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", euser.Password, err)
	}
	hostname, err := mtcrypto.AESDecrypt(key, euser.Hostname)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", euser.Hostname, err)
	}
	lastlogintime, err := mtcrypto.AESDecrypt(key, euser.LastLoginTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", euser.LastLoginTime, err)
	}
	lastloginip, err := mtcrypto.AESDecrypt(key, euser.LastLoginIP)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", euser.LastLoginIP, err)
	}
	return User{
		ID:            euser.ID,
		Username:      string(username),
		Password:      string(password),
		Nick:          euser.Nick,
		Hostname:      string(hostname),
		LastLoginTime: string(lastlogintime),
		LastLoginIP:   string(lastloginip),
		Operation:     euser.OperationEncrypted.Decryt(),
	}
}
func (User User) Count() (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from User where IsDeleted=0", nil).Scan(&count)
	return count
}
func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var user User
		if r.URL.Query().Get("action") == delAction {
			userid := r.FormValue("id")
			id, err := strconv.Atoi(userid)
			if err != nil {
				glog.Errorf("convert userid %s err: %v", userid, err)
			}
			user.ID = id
			user.Operation.DeletedTime = time.Now().Format(LongFormat)
			user.Operation.DeletedBy = 0
			rowsAffected := user.DelEntity()
			if rowsAffected > 0 {
				//successful
			}
		}
		page := r.URL.Query().Get("page")
		count := user.Count()
		pagination := GetPagination(page, count)
		users := user.GetEntity(pagination)
		data := struct {
			Title      string
			Users      []User
			Pagination Pagination
		}{
			Title:      "User",
			Users:      users,
			Pagination: pagination,
		}
		templates.ExecuteTemplate(w, "user.gtpl", data)
	} else if r.Method == "POST" {
		pageIndex := r.FormValue("pageIndex")
		if r.FormValue("update") == "Update" {
			uid := r.FormValue("updatedid")
			id, err := strconv.Atoi(uid)
			if err != nil {
				glog.Errorf("convert uid %s err: %v", uid, err)
			}
			un := r.FormValue("updatedname")
			pwd := r.FormValue("updatedpassword")
			nick := r.FormValue("updatednick")
			hn := r.FormValue("updatedhost")
			var user User
			user.ID = id
			user.Username = un
			user.Password = pwd
			user.Nick = nick
			user.Hostname = hn
			user.Operation.UpdatedTime = time.Now().Format(LongFormat)
			user.Operation.UpdatedBy = 0
			rowsAffected := user.UpdEntity()
			if rowsAffected > 0 {
				//successful
			}
		} else if r.FormValue("create") == "Create" {
			un := r.FormValue("createdusername")
			pwd := r.FormValue("createdpassword")
			nick := r.FormValue("creatednick")
			hn := r.FormValue("createdhostname")
			var user User
			user.Username = un
			user.Password = pwd
			user.Nick = nick
			user.Hostname = hn
			user.Operation.CreatedTime = time.Now().Format(LongFormat)
			user.Operation.CreatedBy = 0
			lastInsertId := user.AddEntity()
			if lastInsertId > -1 {
				//successful
			}
		}
		http.Redirect(w, r, "/user?page="+pageIndex, http.StatusMovedPermanently)
	}
}
func UserDelHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		glog.Errorf("convert userid %s err: %v", userid, err)
		fmt.Fprint(w, false)
		return
	}
	var user User
	user.ID = id
	user.Operation.DeletedTime = time.Now().Format(LongFormat)
	user.Operation.DeletedBy = 0
	rowsAffected := user.DelEntity()
	if rowsAffected > 0 {
		//successful
		fmt.Fprint(w, true)
		return
	}
	fmt.Fprint(w, false)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := struct {
			Username string
			HasError string
		}{
			Username: "",
			HasError: "",
		}
		// logintemplate.Execute(w, data)
		templates.ExecuteTemplate(w, "login.gtpl", data)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		// hostname, err := os.Hostname()
		// if err != nil {
		// 	glog.Errorf("get hostname err: %v", err)
		// }
		var user User
		users := user.GetAllEntity()
		usernamemd5 := fmt.Sprintf("%X", mtcrypto.MD5(username))
		passwordmd5 := fmt.Sprintf("%X", mtcrypto.MD5(password))
		// hostnamemd5 := fmt.Sprintf("%X", mtcrypto.MD5(hostname))

		for _, u := range users {
			// if u.Username == usernamemd5 &&
			// 	u.Password == passwordmd5 &&
			// 	u.Hostname == hostnamemd5 {
			if u.Username == usernamemd5 &&
				u.Password == passwordmd5 {
				session, err := store.Get(r, sessionName)
				if err != nil {
					glog.Errorf("get session err: %v", err)
				}
				session.Options = &sessions.Options{
					Path:     "/",
					MaxAge:   86400,
					HttpOnly: true,
				}
				if flashes := session.Flashes(); len(flashes) > 0 {
					// Use the flash values.
					glog.Infof("get session flashes: %v", flashes)
				} else {
					// Set a new flash.
					session.AddFlash(u.ID)
				}
				user.ID = u.ID
				user.LastLoginTime = time.Now().Format(LongFormat)
				user.LastLoginIP = GetLocalIP()
				user.Operation.UpdatedTime = time.Now().Format(LongFormat)
				user.Operation.UpdatedBy = u.ID
				user.UpdLoginInfo()
				session.Save(r, w)
				http.Redirect(w, r, "/category", http.StatusMovedPermanently)
			}
		}
		data := struct {
			Username string
			HasError string
		}{
			Username: username,
			HasError: "has-error",
		}
		// logintemplate.Execute(w, data)
		templates.ExecuteTemplate(w, "login.gtpl", data)
	}
}
func (op Operation) Encryt() OperationEncrypted {
	ctime, err := mtcrypto.AESEncrypt(key, op.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.CreatedTime, err)
	}
	utime, err := mtcrypto.AESEncrypt(key, op.UpdatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.UpdatedTime, err)
	}
	dtime, err := mtcrypto.AESEncrypt(key, op.DeletedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.DeletedTime, err)
	}
	return OperationEncrypted{
		CreatedTime: ctime,
		CreatedBy:   op.CreatedBy,
		UpdatedTime: utime,
		UpdatedBy:   op.UpdatedBy,
		DeletedTime: dtime,
		DeletedBy:   op.DeletedBy,
	}
}
func (op OperationEncrypted) Decryt() Operation {
	ctime, err := mtcrypto.AESDecrypt(key, op.CreatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.CreatedTime, err)
	}
	utime, err := mtcrypto.AESDecrypt(key, op.UpdatedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.UpdatedTime, err)
	}
	dtime, err := mtcrypto.AESDecrypt(key, op.DeletedTime)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", op.DeletedTime, err)
	}
	return Operation{
		CreatedTime: string(ctime),
		CreatedBy:   op.CreatedBy,
		UpdatedTime: string(utime),
		UpdatedBy:   op.UpdatedBy,
		DeletedTime: string(dtime),
		DeletedBy:   op.DeletedBy,
	}
}
func (cate Category) GetEntity(pagination Pagination) []Category {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->GetEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID, Name,CreatedTime, CreatedBy FROM Category where IsDeleted=0 limit ? offset ?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	rows, err := stmt.Query(pagination.Size, pagination.Size*(pagination.Index-1))
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
func (cate Category) GetAllEntity() []Category {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->GetEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID, Name,CreatedTime, CreatedBy FROM Category where IsDeleted=0")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	rows, err := stmt.Query()
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
	res, err := stmt.Exec(ecate.Name, ecate.OperationEncrypted.CreatedTime, ecate.OperationEncrypted.CreatedBy)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return id
}
func (cate Category) UpdEntity() int64 {
	ecate := cate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Category set Name=?,UpdatedTime=?,UpdatedBy=? where id=?")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(ecate.Name, ecate.OperationEncrypted.UpdatedTime, ecate.OperationEncrypted.UpdatedBy, ecate.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
}
func (cate Category) DelEntity() int64 {
	ecate := cate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Category set IsDeleted=1,DeletedTime=?,DeletedBy=? where id=?")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(ecate.OperationEncrypted.DeletedTime, ecate.OperationEncrypted.DeletedBy, ecate.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
}
func (cate Category) Encrypt() CategoryEncrypted {
	name, err := mtcrypto.AESEncrypt(key, cate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", cate.Name, err)
	}
	return CategoryEncrypted{
		ID:                 cate.ID,
		Name:               name,
		OperationEncrypted: cate.Operation.Encryt(),
	}
}
func (ecate CategoryEncrypted) Decrypt() Category {
	name, err := mtcrypto.AESDecrypt(key, ecate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", ecate.Name, err)
	}
	return Category{
		ID:        ecate.ID,
		Name:      string(name),
		Operation: ecate.OperationEncrypted.Decryt(),
	}
}
func (cate Category) Count() (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from Category where IsDeleted=0", nil).Scan(&count)
	return count
}
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	if r.Method == "GET" {
		var cate Category
		page := r.URL.Query().Get("page")
		pagination := GetPagination(page, cate.Count())
		cates := cate.GetEntity(pagination)
		data := struct {
			Title      string
			Categories []Category
			Pagination Pagination
		}{
			Title:      "Category",
			Categories: cates,
			Pagination: pagination,
		}
		// categorytemplate.Execute(w, data)
		templates.ExecuteTemplate(w, "category.gtpl", data)
	} else if r.Method == "POST" {
		pageIndex := r.FormValue("pageIndex")
		if r.FormValue("update") == "Update" {
			updatedname := r.FormValue("updatedname")
			updatedid := r.FormValue("updatedid")
			id, err := strconv.Atoi(updatedid)
			if err != nil {
				glog.Errorf("convert string %s to int err: %v", updatedid, err)
			}
			var cate Category
			cate.ID = id
			cate.Name = updatedname
			cate.Operation.UpdatedTime = time.Now().Format(LongFormat)
			cate.Operation.UpdatedBy = 0
			rowsAffected := cate.UpdEntity()
			if rowsAffected > 0 {
				//successful
			}

		} else if r.FormValue("create") == "Create" {
			cateName := r.FormValue("createdname")
			var cate Category
			cate.Name = cateName
			cate.Operation.CreatedTime = time.Now().Format(LongFormat)
			cate.Operation.CreatedBy = 0

			lastInsertId := cate.AddEntity()
			if lastInsertId > -1 {
				//insert successful
			}
		}
		http.Redirect(w, r, "/category?page="+pageIndex, http.StatusMovedPermanently)
	}
}
func CategoryDelHandler(w http.ResponseWriter, r *http.Request) {
	var cate Category
	cid := r.FormValue("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		glog.Errorf("convert string %s to int err: %v", cid, err)
		fmt.Fprint(w, false)
		return
	}
	cate.ID = id
	cate.Operation.DeletedTime = time.Now().Format(LongFormat)
	cate.Operation.DeletedBy = 0
	rowsAffected := cate.DelEntity()
	if rowsAffected > 0 {
		//successful
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
}
func (subcate Subcategory) GetEntity(pagination Pagination) []Subcategory {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID,CategoryID, Name, CreatedTime, CreatedBy FROM Subcategory where IsDeleted=0 and CategoryID=? limit ? offset ?")
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->stmt err: %v\n", err)
	}
	rows, err := stmt.Query(subcate.Category.ID, pagination.Size, pagination.Size*(pagination.Index-1))
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->rows err: %v\n", err)
	}
	var subcates []Subcategory
	for rows.Next() {
		var esubcate SubcategoryEncrypted
		rows.Scan(&esubcate.ID, &esubcate.CategoryEncrypted.ID, &esubcate.Name,
			&esubcate.OperationEncrypted.CreatedTime, &esubcate.OperationEncrypted.CreatedBy)
		subcates = append(subcates, esubcate.Decrypt())
	}
	return subcates
}
func (subcate Subcategory) GetAllEntity() []Subcategory {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT ID,CategoryID, Name, CreatedTime, CreatedBy FROM Subcategory where IsDeleted=0 and CategoryID=?")
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
		rows.Scan(&esubcate.ID, &esubcate.CategoryEncrypted.ID, &esubcate.Name,
			&esubcate.OperationEncrypted.CreatedTime, &esubcate.OperationEncrypted.CreatedBy)
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
	res, err := stmt.Exec(esubcate.CategoryEncrypted.ID, esubcate.Name,
		esubcate.OperationEncrypted.CreatedTime, esubcate.OperationEncrypted.CreatedBy)
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func (subcate Subcategory) UpdEntity() int64 {
	esubcate := subcate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->open db err: %v\n", err)
	}
	stmt, err := db.Prepare("update Subcategory set Name=?,UpdatedTime=?,UpdatedBy=? where ID=?")
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(esubcate.Name, esubcate.OperationEncrypted.UpdatedTime,
		esubcate.OperationEncrypted.UpdatedBy, esubcate.ID)
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Subcategory->AddEntity->get LastInsertId err: %v\n", err)
	}
	return rowsAffected
}
func (subcate Subcategory) DelEntity() int64 {
	ecate := subcate.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Subcategory set IsDeleted=1,DeletedTime=?,DeletedBy=? where id=?")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(ecate.OperationEncrypted.DeletedTime, ecate.OperationEncrypted.DeletedBy, ecate.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
}
func (subcate Subcategory) Encrypt() SubcategoryEncrypted {
	name, err := mtcrypto.AESEncrypt(key, subcate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", subcate.Name, err)
	}
	return SubcategoryEncrypted{
		ID:                 subcate.ID,
		Name:               name,
		OperationEncrypted: subcate.Operation.Encryt(),
		CategoryEncrypted:  subcate.Category.Encrypt(),
	}
}
func (esubcate SubcategoryEncrypted) Decrypt() Subcategory {
	name, err := mtcrypto.AESDecrypt(key, esubcate.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", esubcate.Name, err)
	}
	return Subcategory{
		ID:        esubcate.ID,
		Name:      string(name),
		Operation: esubcate.OperationEncrypted.Decryt(),
		Category:  esubcate.CategoryEncrypted.Decrypt(),
	}
}
func (subcate Subcategory) Count() (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from Subcategory where IsDeleted=0", nil).Scan(&count)
	return count
}
func (subcate Subcategory) CountByCategoryId(id int) (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from Subcategory where IsDeleted=0 and CategoryID=?", id).Scan(&count)
	return count
}
func SubcategoryHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	if r.Method == "GET" {
		var subcate Subcategory
		cateIDFromURL := r.URL.Query().Get("id")
		cateID, err := strconv.Atoi(cateIDFromURL)
		subcate.Category.ID = 0
		var cate Category
		cates := cate.GetAllEntity()

		for i, _ := range cates {
			if cates[i].ID == cateID {
				cates[i].Selected = true
			}
		}
		if err != nil {
			if len(cates) > 0 {
				subcate.Category.ID = cates[0].ID
			}
			glog.Infof("Subcategory->convert id err: %v", err)
		} else {
			subcate.Category.ID = cateID
		}
		page := r.URL.Query().Get("page")
		count := subcate.CountByCategoryId(subcate.Category.ID)
		pagination := GetPagination(page, count)
		subcates := subcate.GetEntity(pagination)
		data := struct {
			Title         string
			Categories    []Category
			Subcategories []Subcategory
			Pagination    Pagination
			CategoryId    int
		}{
			Title:         "Subcategory",
			Categories:    cates,
			Subcategories: subcates,
			Pagination:    pagination,
			CategoryId:    subcate.Category.ID,
		}
		templates.ExecuteTemplate(w, "subcategory.gtpl", data)
	} else if r.Method == "POST" {
		cateIDForm := r.FormValue("category")
		cateID, err := strconv.Atoi(cateIDForm)
		if err != nil {
			glog.Errorf("SubcategoryHandler->convert id err: %v\n", err)
		}
		if r.FormValue("update") == "Update" {
			updatedname := r.FormValue("updatedname")
			updatedid := r.FormValue("updatedid")
			id, err := strconv.Atoi(updatedid)
			if err != nil {
				glog.Errorf("convert string %s to int err: %v", updatedid, err)
			}
			var subcate Subcategory
			subcate.ID = id
			subcate.Name = updatedname
			subcate.Operation.UpdatedTime = time.Now().Format(LongFormat)
			subcate.Operation.UpdatedBy = 0
			rowsAffected := subcate.UpdEntity()
			if rowsAffected > 0 {
				//successful
			}
		} else if r.FormValue("create") == "Create" {
			subcateName := r.FormValue("createdsubcatename")
			var subcate Subcategory
			subcate.Name = subcateName

			subcate.Category.ID = cateID
			subcate.CreatedTime = time.Now().Format(LongFormat)
			subcate.CreatedBy = 0

			lastInsertId := subcate.AddEntity()
			if lastInsertId > -1 {
				//insert successful
			}
		}
		http.Redirect(w, r, "/subcategory?id="+cateIDForm, http.StatusMovedPermanently)
	}
}
func SubcategoryDelHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	sid := r.FormValue("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		glog.Errorf("Subcategory->convert id err: %v", err)
		fmt.Fprint(w, false)
		return
	}
	var subcate Subcategory
	subcate.ID = id
	subcate.Operation.DeletedTime = time.Now().Format(LongFormat)
	subcate.Operation.DeletedBy = 0
	rowsAffected := subcate.DelEntity()
	if rowsAffected > 0 {
		//successful
		fmt.Fprint(w, true)
		return
	}
	fmt.Fprint(w, false)
}
func (item Item) GetEntity(pagination Pagination) []Item {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select ID,Store,Address,PurchasedDate,Amount,ReceiptImage,Remark,CreatedTime,CreatedBy,SubcategoryID,SubcategoryName,CategoryID,CategoryName from vw_Item where IsDeleted=0 limit ? offset ?")
	if err != nil {
		glog.Errorf("Item->GetEntity->stmt err: %v\n", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(pagination.Size, pagination.Size*(pagination.Index-1))
	if err != nil {
		glog.Errorf("Item->GetEntity->query err: %v\n", err)
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var eitem ItemEncrypted
		rows.Scan(&eitem.ID, &eitem.Store, &eitem.Address, &eitem.PurchasedDate, &eitem.Amount,
			&eitem.Receipt, &eitem.Remark, &eitem.OperationEncrypted.CreatedTime,
			&eitem.OperationEncrypted.CreatedBy, &eitem.SubcategoryEncrypted.ID,
			&eitem.SubcategoryEncrypted.Name, &eitem.SubcategoryEncrypted.CategoryEncrypted.ID,
			&eitem.SubcategoryEncrypted.CategoryEncrypted.Name)
		item := eitem.Decrypt()
		item.Receipt = mtcrypto.Base64Encode([]byte(item.Receipt))
		items = append(items, item)
	}
	return items
}
func (item Item) GetAllEntity() []Item {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->GetEntity->open db err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select ID,Store,Address,PurchasedDate,Amount,ReceiptImage,Remark,CreatedTime,CreatedBy,SubcategoryID,SubcategoryName,CategoryID,CategoryName from vw_Item where IsDeleted=0")
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
		var eitem ItemEncrypted
		rows.Scan(&eitem.ID, &eitem.Store, &eitem.Address, &eitem.PurchasedDate, &eitem.Amount,
			&eitem.Receipt, &eitem.Remark, &eitem.OperationEncrypted.CreatedTime,
			&eitem.OperationEncrypted.CreatedBy, &eitem.SubcategoryEncrypted.ID,
			&eitem.SubcategoryEncrypted.Name, &eitem.SubcategoryEncrypted.CategoryEncrypted.ID,
			&eitem.SubcategoryEncrypted.CategoryEncrypted.Name)
		item := eitem.Decrypt()
		item.Receipt = mtcrypto.Base64Encode([]byte(item.Receipt))
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
	res, err := stmt.Exec(eitem.Store, eitem.Address, eitem.PurchasedDate, eitem.Receipt,
		eitem.Remark, eitem.OperationEncrypted.CreatedTime, eitem.OperationEncrypted.CreatedBy,
		eitem.SubcategoryEncrypted.ID)
	if err != nil {
		glog.Errorf("Item->AddEntity->exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Item->AddEntity->get LastInsertId err: %v\n", err)
	}
	return id
}
func (item Item) DelEntity() int64 {
	eitem := item.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Item set IsDeleted=1,DeletedTime=?,DeletedBy=? where id=?")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(eitem.OperationEncrypted.DeletedTime, eitem.OperationEncrypted.DeletedBy, eitem.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
}
func (item Item) UpdEntity() int64 {
	var stmt *sql.Stmt
	var res sql.Result
	var err error
	eitem := item.Encrypt()

	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Item->AddEntity->open db err: %v\n", err)
	}
	defer db.Close()
	if len(item.Receipt) > 0 {
		stmt, err = db.Prepare("update Item set Store=?,Address=?,PurchasedDate=?,ReceiptImage=?,Remark=?,UpdatedTime=?,UpdatedBy=?,SubcategoryID=? where id=?")
		if err != nil {
			glog.Errorf("Item->AddEntity->stmt err: %v\n", err)
		}
		res, err = stmt.Exec(eitem.Store, eitem.Address, eitem.PurchasedDate, eitem.Receipt,
			eitem.Remark, eitem.OperationEncrypted.UpdatedTime, eitem.OperationEncrypted.UpdatedBy,
			eitem.SubcategoryEncrypted.ID, eitem.ID)
		if err != nil {
			glog.Errorf("Item->AddEntity->exec err: %v\n", err)
		}
	} else {
		stmt, err = db.Prepare("update Item set Store=?,Address=?,PurchasedDate=?,Remark=?,UpdatedTime=?,UpdatedBy=?,SubcategoryID=? where id=?")
		if err != nil {
			glog.Errorf("Item->AddEntity->stmt err: %v\n", err)
		}
		res, err = stmt.Exec(eitem.Store, eitem.Address, eitem.PurchasedDate,
			eitem.Remark, eitem.OperationEncrypted.UpdatedTime, eitem.OperationEncrypted.UpdatedBy,
			eitem.SubcategoryEncrypted.ID, eitem.ID)
		if err != nil {
			glog.Errorf("Item->AddEntity->exec err: %v\n", err)
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Item->AddEntity->get LastInsertId err: %v\n", err)
	}
	return rowsAffected
}
func (item Item) RemoveRceipt() int64 {
	eitem := item.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Item set ReceiptImage='',UpdatedTime=?,UpdatedBy=? where id=?")
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(eitem.OperationEncrypted.UpdatedTime, eitem.OperationEncrypted.UpdatedBy, eitem.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
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
	amount, err := mtcrypto.AESEncrypt(key, mtconverter.Float642String(item.Amount))
	if err != nil {
		glog.Errorf("enctypt name %v err: %v", item.Amount, err)
	}
	remark, err := mtcrypto.AESEncrypt(key, item.Remark)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", item.Remark, err)
	}

	return ItemEncrypted{
		ID:                   item.ID,
		Store:                store,
		Address:              address,
		PurchasedDate:        purdate,
		Receipt:              receipt,
		Amount:               amount,
		Remark:               remark,
		OperationEncrypted:   item.Operation.Encryt(),
		SubcategoryEncrypted: item.Subcategory.Encrypt(),
	}
}
func (eitem ItemEncrypted) Decrypt() Item {
	store, err := mtcrypto.AESDecrypt(key, eitem.Store)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Store, err)
	}
	address, err := mtcrypto.AESDecrypt(key, eitem.Address)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Address, err)
	}
	purdate, err := mtcrypto.AESDecrypt(key, eitem.PurchasedDate)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.PurchasedDate, err)
	}
	amt, err := mtcrypto.AESDecrypt(key, eitem.Amount)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Amount, err)
	}
	amount, err := mtconverter.Bytes2Float64(amt)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Amount, err)
	}
	receipt, err := mtcrypto.AESDecrypt(key, eitem.Receipt)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Receipt, err)
	}
	remark, err := mtcrypto.AESDecrypt(key, eitem.Remark)
	if err != nil {
		glog.Errorf("decrypt name %s err: %v", eitem.Remark, err)
	}
	return Item{
		ID:            eitem.ID,
		Store:         string(store),
		Address:       string(address),
		PurchasedDate: string(purdate),
		Amount:        amount,
		Receipt:       string(receipt),
		Remark:        string(remark),
		Operation:     eitem.OperationEncrypted.Decryt(),
		Subcategory:   eitem.SubcategoryEncrypted.Decrypt(),
	}
}
func (item Item) Count() (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from vw_Item where IsDeleted=0", nil).Scan(&count)
	return count
}
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	if r.Method == "GET" {

		var item Item
		var cate Category
		var subcate Subcategory

		cates := cate.GetAllEntity()
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
		subcates := subcate.GetAllEntity()
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
		page := r.URL.Query().Get("page")
		count := item.Count()
		pagination := GetPagination(page, count)
		items := item.GetEntity(pagination)
		data := struct {
			Title         string
			Items         []Item
			Categories    []Category
			Subcategories []Subcategory
			Pagination    Pagination
		}{
			Title:         "Item",
			Items:         items,
			Categories:    cates,
			Subcategories: subcates,
			Pagination:    pagination,
		}
		// itemtemplate.Execute(w, data)
		templates.ExecuteTemplate(w, "item.gtpl", data)
	} else if r.Method == "POST" {
		subcateID := r.FormValue("subcategory")
		cateID := r.FormValue("category")
		sid, err := strconv.Atoi(subcateID)
		if err != nil {
			glog.Errorf("convert sid to int err: %v \n", err)
		}
		if r.FormValue("update") == "Update" {
			itemID := r.FormValue("updatedid")
			id, err := strconv.Atoi(itemID)
			if err != nil {
				glog.Errorf("convert string %s to int err: %v", itemID, err)
			}
			// ucate := r.FormValue("updatedcategory")
			// cid, err := strconv.Atoi(ucate)
			// if err != nil {
			// 	glog.Errorf("convert string %s to int err: %v", ucate, err)
			// }
			usubcate := r.FormValue("updatedsubcategory")
			sid, err = strconv.Atoi(usubcate)
			if err != nil {
				glog.Errorf("convert string %s to int err: %v", usubcate, err)
			}
			ustore := r.FormValue("updatedstore")
			uaddr := r.FormValue("updatedaddress")
			upurdate := r.FormValue("updatedpurchaseddate")
			file, _, err := r.FormFile("updatedreceipt")
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
			uremark := r.FormValue("purchaseddateremark")

			var item Item
			item.ID = id
			item.Subcategory.ID = sid
			item.Store = ustore
			item.Address = uaddr
			item.PurchasedDate = upurdate
			item.Receipt = string(receiptData)
			item.Remark = uremark
			item.Operation.UpdatedTime = time.Now().Format(LongFormat)
			item.Operation.UpdatedBy = 0
			item.UpdEntity()
		} else if r.FormValue("create") == "Create" {

			purchasedDate := r.FormValue("createdpurchaseddate")
			store := r.FormValue("createdstore")
			address := r.FormValue("createdaddress")
			remark := r.FormValue("createdremark")
			file, _, err := r.FormFile("createdreceiptimage")
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
		}
		http.Redirect(w, r, "/item?sid="+strconv.Itoa(sid)+"&cid="+cateID, http.StatusMovedPermanently)
	}
}
func ItemDelHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	eid := r.FormValue("id")
	id, err := strconv.Atoi(eid)
	if err != nil {
		glog.Errorf("convert string %s to int err: %v", eid, err)
		fmt.Fprint(w, false)
		return
	}
	var item Item
	item.ID = id
	item.Operation.DeletedTime = time.Now().Format(LongFormat)
	item.Operation.DeletedBy = 0
	rowsAffected := item.DelEntity()
	if rowsAffected > 0 {
		//successful
		fmt.Fprint(w, true)
		return
	}
	fmt.Fprint(w, false)
}
func (detail Detail) GetEntity(pagination Pagination) []Detail {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("select ID,Name,Price,Quantity,LabelOne,LabelTwo,Remark,CreatedTime,CreatedBy from Detail where IsDeleted=0 and ItemID=? limit ? offset ?")
	if err != nil {
		glog.Errorf("db prepare err: %v\n", err)
	}
	rows, err := stmt.Query(detail.Item.ID, pageSize, pageSize*(pagination.Index-1))
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	var details []Detail
	for rows.Next() {
		var edetail DetailEncrypted
		rows.Scan(&edetail.ID, &edetail.Name, &edetail.Price, &edetail.Quantity,
			&edetail.LabelOne, &edetail.LabelTwo, &edetail.Remark,
			&edetail.CreatedTime, &edetail.CreatedBy)
		detail := edetail.Decrypt()
		detail.LabelOne = mtcrypto.Base64Encode([]byte(detail.LabelOne))
		detail.LabelTwo = mtcrypto.Base64Encode([]byte(detail.LabelTwo))

		details = append(details, detail)
	}
	return details
}
func (detail Detail) GetAllEntity() []Detail {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	stmt, err := db.Prepare("select ID,Name,Price,Quantity,LabelOne,LabelTwo,Remark,CreatedTime,CreatedBy from Detail where IsDeleted=0 and ItemID=?")
	if err != nil {
		glog.Errorf("db prepare err: %v\n", err)
	}
	rows, err := stmt.Query(detail.Item.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	var details []Detail
	for rows.Next() {
		var edetail DetailEncrypted
		rows.Scan(&edetail.ID, &edetail.Name, &edetail.Price, &edetail.Quantity,
			&edetail.LabelOne, &edetail.LabelTwo, &edetail.Remark,
			&edetail.CreatedTime, &edetail.CreatedBy)
		detail := edetail.Decrypt()
		detail.LabelOne = mtcrypto.Base64Encode([]byte(detail.LabelOne))
		detail.LabelTwo = mtcrypto.Base64Encode([]byte(detail.LabelTwo))

		details = append(details, detail)
	}
	return details
}
func (detail Detail) AddEntity() int64 {
	edetail := detail.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	////transaction begin
	tx, err := db.Begin()
	if err != nil {
		glog.Errorf("tx begin err: %v\n", err)
	}
	//insert detail
	stmt, err := tx.Prepare("insert into Detail(ItemID,Name,Price,Quantity,LabelOne,LabelTwo,Remark,CreatedTime,CreatedBy) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(edetail.ItemEncrypted.ID, edetail.Name, edetail.Price,
		edetail.Quantity, edetail.LabelOne, edetail.LabelTwo, edetail.Remark,
		edetail.CreatedTime, edetail.CreatedBy)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("get LastInsertId err: %v\n", err)
	}
	//retrive amount from item
	var eitem ItemEncrypted
	eitem.ID = detail.Item.ID
	stmt, err = tx.Prepare("select Amount from Item where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	err = stmt.QueryRow(detail.Item.ID).Scan(&eitem.Amount)
	switch {
	case err == sql.ErrNoRows:
		glog.Infof("no row err: %v\n", err)
	case err != nil:
		glog.Errorf("get item amount err: %v\n", err)
	default:
		//
	}
	item := eitem.Decrypt()
	//update amount
	item.Amount = item.Amount + GetAmount(detail.Price, detail.Quantity)
	stmt, err = tx.Prepare("update Item set Amount=? where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err = stmt.Exec(item.Encrypt().Amount, detail.Item.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	//transaction commit
	err = tx.Commit()
	if err != nil {
		id = -1
		glog.Errorf("tx commit err: %v\n", err)
		tx.Rollback()
	}
	return id
}
func (detail Detail) DelEntity() int64 {
	edetail := detail.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	////transaction begin
	tx, err := db.Begin()
	if err != nil {
		glog.Errorf("tx begin err: %v\n", err)
	}
	//delete detail
	stmt, err := tx.Prepare("update Detail set IsDeleted=1,DeletedTime=?,DeletedBy=? where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err := stmt.Exec(edetail.OperationEncrypted.DeletedTime,
		edetail.OperationEncrypted.DeletedBy, edetail.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("get LastInsertId err: %v\n", err)
	}
	//retrive amount from item
	var eitem ItemEncrypted
	stmt, err = tx.Prepare("select ItemID,Price,Quantity from Detail where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	rows, err := stmt.Query(edetail.ID)
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	for rows.Next() {
		rows.Scan(&edetail.ItemEncrypted.ID, &edetail.Price, &edetail.Quantity)
	}
	stmt, err = tx.Prepare("select Amount from Item where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	err = stmt.QueryRow(edetail.ItemEncrypted.ID).Scan(&eitem.Amount)
	switch {
	case err == sql.ErrNoRows:
		glog.Infof("no row err: %v\n", err)
	case err != nil:
		glog.Errorf("get item amount err: %v\n", err)
	default:
		//
	}
	item := eitem.Decrypt()
	detail = edetail.Decrypt()
	//update amount
	item.Amount = item.Amount - GetAmount(detail.Price, detail.Quantity)
	stmt, err = tx.Prepare("update Item set Amount=? where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err = stmt.Exec(item.Encrypt().Amount, detail.Item.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	//transaction commit
	err = tx.Commit()
	if err != nil {
		rowsAffected = -1
		glog.Errorf("tx commit err: %v\n", err)
		tx.Rollback()
	}
	return rowsAffected
}
func (detail Detail) UpdEntity() int64 {
	edetail := detail.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open db err: %v\n", err)
	}
	//transaction begin
	tx, err := db.Begin()
	if err != nil {
		glog.Errorf("tx begin err: %v\n", err)
	}
	//retrive price and quantity
	var edetail_old DetailEncrypted
	stmt, err := tx.Prepare("select ID,ItemID,Price,Quantity from Detail where id=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	err = stmt.QueryRow(edetail.ID).Scan(&edetail_old.ID, &edetail_old.ItemEncrypted.ID,
		&edetail_old.Price, &edetail_old.Quantity)
	if err != nil {
		glog.Errorf("stmt.QueryRow err: %v\n", err)
	}
	detail_old := edetail_old.Decrypt()
	//update detail
	sqlstr := "update Detail set Name=?,Price=?,Quantity=?,Remark=?,UpdatedTime=?,UpdatedBy=? "
	if len(detail.LabelOne) > 0 {
		sqlstr += ",LabelOne=? "
	}
	if len(detail.LabelTwo) > 0 {
		sqlstr += ",LabelTwo=? "
	}
	sqlstr += " where id=?"
	stmt, err = tx.Prepare(sqlstr)
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	var res sql.Result
	//var err error
	if len(detail.LabelOne) > 0 && len(detail.LabelTwo) > 0 {
		res, err = stmt.Exec(edetail.Name, edetail.Price, edetail.Quantity, edetail.Remark,
			edetail.OperationEncrypted.UpdatedTime, edetail.OperationEncrypted.UpdatedBy,
			edetail.LabelOne, edetail.LabelTwo, edetail.ID)
	} else if len(detail.LabelOne) > 0 {
		res, err = stmt.Exec(edetail.Name, edetail.Price, edetail.Quantity, edetail.Remark,
			edetail.OperationEncrypted.UpdatedTime, edetail.OperationEncrypted.UpdatedBy,
			edetail.LabelOne, edetail.ID)
	} else if len(detail.LabelTwo) > 0 {
		res, err = stmt.Exec(edetail.Name, edetail.Price, edetail.Quantity, edetail.Remark,
			edetail.OperationEncrypted.UpdatedTime, edetail.OperationEncrypted.UpdatedBy,
			edetail.LabelTwo, edetail.ID)
	} else {
		res, err = stmt.Exec(edetail.Name, edetail.Price, edetail.Quantity, edetail.Remark,
			edetail.OperationEncrypted.UpdatedTime, edetail.OperationEncrypted.UpdatedBy,
			edetail.ID)
	}
	if err != nil {
		glog.Errorf("stmt.Exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("get LastInsertId err: %v\n", err)
	}
	//retrive amount from Item
	var eitem ItemEncrypted
	stmt, err = tx.Prepare("select Amount from Item where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	err = stmt.QueryRow(detail_old.Item.ID).Scan(&eitem.Amount)
	switch {
	case err == sql.ErrNoRows:
		glog.Infof("no row err: %v\n", err)
	case err != nil:
		glog.Errorf("get item amount err: %v\n", err)
	default:
		glog.Infof("get item amount %v", eitem.Amount)
	}
	//update amount
	item := eitem.Decrypt()
	item.Amount = item.Amount - GetAmount(detail_old.Price, detail_old.Quantity) + GetAmount(detail.Price, detail.Quantity)
	stmt, err = tx.Prepare("update Item set Amount=? where ID=?")
	if err != nil {
		glog.Errorf("stmt err: %v\n", err)
	}
	res, err = stmt.Exec(item.Encrypt().Amount, detail_old.Item.ID)
	if err != nil {
		glog.Errorf("exec err: %v\n", err)
	}
	//transaction commit
	err = tx.Commit()
	if err != nil {
		rowsAffected = -1
		glog.Errorf("tx commit err: %v\n", err)
		tx.Rollback()
	}

	return rowsAffected
}
func (detail Detail) RemoveLabel(label string) int64 {
	edetail := detail.Encrypt()
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("Category->AddEntity->open sqlite err: %v\n", err)
	}
	defer db.Close()
	sqlstr := "update Detail set UpdatedTime=?,UpdatedBy=? "
	if label == "1" {
		sqlstr += ",LabelOne=''  where id=?"
	} else if label == "2" {
		sqlstr += ",LabelTwo=''  where id=?"
	}
	stmt, err := db.Prepare(sqlstr)
	if err != nil {
		glog.Errorf("Category->AddEntity->stmt err: %v\n", err)
	}
	res, err := stmt.Exec(edetail.OperationEncrypted.UpdatedTime, edetail.OperationEncrypted.UpdatedBy, edetail.ID)
	if err != nil {
		glog.Errorf("Category->AddEntity->exec err: %v\n", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("Category->AddEntity->get lastinsertid err: %v\n", err)
	}
	return rowsAffected
}
func (detail Detail) Encrypt() DetailEncrypted {
	name, err := mtcrypto.AESEncrypt(key, detail.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", detail.Name, err)
	}
	price, err := mtcrypto.AESEncrypt(key, mtconverter.Float642String(detail.Price))
	if err != nil {
		glog.Errorf("enctypt name %v err: %v", detail.Price, err)
	}
	quantity, err := mtcrypto.AESEncrypt(key, mtconverter.Int642String(detail.Quantity))
	if err != nil {
		glog.Errorf("enctypt name %d err: %v", detail.Quantity, err)
	}
	labelone, err := mtcrypto.AESEncrypt(key, detail.LabelOne)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", detail.LabelOne, err)
	}
	labeltwo, err := mtcrypto.AESEncrypt(key, detail.LabelTwo)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", detail.LabelTwo, err)
	}
	remark, err := mtcrypto.AESEncrypt(key, detail.Remark)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", detail.Remark, err)
	}
	return DetailEncrypted{
		ID:                 detail.ID,
		Name:               name,
		Price:              price,
		Quantity:           quantity,
		LabelOne:           labelone,
		LabelTwo:           labeltwo,
		Remark:             remark,
		OperationEncrypted: detail.Operation.Encryt(),
		ItemEncrypted:      detail.Item.Encrypt(),
	}
}
func (edetail DetailEncrypted) Decrypt() Detail {
	name, err := mtcrypto.AESDecrypt(key, edetail.Name)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", edetail.Name, err)
	}
	pri, err := mtcrypto.AESDecrypt(key, edetail.Price)
	if err != nil {
		glog.Errorf("enctypt name %v err: %v", edetail.Price, err)
	}
	price, err := mtconverter.Bytes2Float64(pri)
	if err != nil {
		glog.Errorf("enctypt name %v err: %v", edetail.Price, err)
	}
	quan, err := mtcrypto.AESDecrypt(key, edetail.Quantity)
	if err != nil {
		glog.Errorf("enctypt name %d err: %v", edetail.Quantity, err)
	}
	quantity, err := mtconverter.Bytes2Int(quan)
	if err != nil {
		glog.Errorf("enctypt name %d err: %v", edetail.Quantity, err)
	}
	labelone, err := mtcrypto.AESDecrypt(key, edetail.LabelOne)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", edetail.LabelOne, err)
	}
	labeltwo, err := mtcrypto.AESDecrypt(key, edetail.LabelTwo)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", edetail.LabelTwo, err)
	}
	remark, err := mtcrypto.AESDecrypt(key, edetail.Remark)
	if err != nil {
		glog.Errorf("enctypt name %s err: %v", edetail.Remark, err)
	}
	return Detail{
		ID:        edetail.ID,
		Name:      string(name),
		Price:     price,
		Quantity:  quantity,
		LabelOne:  string(labelone),
		LabelTwo:  string(labeltwo),
		Operation: edetail.OperationEncrypted.Decryt(),
		Remark:    string(remark),
		Item:      edetail.ItemEncrypted.Decrypt(),
	}
}
func (detail Detail) Count(iid int) (count int) {
	db, err := sql.Open(dbDrive, "./data.db")
	if err != nil {
		glog.Errorf("open sqlite err: %v\n", err)
	}
	db.QueryRow("select count(*) from Detail where IsDeleted=0 and ItemID=?", iid).Scan(&count)
	return count
}
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
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
		page := r.URL.Query().Get("page")
		count := detail.Count(iid)
		pagination := GetPagination(page, count)
		details := detail.GetEntity(pagination)
		data := struct {
			Title      string
			ItemID     int
			Details    []Detail
			Pagination Pagination
		}{
			Title:      "Detail",
			ItemID:     iid,
			Details:    details,
			Pagination: pagination,
		}
		// detailtemplate.Execute(w, data)
		templates.ExecuteTemplate(w, "detail.gtpl", data)
	} else if r.Method == "POST" {
		itemid := r.FormValue("itemid")
		if r.FormValue("update") == "Update" {
			var detail Detail
			detailID := r.FormValue("updatedid")
			id, err := strconv.Atoi(detailID)
			if err != nil {
				glog.Errorf("convert string %s to int err: %v", detailID, err)
			}
			detail.ID = id
			name := r.FormValue("updatedname")
			detail.Name = name
			price := r.FormValue("updatedprice")
			pri, err := strconv.ParseFloat(price, 64)
			if err != nil {
				detail.Price = 0.0
				glog.Errorf("parse float %s err: %v", price, err)
			} else {
				detail.Price = pri
			}
			quantity := r.FormValue("updatedquantity")
			quan, err := strconv.ParseInt(quantity, 10, 64)
			if err != nil {
				detail.Quantity = 1
				glog.Errorf("parse float %s err: %v", price, err)
			} else {
				detail.Quantity = quan
			}
			var labeloneData, labeltwoData []byte
			labelone, _, err := r.FormFile("updatedlone")
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
			labeltwo, _, err := r.FormFile("updatedltwo")
			switch err {
			case nil:
				labeltwoData, err = ioutil.ReadAll(labeltwo)
				if err != nil {
					glog.Errorf("read file err: %v\n", err)
				}
			case http.ErrMissingFile:
				glog.Infof("no file uploaded \n")
			default:
				glog.Errorf("upload file err: %v\n", err)
			}
			detail.LabelOne = string(labeloneData)
			detail.LabelTwo = string(labeltwoData)
			remark := r.FormValue("updatedremark")
			detail.Remark = remark
			detail.Operation.UpdatedTime = time.Now().Format(LongFormat)
			detail.Operation.UpdatedBy = 0
			rowsAffected := detail.UpdEntity()
			if rowsAffected > 0 {
				//insert successful
			}
		} else if r.FormValue("create") == "Create" {
			var detail Detail
			iid, err := strconv.Atoi(itemid)
			if err != nil {
				glog.Fatalf("get item id %s err: %v", itemid, err)
			}
			detail.Item.ID = iid
			name := r.FormValue("createdname")
			detail.Name = name
			price := r.FormValue("createdprice")
			pri, err := strconv.ParseFloat(price, 64)
			if err != nil {
				detail.Price = 0.0
				glog.Errorf("parse float %s err: %v", price, err)
			} else {
				detail.Price = pri
			}
			quantity := r.FormValue("createdquantity")
			quan, err := strconv.ParseInt(quantity, 10, 64)
			if err != nil {
				detail.Quantity = 1
				glog.Errorf("parse float %s err: %v", price, err)
			} else {
				detail.Quantity = quan
			}
			var labeloneData, labeltwoData []byte
			labelone, _, err := r.FormFile("createdlone")
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
			labeltwo, _, err := r.FormFile("createdltwo")
			switch err {
			case nil:
				labeltwoData, err = ioutil.ReadAll(labeltwo)
				if err != nil {
					glog.Errorf("read file err: %v\n", err)
				}
			case http.ErrMissingFile:
				glog.Infof("no file uploaded \n")
			default:
				glog.Errorf("upload file err: %v\n", err)
			}
			detail.LabelOne = string(labeloneData)
			detail.LabelTwo = string(labeltwoData)
			remark := r.FormValue("createdremark")
			detail.Remark = remark
			detail.CreatedTime = time.Now().Format(LongFormat)
			detail.CreatedBy = 0
			lastInsertId := detail.AddEntity()
			if lastInsertId > 0 {
				//insert successful
			}
		}
		http.Redirect(w, r, "/detail?id="+itemid, http.StatusMovedPermanently)
	}
}
func DetailDelHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	did := r.FormValue("id")
	var detail Detail
	id, err := strconv.Atoi(did)
	if err != nil {
		glog.Errorf("get detail by item id err: %v", err)
		fmt.Fprint(w, false)
		return
	}
	detail.ID = id
	detail.Operation.DeletedTime = time.Now().Format(LongFormat)
	detail.Operation.DeletedBy = 0
	rowsAffected := detail.DelEntity()
	if rowsAffected > 0 {
		//successful
		fmt.Fprint(w, true)
		return
	}
	fmt.Fprint(w, false)
}
func GetSubcategoryHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	var subcate Subcategory
	cateIDForm := r.URL.Query().Get("id")
	subcate.Category.ID = 0
	cateID, err := strconv.Atoi(cateIDForm)
	if err != nil {
		glog.Infof("convert cateID %s to int err: %v", cateIDForm, err)
	} else {
		subcate.Category.ID = cateID
	}
	subcates := subcate.GetAllEntity()
	data, err := json.Marshal(subcates)
	if err != nil {
		glog.Errorf("convert %T to json err: %v", subcates, err)
	}
	fmt.Fprint(w, string(data))
}
func GetAmount(price float64, quantity int64) float64 {
	return price * float64(quantity)
}
func RemoveReceiptHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	var item Item
	itemID := r.FormValue("id")
	id, err := strconv.Atoi(itemID)
	if err != nil {
		glog.Errorf("err :%v", err)
	}
	item.ID = id
	item.Operation.UpdatedTime = time.Now().Format(LongFormat)
	item.Operation.UpdatedBy = 0
	if item.RemoveRceipt() > 0 {
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
}
func RemoveLabelHandler(w http.ResponseWriter, r *http.Request) {
	CheckSessions(w, r)
	var detail Detail
	detailID := r.FormValue("id")
	id, err := strconv.Atoi(detailID)
	if err != nil {
		glog.Errorf("err :%v", err)
	}
	label := r.FormValue("label")
	detail.ID = id
	detail.Operation.UpdatedTime = time.Now().Format(LongFormat)
	detail.Operation.UpdatedBy = 0
	if detail.RemoveLabel(label) > 0 {
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
}
func GetIPFromRequest(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		glog.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		glog.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	return userIP.String()
}
func GetLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}
func GetPagination(page string, count int) Pagination {
	pageIndex, err := strconv.Atoi(page)
	if err != nil {
		pageIndex = 1
		glog.Infof("get page index err: %v", err)
	}

	var pagination Pagination
	pagination.Size = pageSize
	pagination.Index = pageIndex
	if count%pagination.Size == 0 {
		pagination.Count = count / pagination.Size
	} else {
		pagination.Count = count/pagination.Size + 1
	}
	pagination.Previous = pageIndex - 1
	pagination.Next = pageIndex + 1
	if pagination.Index > pagination.Count {
		pagination.Index -= 1
	}
	if pagination.Index < 1 {
		pagination.Index = 1
	}
	return pagination
}
