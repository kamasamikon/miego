package roast

import (
	"github.com/kamasamikon/miego/klog"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var sqlMap map[string]*sql.DB
var ormMap map[string]*gorm.DB

func SQL(prefix string) *sql.DB {
	var db *sql.DB
	var ok bool

	if prefix == "" {
		prefix = "db/my"
	}

	add := func(prefix string) *sql.DB {
		db, err := sql.Open("mysql", DSN(prefix))
		if err != nil {
			klog.E(err.Error())
			return nil
		}
		sqlMap[prefix] = db
		return db
	}

	if db, ok = sqlMap[prefix]; !ok {
		return add(prefix)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		klog.E(err.Error())
		return add(prefix)
	}

	return db
}

func ORM(prefix string) *gorm.DB {
	var db *gorm.DB
	var ok bool

	if prefix == "" {
		prefix = "db/my"
	}

	add := func(prefix string) *gorm.DB {
		db := OpenByConf(prefix)
		ormMap[prefix] = db
		return db
	}

	if db, ok = ormMap[prefix]; !ok {
		return add(prefix)
	}

	return db
}

func Hey() {
	if sqlMap == nil {
		sqlMap = make(map[string]*sql.DB)
	}
	if ormMap == nil {
		ormMap = make(map[string]*gorm.DB)
	}
}

func Bye() {
	for _, v := range sqlMap {
		v.Close()
	}
	for _, v := range ormMap {
		v.Close()
	}
}

func init() {
	Hey()
}
