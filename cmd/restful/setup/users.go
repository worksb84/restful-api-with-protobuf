package setup

import (
	"restful/db"
)

func Users(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("users", "uid")
	db.Uniquekey("users", "uid")
	db.CreateIndex("users", "uid")
	db.AlterColumnPrimaryKey("users", []string{"uid", "email"})
	db.AlterColumnType("users", "email", "VARCHAR(255)", nil)
	db.AlterColumnType("users", "nickname", "VARCHAR(50)", nil)
	db.AlterColumnType("users", "image", "VARCHAR(255)", nil)
	db.AlterColumnType("users", "account", "VARCHAR(255)", nil)
	db.AlterColumnType("users", "is_subscribe", "VARCHAR(1)", "'N'")
	db.AlterColumnType("users", "subscribe_at", "VARCHAR(20)", nil)
	db.AlterColumnType("users", "subscribe_end_at", "VARCHAR(20)", nil)
	db.AlterColumnType("users", "login_at", "VARCHAR(20)", nil)
	db.AlterColumnType("users", "create_at", "VARCHAR(20)", nil)
	db.AlterColumnType("users", "modify_at", "VARCHAR(20)", nil)
	db.AlterColumnType("users", "delete_at", "VARCHAR(20)", nil)
}
