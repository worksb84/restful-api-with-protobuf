package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pbm"
	"restful/db"
	r_v1 "restful/routers/v1"
	"restful/setup"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
)

func main() {
	time.LoadLocation("Asia/Seoul")

	dbConfig := &db.Config{
		Host:     "xxxxxxxxxxxxxxxx.ap-northeast-2.rds.amazonaws.com",
		Port:     "xxxx",
		Username: "xxxx",
		Password: "xxxx",
		DBname:   "xxxx",
	}

	dbConn := db.NewDB(dbConfig)
	setup.BalanceMasters(dbConn, &pbm.BalanceMasters{})
	setup.BalanceDetails(dbConn, &pbm.BalanceDetails{})
	setup.DeckMasters(dbConn, &pbm.DeckMasters{})
	setup.DeckDetails(dbConn, &pbm.DeckDetails{})
	setup.Favorities(dbConn, &pbm.Favorities{})
	setup.GameRecordMasters(dbConn, &pbm.GameRecordMasters{})
	setup.GameRecordDetails(dbConn, &pbm.GameRecordDetails{})
	setup.PokerGameMasters(dbConn, &pbm.PokerGameMasters{})
	setup.PokerGameDetails(dbConn, &pbm.PokerGameDetails{})
	setup.Purchases(dbConn, &pbm.Purchases{})
	setup.Servers(dbConn, &pbm.Servers{})
	setup.Users(dbConn, &pbm.Users{})

	setup.SnapshotLogs(dbConn, &pbm.SnapshotLogs{})
	setup.Ratios(dbConn, &pbm.Ratios{})

	r := gin.Default()

	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			MaxAge:       12 * time.Hour,
		}))

	val, _ := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%s", "0.0.0.0", "0000")},
		TLSConfig:   &tls.Config{},
		SelectDB:    0,
		ClientName:  "-",
	})

	r_v1.ServersRouter(r, dbConn)
	r_v1.UsersRouter(r, dbConn)
	r_v1.BalanaceRouter(r, dbConn)
	r_v1.DeckRouter(r, dbConn)
	r_v1.FavoriteRouter(r, dbConn)
	r_v1.PurchaseRouter(r, dbConn)
	r_v1.SnapshotRouter(r, dbConn, val)
	r_v1.RatioRouter(r, dbConn)
	r_v1.SECRouter(r, dbConn, val)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", "8001"),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println(sig)
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
