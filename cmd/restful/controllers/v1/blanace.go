package c_v1

import (
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
)

type Balance struct {
	DB *db.DB
}

func (b *Balance) UpdateBalance(c *gin.Context) {
	var body pbm.ReqUpdateBalance

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	var users *pbm.Users
	query := `select u.account from users u where u.uid = ?`
	b.DB.Conn.Raw(query, body.Uid).Scan(&users)

	var balanceMasters *pbm.BalanceMasters
	query = `select bm.bmid from balance_masters bm where bm.account = ?`
	b.DB.Conn.Raw(query, users.Account).Scan(&balanceMasters)

	at := b.DB.Now()
	bdid := b.DB.NextSequence("balance_details", "bdid")
	balanceDetails := &pbm.BalanceDetails{
		Bdid:     bdid,
		Bmid:     balanceMasters.Bmid,
		Money:    body.Money,
		CreateAt: at,
	}
	b.DB.Conn.Create(&balanceDetails)

	query = `
	update balance_masters
	set total_money = total_money + ?,
		modify_at = ?
	where account = ? and bmid = ?
	`
	b.DB.Conn.Exec(query, body.Money, at, users.Account, balanceMasters.Bmid)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}
