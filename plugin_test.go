package gcache

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
	"time"
)

type user struct {
	Id    int
	Books []book
}

type book struct {
	Id     int
	UserId int
}

func (user) TableName() string {
	return "users"
}

func (book) TableName() string {
	return "books"
}

var db *gorm.DB

func init() {
	addr := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	_db, e := gorm.Open("mysql", addr)
	if e != nil {
		panic(e)
	}
	plugin := NewPlugin(NewMemoryDriver())
	plugin.Bind(_db)
	_db.LogMode(true)
	db = _db
}

func TestPluginBind(t *testing.T) {

	if db.Callback().Query().Get(CacheBeforeQuery) == nil {
		t.Error("cache before query does not register")
	}

	if db.Callback().Query().Get(CacheAfterQuery) == nil {
		t.Error("cache after query does not register")
	}
}

func TestSelect(t *testing.T) {

	db.Exec("truncate users")
	db.Create(user{
		Id: 1,
	})
	var users []user
	db.Table("users").Find(&users)
	if len(users) == 0 {
		t.Error("formal select error")
	}
}

func TestCacheSelect(t *testing.T) {

	db.Exec("truncate users")
	db.Create(user{Id: 1})
	db.Create(user{Id: 2})

	var users []user
	//第一次查询, 加载如缓存
	db.Set(CacheOption, CacheParam{ttl: time.Hour * 10, Key: "test_key"}).Table("users").Find(&users)
	if len(users) != 2 {
		t.Error("formal select1 error")
	}

	db.Table("users").Delete(user{})

	//从缓存中查出
	db.Set(CacheOption, CacheParam{
		ttl: time.Hour * 10,
		Key: "test_key",
	}).Table("users").Find(&users)
	if len(users) != 2 {
		t.Error("cache select error")
	}

	//不使用缓存
	db.Table("users").Find(&users)
	if len(users) != 0 {
		t.Error("formal select2 error")
	}

	db.Create(user{Id: 3})
	db.Table("users").Find(&users)
	if len(users) != 1 {
		t.Error("formal select3 error")
	}

}

func TestPreload(t *testing.T) {
	db.Exec("truncate users")
	db.Exec("truncate books")

	var users []user
	db.Create(user{Id: 1})
	db.Create(book{Id: 1, UserId: 1})
	db.Set(NewCacheParam("test_key_2", time.Hour)).Table("users").Preload("Books").Find(&users)
	if len(users) != 1 && len(users[0].Books) != 1 {
		t.Error("formal select error")
	}
	db.Exec("truncate books")
	db.Set(CacheOption, CacheParam{ttl: time.Hour, Key: "test_key2"}).Table("users").Preload("Books").Find(&users)
	if len(users) != 1 && len(users[0].Books) != 1 {
		t.Error("formal select error")
	}
}
