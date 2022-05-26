package main

import (
	// "database/sql"
	"database/sql"
	"net/http"
	"os"
	"stockcontent-monitor-demo-back/controller"
	"stockcontent-monitor-demo-back/model"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	db := model.MysqlRepo()

	e := echo.New()

	e.GET("/deny-tag", func(c echo.Context) error {

		var denyTag []DenyTagEntity
		err := db.Find(&denyTag).Error

		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Content entity not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, denyTag)

	})
	e.POST("/content/:id/deny", func(c echo.Context) error {
		var ChangeDenyEntity struct {
			Content_id uuid.UUID `json:"-" param:"id"`
			Reason     string    `json:"reason"`
			Tag_id     []int     `json:"tag"`
		}

		err := c.Bind(&ChangeDenyEntity)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = db.Transaction(func(tx *gorm.DB) (err error) {

			var video VideoEntity
			var logId DenyLogEntity

			tx.Model(&video).Where("content_id = ?", ChangeDenyEntity.Content_id).Update("state_label", "DENY")
			tx.Exec("INSERT INTO deny_log(content_id, reason, denied_at) VALUES ( ? , ? , ?);", ChangeDenyEntity.Content_id, ChangeDenyEntity.Reason, time.Now())
			tx.Last(&logId)
			for i := 0; i < len(ChangeDenyEntity.Tag_id); i++ {
				tx.Exec("INSERT INTO stock_content_deny_tag VALUES ( ?, ? );", logId.LogId, ChangeDenyEntity.Tag_id[i])
			}

			return
		}, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "DENY entity not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, " message : Change DENY SUCCESS ")

	})

	e.POST("/content/:id/approve", func(c echo.Context) error {
		var ChangeApproceEntity struct {
			Content_id uuid.UUID `json:"-" param:"id"`
		}

		err := c.Bind(&ChangeApproceEntity)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var video VideoEntity
		db.Model(&video).Where("content_id = ?", ChangeApproceEntity.Content_id).Update("state_label", "APPROVE")

		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "APPROVE entity not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, " message : Change APPROVE SUCCESS ")

	})

	content := e.Group("/content")
	{
		// GET
		content.GET("/:id/monitoring", controller.CheckMonotoring)
		content.GET("/:id",
			func(c echo.Context) error {
				var binder struct {
					ContentId uuid.UUID `param:"id"`
				}

				err := c.Bind(&binder)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}

				var video VideoEntity
				err = db.First(&video, binder.ContentId).Error
				if err == gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusNotFound, "video entity not found")
				} else if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				return c.JSON(http.StatusOK, video)
			})
		content.GET("/", func(c echo.Context) error {
			var items []VideoEntity
			err := db.Find(&items).Error
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if len(items) == 0 {
				return c.NoContent(http.StatusNoContent)
			}

			return c.JSON(http.StatusOK, items)
		})
		// POST
		content.POST("/:id/monitoring", controller.MarkMonotoring)
	}

	e.Start(os.Getenv("SERVE_ADDR"))
}

// hello entity

type HelloEntity struct {
	Id   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"size:300;not null" json:"name"`
}

func (HelloEntity) TableName() string {
	return "hello"
}

// video entity

type VideoEntity struct {
	ContentId     uuid.UUID      `gorm:"type:varchar(36);primaryKey;not null;" json:"contentId"`
	StateLabel    Videostate     `gorm:"type:varchar(30);not null;default:NONE;"  json:"stateLabel" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
	MonitorExp    int64          `gorm:"autoUpdateTime:milli;" json:"monitorExp"`
	Subject       string         `gorm:"type:varchar(60);not null" json:"subject"`
	Description   string         `gorm:"type:varchar(300);not null" json:"description"`
	Thumb         string         `gorm:"not null" json:"thumb"`
	SampleContent string         `gorm:"not null" json:"sampleContent"`
	Tags          datatypes.JSON `gorm:"type:json" json:"tags"`
	UploadedAt    time.Time      `gorm:"type:datetime(6);not null;" json:"uploadedAt"`
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
	LogId int64 `gorm:"primaryKey;auto_increment"`
	// ContentId uuid.UUID `gorm:"type:varchar(36);not null" json:"contentId"`
	ContentId     []VideoEntity   `gorm:"foreignKey:ContentId" json:"contentId"`
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
