package gcache

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

const CacheBeforeQuery = "CACHE:BEFORE_QUERY"
const CacheAfterQuery = "CACHE:AFTER_QUERY"
const CacheOption = "gorm:cache"

//Plugin
type Plugin struct {
	cacheDriver CacheDriver
}

//NewPlugin
func NewPlugin(driver CacheDriver) Plugin {
	return Plugin{
		cacheDriver: driver,
	}
}

//Bind
func (p Plugin) Bind(db *gorm.DB) {

	db.Callback().Query().Before("gorm:query").Register(CacheBeforeQuery, p.beforeQueryInvoke)
	db.Callback().Query().After("gorm:after_query").Register(CacheAfterQuery, p.afterQueryInvoke)

}

func (p Plugin) beforeQueryInvoke(scope *gorm.Scope) {
	cacheOption, ok := scope.Get(CacheOption)
	if !ok {
		return
	}
	param, ok := cacheOption.(CacheParam)
	if !ok {
		scope.Log("cache param error")
		return
	}
	val, ok, err := p.cacheDriver.Get(param.Key)
	if err != nil {
		scope.Log(errors.Wrapf(err, "gcache plugin get key %s fail", param.Key))
		return
	}
	if ok {
		scope.Value = val
		scope.SkipLeft()
	}

}

func (p Plugin) afterQueryInvoke(scope *gorm.Scope) {
	cacheOption, ok := scope.Get(CacheOption)
	if !ok {
		return
	}
	param, ok := cacheOption.(CacheParam)
	if !ok {
		scope.Log("cache param error")
		return
	}

	if scope.HasError() {
		return
	}

	reflectValue := scope.IndirectValue()
	if reflectValue.IsNil() || reflectValue.Len() == 0 {
		return
	}

	if err := p.cacheDriver.Set(param.Key, scope.Value, param.ttl); err != nil {
		scope.Log(errors.Wrapf(err, "gcache plugin set key %s fail", param.Key))
		return
	}
}
