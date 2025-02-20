package setup

import (
	"restful/db"
)

func Purchases(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("purchases", "pid")
	db.Uniquekey("purchases", "pid")
	db.CreateIndex("purchases", "pid")
	db.CreateIndex("purchases", "uid")
	db.AlterColumnPrimaryKey("purchases", []string{"pid"})
	db.AlterColumnType("purchases", "product_id", "VARCHAR(50)", nil)
	db.AlterColumnType("purchases", "price", "float8", nil)
	db.AlterColumnType("purchases", "create_at", "varchar(20)", nil)
}
