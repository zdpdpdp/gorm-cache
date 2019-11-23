## gcache

gorm 简单的缓存支持

## Install

```
go get -u github.com/zdpdpdp/gcache
```

## Support

- 普通查询
- preload 查询

## Usage

```go
import . "github.com/zdpdpdp/gcache"

addr := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
db, e := gorm.Open("mysql", addr)
if e != nil {
    panic(e)
}
plugin := NewPlugin(NewMemoryDriver())  //内置了不完善的内存驱动,建议自定义 driver
plugin.Bind(db)

var users []user
//使用 cache_key 作为缓存 key, 缓存一小时
db.Set(NewCacheParam("cache_key", time.Hour)).Table("users").Preload("Books").Find(&users)
```