package model

import (
	"database/sql/driver"
	"os"
	"stockcontent-monitor-demo-back/util/sqlx"

	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MysqlRepo() *gorm.DB {
	db := getDB()
	migrate(db)
	return db
}

func getDB() *gorm.DB {
	conn := os.Getenv("DB_CONN")

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(15)
	return db
}

// example

type HelloEntity struct {
	Id   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"size:300;not null" json:"name"`
}

func (HelloEntity) TableName() string {
	return "hello"
}

// video entity

type VideoEntity struct {
	ContentId     uuid.UUID  `gorm:"type:varchar(36);primaryKey;not null;" json:"contentId"`
	StateLabel    Videostate `gorm:"type:varchar(30);not null;default:NONE;"  json:"stateLabel" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
	MonitorExp    int64      `gorm:"autoUpdateTime:milli;" json:"monitorExp"`
	Subject       string     `gorm:"type:varchar(60);not null" json:"subject"`
	Description   string     `gorm:"type:varchar(300);not null" json:"description"`
	Thumb         string     `gorm:"not null" json:"thumb"`
	SampleContent string     `gorm:"not null" json:"sampleContent"`
	Tags          []string   `gorm:"type:json" json:"tags"`
	// Tags       datatypes.JSON `gorm:"type:json" json:"tags"`
	UploadedAt time.Time `gorm:"type:datetime(6);not null;" json:"uploadedAt"`

	DenyLogs []DenyLogEntity `gorm:"foreignKey:ContentId" json:"denyLog"`
}

func (t DenyTagEntity) Value() (driver.Value, error) {
	return sqlx.JsonValue(t)
}

func (t *DenyTagEntity) Scan(src interface{}) error {
	return sqlx.JsonScan(t, src)
}

type Videostate string

const (
	NONE    Videostate = "NONE"
	DENY    Videostate = "DENY"
	APPROVE Videostate = "APPROVE"
)

func (VideoEntity) TableName() string {
	return "video"
}

// deny log

type DenyLogEntity struct {
	LogId     int64     `gorm:"primaryKey;auto_increment"`
	ContentId uuid.UUID `gorm:"type:varchar(36);not null" json:"contentId"`
	// ContentID     []VideoEntity   `gorm:"foreignKey:ContentId" json:"contentID"`
	DenyTagEntity []DenyTagEntity `gorm:"many2many:stock_content_deny_tag"`
	Reason        string          `gorm:"type:varchar(500);" json:"reason"`
	DeniedAt      time.Time       `gorm:"type:datetime(6);not null"`
}

func (DenyLogEntity) TableName() string {
	return "deny_log"
}

// deny Tag

type DenyTagEntity struct {
	TagId   int16  `gorm:"primaryKey;auto_increment"`
	Content string `gorm:"type:varchar(100);not null" json:"content"`
}

func (DenyTagEntity) TableName() string {
	return "deny_tag"
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&DenyTagEntity{}, &VideoEntity{}, &DenyLogEntity{})
	// db.Exec("ALTER TABLE deny_log ADD FOREIGN KEY (content_id) REFERENCES video(content_id);")
	if err != nil {
		panic(err)
	}
}

func isValid(state VideoEntity) bool {

	switch state.StateLabel {
	case DENY, APPROVE, NONE:
		return true
	default:
		return false
	}
}

type inputBody struct {
	StateLabel Videostate `json:"stateLabel" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
}
