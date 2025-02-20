package setup

import (
	"restful/db"
)

func PokerGameMasters(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("poker_game_masters", "pkgmid")
	db.Uniquekey("poker_game_masters", "pkgmid")
	db.CreateIndex("poker_game_masters", "pkgmid")
	db.CreateIndex("poker_game_masters", "grmid")
	db.AlterColumnPrimaryKey("poker_game_masters", []string{"pkgmid"})
	db.AlterColumnType("poker_game_masters", "total_betting_price", "float8", nil)
	db.AlterColumnType("poker_game_masters", "create_at", "VARCHAR(20)", nil)
	db.AlterColumnType("poker_game_masters", "modify_at", "VARCHAR(20)", nil)
	db.AlterColumnType("poker_game_masters", "delete_at", "VARCHAR(20)", nil)
}

func PokerGameDetails(db *db.DB, model interface{}) {
	db.CreateTable(model)
	db.AutoIncrement("poker_game_details", "pkgdbid")
	db.Uniquekey("poker_game_details", "pkgdbid")
	db.CreateIndex("poker_game_details", "pkgdbid")
	db.CreateIndex("poker_game_details", "pkgmid")
	db.AlterColumnPrimaryKey("poker_game_details", []string{"pkgdbid"})
	db.CreateIndex("poker_game_details", "uid")
	db.AlterColumnType("poker_game_details", "betting_price", "float8", nil)
	db.AlterColumnType("poker_game_details", "reward_price", "float8", nil)
}
