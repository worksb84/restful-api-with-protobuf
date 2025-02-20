package setup

import (
	"restful/db"
)

func GameRecordMasters(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("game_record_masters", "grmid")
	db.Uniquekey("game_record_masters", "grmid")
	db.CreateIndex("game_record_masters", "grmid")
	db.AlterColumnPrimaryKey("game_record_masters", []string{"grmid"})
	db.AlterColumnType("game_record_masters", "grtype", "VARCHAR(10)", nil)
	db.AlterColumnType("game_record_masters", "rid", "VARCHAR(255)", nil)
	db.AlterColumnType("game_record_masters", "app_id", "VARCHAR(255)", nil)
	db.AlterColumnType("game_record_masters", "create_at", "VARCHAR(20)", nil)
	db.AlterColumnType("game_record_masters", "modify_at", "VARCHAR(20)", nil)
	db.AlterColumnType("game_record_masters", "delete_at", "VARCHAR(20)", nil)
}

func GameRecordDetails(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("game_record_details", "grdid")
	db.Uniquekey("game_record_details", "grdid")
	db.AlterColumnPrimaryKey("game_record_details", []string{"grdid"})
	db.CreateIndex("game_record_details", "grdid")
	db.CreateIndex("game_record_details", "grmid")
	db.CreateIndex("game_record_details", "uid")
	db.AlterColumnType("game_record_details", "record", "VARCHAR(20)", nil)
}
