package c_v1

import (
	"errors"
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	DB *db.DB
}

func (u *Users) Login(c *gin.Context) {
	var body pbm.ReqLogin

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	var user *pbm.Users
	if err := u.DB.Conn.Where(&pbm.Users{Email: body.Email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			at := u.DB.Now()
			uid := u.DB.NextSequence("users", "uid")
			account := uuid.New().String()
			user = &pbm.Users{
				Uid:         uid,
				Email:       body.Email,
				Nickname:    body.Nickname,
				Image:       body.Image,
				Account:     account,
				IsSubscribe: "N",
				LoginAt:     at,
				CreateAt:    at,
				ModifyAt:    at,
			}
			u.DB.Conn.Create(&user)

			bmid := u.DB.NextSequence("balance_masters", "bmid")
			balanceMaster := &pbm.BalanceMasters{
				Bmid:       bmid,
				Uid:        uid,
				Account:    account,
				TotalMoney: 10000000,
				CreateAt:   at,
				ModifyAt:   at,
			}
			u.DB.Conn.Create(&balanceMaster)

			bdid := u.DB.NextSequence("balance_details", "bdid")
			balanceDetails := &pbm.BalanceDetails{
				Bdid:     bdid,
				Bmid:     bmid,
				Money:    10000000,
				CreateAt: at,
			}
			u.DB.Conn.Create(&balanceDetails)
		}
	}

	subscribeEndAt := u.DB.StringToTime(user.SubscribeEndAt)
	now := u.DB.StringToTime(u.DB.Now())
	if subscribeEndAt.Before(now) {
		user.IsSubscribe = "N"
	}

	user.Nickname = body.Nickname
	user.Image = body.Image
	user.LoginAt = u.DB.Now()

	u.DB.Conn.Where("uid = ?", user.Uid).Save(user)

	query := `
	select
		u.uid,
		u.email,
		u.nickname,
		u.image,
		u.account,
		bm.total_money,
		u.is_subscribe,
		u.subscribe_at,
		u.subscribe_end_at
	from
		users u
	left join balance_masters bm on
		u.uid = bm.uid
		and u.account = bm.account
	where
		u.email = ?
	`

	var resLogin *pbm.ResLogin
	u.DB.Conn.Raw(query, body.Email).Scan(&resLogin)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resLogin)
}

func (u *Users) Profile(c *gin.Context) {
	var body pbm.ReqByUID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `
	select
		g1.uid,
		g1.email,
		g1.nickname,
		g1.image,
		g1.account,
		g1.total_money,
		g1.is_subscribe,
		g1.subscribe_at,
		g1.subscribe_end_at,
		max(g1.win) + max(g1.defeat) + max(g1.draw) games,
		max(g1.win) win,
		max(g1.defeat) defeat,
		max(g1.draw) draw
	from
		(
		select
			u.uid,
			u.email,
			u.nickname,
			u.image,
			u.account,
			bm.total_money,
			u.is_subscribe,
			u.subscribe_at,
			u.subscribe_end_at,
			coalesce (case
				when grd.record = 'W' then count(grd.record)
			end,
			0) win,
			coalesce (case
				when grd.record = 'F' then count(grd.record)
			end,
			0) defeat,
			coalesce (case
				when grd.record = 'D' then count(grd.record)
			end,
			0) draw
		from
			users u
		left join balance_masters bm on
			u.uid = bm.uid
			and u.account = bm.account
		left join game_record_details grd on
			u.uid = grd.uid
		where
			u.uid = ?
		group by
			u.uid,
			u.email,
			u.nickname,
			u.image,
			u.account,
			bm.total_money,
			u.is_subscribe,
			u.subscribe_at,
			u.subscribe_end_at,
			grd.record
	) as g1
	group by 
		g1.uid,
		g1.email,
		g1.nickname,
		g1.image,
		g1.account,
		g1.total_money,
		g1.is_subscribe,
		g1.subscribe_at,
		g1.subscribe_end_at
	`

	var resProfile *pbm.ResProfile
	u.DB.Conn.Raw(query, body.Uid).Scan(&resProfile)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resProfile)
}
