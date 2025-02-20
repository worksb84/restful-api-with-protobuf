package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Conn *gorm.DB
}

func NewDB(config *Config) *DB {
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		config.Host,
		config.Username,
		config.Password,
		config.DBname,
		config.Port,
	)
	conn, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Logger: logger.Discard,
	})

	if err != nil {
		log.Fatalln(err)
	}

	sqlDB, _ := conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)

	return &DB{
		Conn: conn,
	}
}

func (db *DB) CreateTable(model interface{}) {
	db.Conn.Migrator().CreateTable(model)
}

func (db *DB) AutoIncrement(tableName, columnName string) {
	db.Conn.Exec(fmt.Sprintf("CREATE SEQUENCE %s_%s_seq", tableName, columnName))
	db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT nextval('%s_%s_seq')", tableName, columnName, tableName, columnName))
}

func (db *DB) Uniquekey(tableName, columnName string) {
	db.Conn.Exec(fmt.Sprintf("CREATE UNIQUE INDEX CONCURRENTLY %s_%s ON %s (%s)", tableName, columnName, tableName, columnName))
	db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT unique_%s UNIQUE USING INDEX %s_%s", tableName, columnName, tableName, columnName))
}

func (db *DB) CreateIndex(tableName, columnName string) {
	db.Conn.Exec(fmt.Sprintf("CREATE INDEX %s_%s_idx ON %s (%s)", tableName, columnName, tableName, columnName))
}

func (db *DB) AlterColumnType(tableName, columnName, columnType string, defaultValue interface{}) {
	db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, columnName, columnType))
	if defaultValue != nil {
		defaultString := defaultValue.(string)
		db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s", tableName, columnName, defaultString))
	}
}

func (db *DB) AlterColumnPrimaryKey(tableName string, columnName []string) {
	db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s_pk", tableName, tableName))
	db.Conn.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s_pk PRIMARY KEY (%s)", tableName, tableName, strings.Join(columnName, ",")))
}

func (db *DB) NextSequence(tableName, columnName string) int32 {
	var seq int32
	db.Conn.Raw(fmt.Sprintf("SELECT nextval('%s_%s_seq')", tableName, columnName)).First(&seq)
	return seq
}

func (db *DB) Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (db *DB) AddDays(at string, days int) string {
	return db.StringToTime(at).AddDate(0, 0, days).Format("2006-01-02 15:04:05")
}

func (db *DB) Day() string {
	return time.Now().Format("2006-01-02")
}

func (db *DB) StringToTime(t string) time.Time {
	dt, _ := time.Parse("2006-01-02 15:04:05", t)
	return dt
}

func (db *DB) TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
