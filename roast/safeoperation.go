package roast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"miego/xmap"

	"github.com/jinzhu/gorm"
)

type tabler interface {
	TableName() string
}

type UpdateInfo struct {
	Table     string
	SetArgs   xmap.Map
	WhereArgs xmap.Map
	Err       error
}

func SafeRemNew(tableName string, set xmap.Map, where xmap.Map) *UpdateInfo {
	return &UpdateInfo{
		Table:     tableName,
		SetArgs:   set,
		WhereArgs: where,
	}
}

func (u *UpdateInfo) Set(dic xmap.Map) *UpdateInfo {
	u.SetArgs = dic
	return u
}

func (u *UpdateInfo) Where(dic xmap.Map) *UpdateInfo {
	u.WhereArgs = dic
	return u
}

func (u *UpdateInfo) Exec(db *gorm.DB) *UpdateInfo {
	var setLines []string
	for k, data := range u.SetArgs {
		switch data.(type) {
		case int:
			v := reflect.ValueOf(data).Int()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case int8:
			v := reflect.ValueOf(data).Int()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case int16:
			v := reflect.ValueOf(data).Int()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case int32:
			v := reflect.ValueOf(data).Int()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case int64:
			v := reflect.ValueOf(data).Int()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint:
			v := reflect.ValueOf(data).Uint()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint8:
			v := reflect.ValueOf(data).Uint()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint16:
			v := reflect.ValueOf(data).Uint()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint32:
			v := reflect.ValueOf(data).Uint()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint64:
			v := reflect.ValueOf(data).Uint()
			setLines = append(setLines, fmt.Sprintf("`%s` = %d", k, v))

		case bool:
			v := reflect.ValueOf(data).Bool()
			setLines = append(setLines, fmt.Sprintf("`%s` = %t", k, v))

		case float32:
			v := reflect.ValueOf(data).Float()
			setLines = append(setLines, fmt.Sprintf("`%s` = %f", k, v))

		case float64:
			v := reflect.ValueOf(data).Float()
			setLines = append(setLines, fmt.Sprintf("`%s` = %f", k, v))

		case string:
			v := reflect.ValueOf(data).String()
			setLines = append(setLines, fmt.Sprintf("`%s` = \"%s\"", k, v))
		}
	}

	var whereLines []string
	for k, data := range u.WhereArgs {
		switch data.(type) {
		case int:
			v := reflect.ValueOf(data).Int()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case int8:
			v := reflect.ValueOf(data).Int()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case int16:
			v := reflect.ValueOf(data).Int()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case int32:
			v := reflect.ValueOf(data).Int()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case int64:
			v := reflect.ValueOf(data).Int()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint:
			v := reflect.ValueOf(data).Uint()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint8:
			v := reflect.ValueOf(data).Uint()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint16:
			v := reflect.ValueOf(data).Uint()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint32:
			v := reflect.ValueOf(data).Uint()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case uint64:
			v := reflect.ValueOf(data).Uint()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %d", k, v))

		case bool:
			v := reflect.ValueOf(data).Bool()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %t", k, v))

		case float32:
			v := reflect.ValueOf(data).Float()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %f", k, v))

		case float64:
			v := reflect.ValueOf(data).Float()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = %f", k, v))

		case string:
			v := reflect.ValueOf(data).String()
			whereLines = append(whereLines, fmt.Sprintf("`%s` = \"%s\"", k, v))
		}
	}

	s := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", u.Table, strings.Join(setLines, ", "), strings.Join(whereLines, " AND "))

	if err := db.Exec(s).Error; err != nil {
		u.Err = err
	}
	return u
}

func (u *UpdateInfo) Error() error {
	return u.Err
}

func SafeRem(db *gorm.DB, tableName string, RemBy string, RemWhy int, where xmap.Map) error {
	if where == nil {
		return nil
	}

	now, _ := strconv.ParseUint(time.Now().Format("20060102150405"), 0, 64)
	set := xmap.Make(
		"RemAt", now,
		"RemBy", RemBy,
		"RemWhy", RemWhy,
	)
	where.SafeMerge(xmap.Make("RemAt", 0))
	return SafeRemNew(tableName, set, where).Exec(db).Error()
}

func SafeAdd(db *gorm.DB, Object interface{}, RemBy string, where xmap.Map) error {
	tx := db.Begin()

	var tableName string
	if tabler, ok := Object.(tabler); ok {
		tableName = tabler.TableName()
	}

	if err := SafeRem(db, tableName, RemBy, RemWhy_Update, where); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(Object).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// rds.Publish("RPS.A."+tableName,  "CrtAt")

	return nil
}
