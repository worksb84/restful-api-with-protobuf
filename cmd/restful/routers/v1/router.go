package r_v1

import (
	c_v1 "restful/controllers/v1"
	"restful/db"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
)

func UsersRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Users{DB: db}

	v1 := router.Group("/v1/users")
	{
		v1.POST("/login", c.Login)
		v1.POST("/profile", c.Profile)
	}
}

func DeckRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Deck{DB: db}

	v1 := router.Group("/v1/deck")
	{
		v1.POST("/add", c.AddDeck)
		v1.POST("/delete", c.DeleteDeck)
		v1.POST("/get", c.GetDeck)
		v1.POST("/list", c.GetDeckList)
		v1.POST("/update", c.UpdateDeck)
	}
}

func FavoriteRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Favorite{DB: db}

	v1 := router.Group("/v1/favorite")
	{
		v1.POST("/add", c.AddFavorite)
		v1.POST("/delete", c.DeleteFavorite)
		v1.POST("/list", c.GetFavoriteList)
	}
}

func ServersRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Server{DB: db}

	v1 := router.Group("/v1/server")
	{
		v1.POST("/list", c.GetList)
	}
}

func BalanaceRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Balance{DB: db}

	v1 := router.Group("/v1/balance")
	{
		v1.POST("/update", c.UpdateBalance)
	}
}

func PurchaseRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Purchase{DB: db}

	v1 := router.Group("/v1/purchase")
	{
		v1.POST("", c.Purchase)
	}
}

func SnapshotRouter(router *gin.Engine, db *db.DB, valkey valkey.Client) {
	c := c_v1.SnapshotLogs{DB: db, Valkey: valkey}

	v1 := router.Group("/v1/snapshotLogs")
	{
		v1.POST("", c.GetSnapshot)
		v1.POST("/add", c.AddSnapshot)
		v1.POST("/all", c.AllSnapshot)
		v1.POST("/random/:no", c.RandomSnapshot)
		v1.POST("/get", c.GetSymbol)
		v1.POST("/snapshot", c.GetSnapshots)
	}
}

func RatioRouter(router *gin.Engine, db *db.DB) {
	c := c_v1.Ratio{DB: db}

	v1 := router.Group("/v1/ratios")
	{
		v1.POST("/add", c.AddRatio)
		v1.POST("/get", c.GetSymbol)
	}
}

func SECRouter(router *gin.Engine, db *db.DB, valkey valkey.Client) {
	c := c_v1.SEC{DB: db, ApiKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Valkey: valkey}

	v1 := router.Group("/v1/sec")
	{
		v1.POST("tickerDetail", c.TickerDetail)
		v1.POST("relatedCompanies", c.RelatedCompanies)
	}
}
