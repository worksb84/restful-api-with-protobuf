package setup

import (
	"restful/db"
)

func Favorities(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("favorities", "fid")
	db.Uniquekey("favorities", "fid")
	db.CreateIndex("favorities", "fid")
	db.AlterColumnPrimaryKey("favorities", []string{"fid", "uid", "symbol"})
	db.CreateIndex("favorities", "uid")
	db.AlterColumnType("favorities", "symbol", "VARCHAR(50)", nil)
	db.AlterColumnType("favorities", "exchange", "VARCHAR(50)", nil)
	db.AlterColumnType("favorities", "name", "VARCHAR(255)", nil)
	db.AlterColumnType("favorities", "price", "float8", nil)
	db.AlterColumnType("favorities", "change", "float8", nil)
	db.AlterColumnType("favorities", "change_rate", "float8", nil)
	db.AlterColumnType("favorities", "create_at", "varchar(20)", nil)
}
