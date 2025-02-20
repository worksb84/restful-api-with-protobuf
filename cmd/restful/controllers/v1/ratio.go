package c_v1

import (
	"encoding/json"
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
)

type Ratio struct {
	DB *db.DB
}

func (r *Ratio) GetSymbol(c *gin.Context) {
	var body pbm.ReqBySymbol
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from ratios where region = ? order by rid desc limit 1`

	var reqresRatios *pbm.ReqResRatios

	r.DB.Conn.Raw(query, body.Region).Scan(&reqresRatios)

	var ratioMap map[string]*pbm.ResRatios
	json.Unmarshal([]byte(reqresRatios.Ratio), &ratioMap)
	ratio := &pbm.ResRatios{}

	if value, exist := ratioMap[body.Symbol]; exist {
		ratio = &pbm.ResRatios{
			S:   value.S,
			Cp:  value.Cp,
			C:   value.C,
			Cr:  value.Cr,
			Eps: value.Eps,
			Per: value.Per,
			Bps: value.Bps,
			Pbr: value.Pbr,
			D:   value.D,
			Dr:  value.Dr,
		}
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, ratio)
}

func (r *Ratio) AddRatio(c *gin.Context) {
	var body pbm.ReqResRatios

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := r.DB.Day()
	rid := r.DB.NextSequence("ratios", "rid")
	snapshotLogs := &pbm.Ratios{
		Rid:      rid,
		Ratio:    body.Ratio,
		Region:   body.Region,
		CreateAt: at,
	}
	r.DB.Conn.Create(&snapshotLogs)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}
