package setup

import (
	"restful/db"
)

func Ratios(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("ratios", "rid")
	db.Uniquekey("ratios", "rid")
	db.AlterColumnPrimaryKey("ratios", []string{"rid", "region"})
	db.CreateIndex("ratios", "rid")
	db.CreateIndex("ratios", "create_at")
	db.AlterColumnType("ratios", "region", "VARCHAR(20)", nil)
	db.AlterColumnType("ratios", "create_at", "VARCHAR(20)", nil)
}
