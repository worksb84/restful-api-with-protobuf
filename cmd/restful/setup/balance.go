package setup

import (
	"restful/db"
)

func BalanceMasters(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("balance_masters", "bmid")
	db.Uniquekey("balance_masters", "bmid")
	db.CreateIndex("balance_masters", "bmid")
	db.CreateIndex("balance_masters", "uid")
	db.CreateIndex("balance_masters", "account")
	db.AlterColumnPrimaryKey("balance_masters", []string{"bmid", "uid", "account"})
	db.AlterColumnType("balance_masters", "account", "varchar(255)", nil)
	db.AlterColumnType("balance_masters", "total_money", "float8", nil)
	db.AlterColumnType("balance_masters", "create_at", "varchar(20)", nil)
	db.AlterColumnType("balance_masters", "modify_at", "varchar(20)", nil)
	db.AlterColumnType("balance_masters", "delete_at", "varchar(20)", nil)
}

func BalanceDetails(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("balance_details", "bdid")
	db.Uniquekey("balance_details", "bdid")
	db.CreateIndex("balance_details", "bdid")
	db.AlterColumnPrimaryKey("balance_details", []string{"bdid"})
	db.AlterColumnType("balance_details", "money", "float8", nil)
	db.AlterColumnType("balance_details", "create_at", "varchar(20)", nil)
}
