package setup

import (
	"restful/db"
)

func SnapshotLogs(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("snapshot_logs", "slid")
	db.Uniquekey("snapshot_logs", "slid")
	db.AlterColumnPrimaryKey("snapshot_logs", []string{"slid", "region"})
	db.CreateIndex("snapshot_logs", "slid")
	db.CreateIndex("snapshot_logs", "create_at")
	db.AlterColumnType("snapshot_logs", "region", "VARCHAR(20)", nil)
	db.AlterColumnType("snapshot_logs", "create_at", "VARCHAR(20)", nil)
}
