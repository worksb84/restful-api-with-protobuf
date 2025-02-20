package c_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pbm"
	"restful/db"
	"time"

	"github.com/gin-gonic/gin"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/valkey-io/valkey-go"
)

type SEC struct {
	DB     *db.DB
	ApiKey string
	Valkey valkey.Client
}

func (s *SEC) TickerDetail(c *gin.Context) {
	var body pbm.ReqBySymbol
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	p := polygon.New(s.ApiKey)

	params := models.GetTickerDetailsParams{
		Ticker: body.Symbol,
	}.WithDate(models.Date(time.Now()))

	res, err := p.GetTickerDetails(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}

	resSECTickerDetail := &pbm.ResSECTickerDetail{}
	resSECTickerDetail.Address = res.Results.Address.Address1 + ", " + res.Results.Address.City + ", " + res.Results.Address.State + ", " + res.Results.Address.PostalCode
	resSECTickerDetail.Branding = res.Results.Branding.IconURL + "?apikey=" + s.ApiKey
	resSECTickerDetail.Cik = res.Results.CIK
	resSECTickerDetail.CompositeFigi = res.Results.CompositeFIGI
	resSECTickerDetail.CurrencyName = res.Results.CurrencyName
	resSECTickerDetail.Description = res.Results.Description
	resSECTickerDetail.HomepageUrl = res.Results.HomepageURL
	resSECTickerDetail.ListDate = time.Time(res.Results.ListDate).Format("2006-01-02")
	resSECTickerDetail.Locale = res.Results.Locale
	resSECTickerDetail.Market = res.Results.Market
	resSECTickerDetail.MarketCap = res.Results.MarketCap
	resSECTickerDetail.Name = res.Results.Name
	resSECTickerDetail.PhoneNumber = res.Results.PhoneNumber
	resSECTickerDetail.PrimaryExchange = res.Results.PrimaryExchange
	resSECTickerDetail.ShareClassFigi = res.Results.ShareClassFIGI
	resSECTickerDetail.ShareClassSharesOutstanding = float64(res.Results.ShareClassSharesOutstanding)
	resSECTickerDetail.SicCode = res.Results.SICCode
	resSECTickerDetail.SicDescription = res.Results.SICDescription
	resSECTickerDetail.Ticker = res.Results.Ticker
	resSECTickerDetail.TickerRoot = res.Results.TickerRoot
	resSECTickerDetail.TotalEmployees = float64(res.Results.TotalEmployees)
	resSECTickerDetail.WeightedSharesOutstanding = float64(res.Results.WeightedSharesOutstanding)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resSECTickerDetail)
}

func (s *SEC) RelatedCompanies(c *gin.Context) {
	var body pbm.ReqBySymbol
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	p := polygon.New(s.ApiKey)
	params := models.GetTickerRelatedCompaniesParams{
		Ticker: body.Symbol,
	}

	res, err := p.GetTickerRelatedCompanies(context.Background(), &params)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	msg, err := s.Valkey.Do(ctx, s.Valkey.B().Get().Key("STOCK:SEC:LIST_").Build()).ToString()
	if err != nil {
		log.Println(err)
	}
	var snapshotMap map[string]*pbm.Snapshot
	json.Unmarshal([]byte(msg), &snapshotMap)

	resRelatedCompanies := &pbm.ResRelatedCompanies{}

	for _, v := range res.Results {
		snapshot := snapshotMap[v.Ticker]

		relatedCompanies := &pbm.RelatedCompanies{
			Ticker: v.Ticker,
			Name:   snapshot.N,
		}

		resRelatedCompanies.RelatedCompanies = append(resRelatedCompanies.RelatedCompanies, relatedCompanies)

	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resRelatedCompanies)
}
