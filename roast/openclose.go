package roast

import (
	"fmt"
	"os"
	"time"

	"miego/conf"
	"miego/klog"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var map_DB_Conf map[*gorm.DB]string

func ConfName(db *gorm.DB) string {
	if c, ok := map_DB_Conf[db]; ok {
		return c
	}
	return ""
}

func DSNByDB(db *gorm.DB) string {
	return DSN(ConfName(db))
}

func DSN(confprefix string) string {
	if confprefix == "" {
		confprefix = "db/my"
	}

	dbDatabase := conf.SGet(confprefix+"/database", "gene")

	dbUser := conf.SGet(confprefix+"/user/name", "root")
	dbPass := conf.SGet(confprefix+"/user/pass", "root")

	dbHost := conf.SGet(confprefix+"/addr/host", os.Getenv("DOCKER_GATEWAY"))
	dbPort := conf.SGet(confprefix+"/addr/port", "3306")

	return fmt.Sprintf("%s:%s@(%s:%s)/%s?collation=utf8mb4_general_ci&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbDatabase)
}

func Open(db string, user string, pass string, host string, port string, verbose bool) *gorm.DB {
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?collation=utf8mb4_general_ci&parseTime=True&loc=Local", user, pass, host, port, db)
	klog.D("ARGS: %s", args)

	x, err := gorm.Open("mysql", args)
	if err != nil {
		klog.BT(20, "%s", err.Error())
		return nil
	}

	x.LogMode(verbose)
	x.SingularTable(true)

	x.DB().SetMaxIdleConns(20)
	x.DB().SetMaxOpenConns(50)
	x.DB().SetConnMaxLifetime(time.Second * 300)

	return x
}

func CreateTable(db *gorm.DB, models ...interface{}) {
	if db != nil {
		for _, model := range models {
			if !db.HasTable(model) {
				db.CreateTable(model)
			}
		}
	}
}

func OpenByConf(confprefix string, models ...interface{}) *gorm.DB {
	if confprefix == "" {
		confprefix = "db/my"
	}

	dbDatabase := conf.SGet(confprefix+"/database", "gene")

	dbUser := conf.SGet(confprefix+"/user/name", "root")
	dbPass := conf.SGet(confprefix+"/user/pass", "root")

	dbHost := conf.SGet(confprefix+"/addr/host", os.Getenv("DOCKER_GATEWAY"))
	dbPort := conf.SGet(confprefix+"/addr/port", "3306")

	verbose := conf.BFalse(confprefix+"/verbose")

	db := Open(dbDatabase, dbUser, dbPass, dbHost, dbPort, verbose)
	CreateTable(db, models...)

	map_DB_Conf[db] = confprefix
	return db
}

func init() {
	map_DB_Conf = make(map[*gorm.DB]string)
}
