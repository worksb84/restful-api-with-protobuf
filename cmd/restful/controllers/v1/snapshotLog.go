package c_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pbm"
	"restful/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
)

type SnapshotLogs struct {
	DB     *db.DB
	Valkey valkey.Client
}

func (s *SnapshotLogs) GetSnapshots(c *gin.Context) {
	var body pbm.ReqBySymbol
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	msg, err := s.Valkey.Do(ctx, s.Valkey.B().Get().Key(fmt.Sprintf("STOCK:%s:%s_", body.Region, body.Symbol)).Build()).ToString()
	if err != nil {
		log.Println(err)
	}
	var snapshotMap map[string]*pbm.Snapshot
	resSnapshots := &pbm.ResSnapshots{}

	json.Unmarshal([]byte(msg), &snapshotMap)

	for _, v := range snapshotMap {
		resSnapshots.Snapshots = append(resSnapshots.Snapshots, &pbm.ResSnapshot{
			N:   v.N,
			Ne:  v.Ne,
			S:   v.S,
			E:   v.E,
			C:   v.C,
			Cr:  v.Cr,
			Pcp: v.Pcp,
			Cp:  v.Cp,
			Op:  v.Op,
			Hp:  v.Hp,
			Lp:  v.Lp,
		})
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resSnapshots)
}

func (s *SnapshotLogs) GetSymbol(c *gin.Context) {
	var body pbm.ReqBySymbol
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	msg, err := s.Valkey.Do(ctx, s.Valkey.B().Get().Key(fmt.Sprintf("STOCK:%s:%s_", body.Region, body.Symbol)).Build()).ToString()
	if err != nil {
		log.Println(err)
	}
	var resSnapshot *pbm.ResSnapshot

	json.Unmarshal([]byte(msg), &resSnapshot)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resSnapshot)
}

func (s *SnapshotLogs) AllSnapshot(c *gin.Context) {
	var body pbm.ReqResSnapshotLogs

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	resSnapshots := &pbm.ResSnapshots{}
	if body.Region == "ALL" {
		for _, region := range []string{"KRX", "SEC"} {
			query := `select * from snapshot_logs where region = ? order by slid desc limit 1`

			var reqResSnapshotLogs *pbm.ReqResSnapshotLogs

			s.DB.Conn.Raw(query, region).Scan(&reqResSnapshotLogs)

			if reqResSnapshotLogs != nil {
				var snapshotMap map[string]*pbm.Snapshot
				json.Unmarshal([]byte(reqResSnapshotLogs.Snapshot), &snapshotMap)

				for _, snapshot := range snapshotMap {
					resSnapshots.Snapshots = append(resSnapshots.Snapshots, (*pbm.ResSnapshot)(snapshot))
				}
			}
		}
	} else {
		query := `select * from snapshot_logs where region = ? order by slid desc limit 1`

		var reqResSnapshotLogs *pbm.ReqResSnapshotLogs

		s.DB.Conn.Raw(query, body.Region).Scan(&reqResSnapshotLogs)

		if reqResSnapshotLogs != nil {
			var snapshotMap map[string]*pbm.Snapshot
			json.Unmarshal([]byte(reqResSnapshotLogs.Snapshot), &snapshotMap)

			for _, snapshot := range snapshotMap {
				resSnapshots.Snapshots = append(resSnapshots.Snapshots, (*pbm.ResSnapshot)(snapshot))
			}
		}
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resSnapshots)
}

func (s *SnapshotLogs) RandomSnapshot(c *gin.Context) {
	var body pbm.ReqResSnapshotLogs
	no, _ := strconv.Atoi(c.Param("no"))

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from snapshot_logs where region = ? order by slid desc limit 1`

	var reqResSnapshotLogs *pbm.ReqResSnapshotLogs

	s.DB.Conn.Raw(query, body.Region).Scan(&reqResSnapshotLogs)

	var snapshotMap map[string]*pbm.Snapshot
	json.Unmarshal([]byte(reqResSnapshotLogs.Snapshot), &snapshotMap)
	resSnapshots := &pbm.ResSnapshots{}

	for _, snapshot := range snapshotMap {
		resSnapshots.Snapshots = append(resSnapshots.Snapshots, (*pbm.ResSnapshot)(snapshot))
	}

	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(resSnapshots.Snapshots), func(i, j int) {
		resSnapshots.Snapshots[i], resSnapshots.Snapshots[j] = resSnapshots.Snapshots[j], resSnapshots.Snapshots[i]
	})

	resSnapshots.Snapshots = resSnapshots.Snapshots[:no]
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resSnapshots)
}

func (s *SnapshotLogs) GetSnapshot(c *gin.Context) {
	var body pbm.ReqResSnapshotLogs

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	query := `select * from snapshot_logs where region = ? order by slid desc limit 1`

	var reqResSnapshotLogs *pbm.ReqResSnapshotLogs

	s.DB.Conn.Raw(query, body.Region).Scan(&reqResSnapshotLogs)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, reqResSnapshotLogs)
}

func (s *SnapshotLogs) AddSnapshot(c *gin.Context) {
	var body pbm.ReqResSnapshotLogs

	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
	}

	at := s.DB.Day()
	slid := s.DB.NextSequence("snapshot_logs", "slid")
	snapshotLogs := &pbm.SnapshotLogs{
		Slid:     slid,
		Snapshot: body.Snapshot,
		Region:   body.Region,
		CreateAt: at,
	}
	s.DB.Conn.Create(&snapshotLogs)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, nil)
}
