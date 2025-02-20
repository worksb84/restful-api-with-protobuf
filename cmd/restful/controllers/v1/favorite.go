package c_v1

import (
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
)

type Favorite struct {
	DB *db.DB
}

func (u *Favorite) AddFavorite(c *gin.Context) {
	var body pbm.ReqResFavorities

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from favorities f where f.symbol = ? and uid = ?`
	u.DB.Conn.Exec(query, body.Symbol, body.Uid)

	at := u.DB.Now()
	fid := u.DB.NextSequence("favorities", "fid")
	favorities := &pbm.Favorities{
		Fid:        fid,
		Uid:        body.Uid,
		Symbol:     body.Symbol,
		Exchange:   body.Exchange,
		Name:       body.Name,
		Price:      body.Price,
		Change:     body.Change,
		ChangeRate: body.ChangeRate,
		CreateAt:   at,
	}
	u.DB.Conn.Create(&favorities)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}

func (u *Favorite) DeleteFavorite(c *gin.Context) {
	var body pbm.ReqBySymbolAndUID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}
	query := `delete from favorities where uid = ? and symbol = ?`
	u.DB.Conn.Exec(query, body.Uid, body.Symbol)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}

func (u *Favorite) GetFavoriteList(c *gin.Context) {
	var body pbm.ReqByUID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from favorities f where f.uid = ?`
	var reqresFavorities []*pbm.ReqResFavorities
	u.DB.Conn.Raw(query, body.Uid).Scan(&reqresFavorities)

	resFavoriteList := &pbm.ResFavoriteList{
		FavoriteList: reqresFavorities,
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resFavoriteList)
}
