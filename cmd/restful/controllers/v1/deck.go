package c_v1

import (
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
)

type Deck struct {
	DB *db.DB
}

func (d *Deck) GetDeckList(c *gin.Context) {
	var body pbm.ReqByAppIDandUID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	var reqresDecks []*pbm.ReqResDeck

	query := `select * from deck_masters dm where dm.app_id = ? and dm.uid = ? and (dm.delete_at is null or dm.delete_at = '')`
	var deckMasters []*pbm.DeckMasters
	d.DB.Conn.Raw(query, body.AppId, body.Uid).Scan(&deckMasters)

	for _, v := range deckMasters {
		var reqresDeckDetails []*pbm.ReqResDeckDetails
		query = `select * from deck_details dd where dd.dmid = ?`
		var deckDetails []*pbm.DeckDetails
		d.DB.Conn.Raw(query, v.Dmid).Scan(&deckDetails)

		for _, w := range deckDetails {
			reqresDeckDetail := &pbm.ReqResDeckDetails{
				Ddid:       w.Ddid,
				Dmid:       w.Dmid,
				Symbol:     w.Symbol,
				Exchange:   w.Exchange,
				Name:       w.Name,
				Price:      w.Price,
				Change:     w.Change,
				ChangeRate: w.ChangeRate,
				CreateAt:   w.CreateAt,
			}
			reqresDeckDetails = append(reqresDeckDetails, reqresDeckDetail)
		}

		reqresDeck := &pbm.ReqResDeck{
			Dmid:        v.Dmid,
			Uid:         v.Uid,
			AppId:       v.AppId,
			Name:        v.Name,
			CreateAt:    v.CreateAt,
			ModifyAt:    v.ModifyAt,
			DeleteAt:    v.DeleteAt,
			DeckDetails: reqresDeckDetails,
		}
		reqresDecks = append(reqresDecks, reqresDeck)
	}

	resDeckList := &pbm.ResDeckList{
		Decklist: reqresDecks,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resDeckList)
}

func (d *Deck) GetDeck(c *gin.Context) {
	var body pbm.ReqByDeckMastersID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from deck_masters dm where dm.dmid = ?`
	var deckMaster *pbm.DeckMasters
	d.DB.Conn.Raw(query, body.Dmid).Scan(&deckMaster)

	var reqresDeckDetails []*pbm.ReqResDeckDetails
	query = `select * from deck_details dd where dd.dmid  = ?`
	var deckDetails []*pbm.DeckDetails
	d.DB.Conn.Raw(query, deckMaster.Dmid).Scan(&deckDetails)

	for _, w := range deckDetails {
		reqresDeckDetail := &pbm.ReqResDeckDetails{
			Ddid:       w.Ddid,
			Dmid:       w.Dmid,
			Symbol:     w.Symbol,
			Exchange:   w.Exchange,
			Name:       w.Name,
			Price:      w.Price,
			Change:     w.Change,
			ChangeRate: w.ChangeRate,
			CreateAt:   w.CreateAt,
		}
		reqresDeckDetails = append(reqresDeckDetails, reqresDeckDetail)
	}

	reqresDeck := &pbm.ReqResDeck{
		Dmid:        deckMaster.Dmid,
		Uid:         deckMaster.Uid,
		AppId:       deckMaster.AppId,
		Name:        deckMaster.Name,
		CreateAt:    deckMaster.CreateAt,
		ModifyAt:    deckMaster.ModifyAt,
		DeleteAt:    deckMaster.DeleteAt,
		DeckDetails: reqresDeckDetails,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, reqresDeck)
}

func (d *Deck) UpdateDeck(c *gin.Context) {
	var body pbm.ReqResDeck

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := d.DB.Now()
	query := `
	update deck_masters set name = ?, modify_at = ? where dmid = ?`
	d.DB.Conn.Exec(query, body.Name, at, body.Dmid)

	query = `delete from deck_details dd where dd.dmid = ?`
	d.DB.Conn.Exec(query, body.Dmid)
	for _, v := range body.DeckDetails {
		ddid := d.DB.NextSequence("deck_details", "ddid")
		deckDetails := &pbm.DeckDetails{
			Ddid:       ddid,
			Dmid:       v.Dmid,
			Symbol:     v.Symbol,
			Exchange:   v.Exchange,
			Name:       v.Name,
			Price:      v.Price,
			Change:     v.Change,
			ChangeRate: v.ChangeRate,
			CreateAt:   at,
		}
		d.DB.Conn.Create(&deckDetails)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}

func (d *Deck) DeleteDeck(c *gin.Context) {
	var body pbm.ReqByDeckMastersID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := d.DB.Now()

	query := `update deck_masters set delete_at = ? where dmid = ?`
	d.DB.Conn.Exec(query, at, body.Dmid)
}

func (d *Deck) AddDeck(c *gin.Context) {
	var body pbm.ReqResDeck

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := d.DB.Now()
	dmid := d.DB.NextSequence("deck_masters", "dmid")
	deckMasters := &pbm.DeckMasters{
		Dmid:     dmid,
		Uid:      body.Uid,
		AppId:    body.AppId,
		Name:     body.Name,
		CreateAt: at,
		ModifyAt: at,
	}
	d.DB.Conn.Create(&deckMasters)

	for _, v := range body.DeckDetails {
		ddid := d.DB.NextSequence("deck_details", "ddid")
		deckDetails := &pbm.DeckDetails{
			Ddid:       ddid,
			Dmid:       dmid,
			Symbol:     v.Symbol,
			Exchange:   v.Exchange,
			Name:       v.Name,
			Price:      v.Price,
			Change:     v.Change,
			ChangeRate: v.ChangeRate,
			CreateAt:   at,
		}
		d.DB.Conn.Create(&deckDetails)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}
