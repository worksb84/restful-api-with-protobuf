package c_v1

import (
	"log"
	"net/http"
	"pbm"
	"regexp"
	"restful/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Purchase struct {
	DB *db.DB
}

func (u *Purchase) Purchase(c *gin.Context) {
	var body pbm.ReqPurchases

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := u.DB.Now()
	pid := u.DB.NextSequence("purchases", "pid")
	purchases := &pbm.Purchases{
		Pid:       pid,
		Uid:       body.Uid,
		ProductId: body.ProductId,
		Price:     body.Price,
		CreateAt:  at,
	}
	u.DB.Conn.Create(&purchases)

	var users *pbm.Users
	u.DB.Conn.Where(&pbm.Users{Uid: body.Uid}).First(&users)

	reStr := regexp.MustCompile("[a-zA-Z]+")
	productType := reStr.FindAllString(body.ProductId, -1)[0]
	reNum := regexp.MustCompile("[0-9]+")
	productValue, _ := strconv.Atoi(reNum.FindAllString(body.ProductId, -1)[0])

	if productType == "ls" {
		// If Subscribe
		users.IsSubscribe = "Y"

		subscribeAt := u.DB.StringToTime(at)
		users.SubscribeAt = at

		if u.DB.StringToTime(users.SubscribeEndAt).After(subscribeAt) && users.IsSubscribe == "Y" {
			users.SubscribeEndAt = u.DB.AddDays(users.SubscribeEndAt, productValue)
		} else {
			users.SubscribeEndAt = u.DB.AddDays(at, productValue)
		}

		u.DB.Conn.Where("uid = ?", body.Uid).Save(users)
	} else {
		var users *pbm.Users
		query := `select u.account from users u where u.uid = ?`
		u.DB.Conn.Raw(query, body.Uid).Scan(&users)

		var balanceMasters *pbm.BalanceMasters
		query = `select bm.bmid from balance_masters bm where bm.account = ?`
		u.DB.Conn.Raw(query, users.Account).Scan(&balanceMasters)

		at := u.DB.Now()
		bdid := u.DB.NextSequence("balance_details", "bdid")
		balanceDetails := &pbm.BalanceDetails{
			Bdid:     bdid,
			Bmid:     balanceMasters.Bmid,
			Money:    float32(productValue) * 1000,
			CreateAt: at,
		}
		u.DB.Conn.Create(&balanceDetails)

		query = `
		update
			balance_masters
		set
			total_money = total_money + ?,
			modify_at = ?
		where
			account = ?
			and bmid = ?
		`
		u.DB.Conn.Exec(query, float32(productValue)*1000, at, users.Account, balanceMasters.Bmid)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, nil)

}
