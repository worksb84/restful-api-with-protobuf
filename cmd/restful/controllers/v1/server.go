package c_v1

import (
	"log"
	"net/http"
	"pbm"
	"restful/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB *db.DB
}

func (s *Server) GetList(c *gin.Context) {
	var body pbm.ReqByAppID

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from servers s where s.app_id = ?`

	var resServers []*pbm.ResServers

	s.DB.Conn.Raw(query, body.AppId).Scan(&resServers)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resServers)
}
