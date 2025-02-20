package setup

import (
	"restful/db"
)

func DeckMasters(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("deck_masters", "dmid")
	db.Uniquekey("deck_masters", "dmid")
	db.CreateIndex("deck_masters", "dmid")
	db.AlterColumnPrimaryKey("deck_masters", []string{"dmid"})
	db.CreateIndex("deck_masters", "uid")
	db.AlterColumnType("deck_masters", "app_id", "VARCHAR(255)", nil)
	db.AlterColumnType("deck_masters", "name", "VARCHAR(255)", nil)
	db.AlterColumnType("deck_masters", "create_at", "VARCHAR(20)", nil)
	db.AlterColumnType("deck_masters", "modify_at", "VARCHAR(20)", nil)
	db.AlterColumnType("deck_masters", "delete_at", "VARCHAR(20)", nil)
}

func DeckDetails(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("deck_details", "ddid")
	db.Uniquekey("deck_details", "ddid")
	db.CreateIndex("deck_details", "ddid")
	db.CreateIndex("deck_details", "dmid")
	db.AlterColumnType("deck_details", "symbol", "VARCHAR(50)", nil)
	db.AlterColumnType("deck_details", "exchange", "VARCHAR(50)", nil)
	db.AlterColumnType("deck_details", "name", "VARCHAR(255)", nil)
	db.AlterColumnType("deck_details", "price", "float8", nil)
	db.AlterColumnType("deck_details", "change", "float8", nil)
	db.AlterColumnType("deck_details", "change_rate", "float8", nil)
	db.AlterColumnType("deck_details", "create_at", "VARCHAR(20)", nil)
}
