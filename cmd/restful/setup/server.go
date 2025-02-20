package setup

import (
	"restful/db"
)

func Servers(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("servers", "sid")
	db.Uniquekey("servers", "sid")
	db.CreateIndex("servers", "sid")
	db.AlterColumnPrimaryKey("servers", []string{"sid"})
	db.AlterColumnType("servers", "app_id", "varchar(255)", nil)
	db.AlterColumnType("servers", "uri", "varchar(255)", nil)
	db.AlterColumnType("servers", "region", "varchar(255)", nil)
	db.AlterColumnType("servers", "sort", "int4", nil)
	db.AlterColumnType("servers", "create_at", "varchar(20)", nil)
}
