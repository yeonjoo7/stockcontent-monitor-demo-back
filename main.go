package main

import (
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

	// e.POST("/content/:id/deny", func(c echo.Context) error {
	// 	var binder struct {
	// 		Content_id uuid.UUID `json:"-" param:"id"`
	// 		Reason     string    `json:"Description"`
	// 		Tag        []string  `json:"tag_id"`
	// 	}

	// 	err := c.Bind(&binder)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// 	}

	// 	// var deny_log DenylogEntity

	// 	err = db.Transaction(func(tx *gorm.DB) (err error) {
	// 		// value update
	// 		// business logic

	// 		err = tx.M(&denyLog, binder.Content_id).Error
	// 		if err != nil {
	// 			return
	// 		}

	// 		hello.Name = binder.Name
	// 		err = tx.Save(&hello).Error
	// 		return
	// 	}, &sql.TxOptions{
	// 		Isolation: sql.LevelSerializable,
	// 	})

	// 	if err == gorm.ErrRecordNotFound {
	// 		return echo.NewHTTPError(http.StatusNotFound, "hello entity not found")
	// 	} else if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.NoContent(http.StatusNoContent)
	// })

	// ===========================================================================
	// ===========================================================================
	// ===========================================================================

	e.POST("/hello", func(c echo.Context) error {
		var hello HelloEntity
		err := c.Bind(&hello)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		hello.Id = uuid.New()
		err = db.Create(&hello).Error

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"id": hello.Id,
		})
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

				// if video.StateLabel {

				// }

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

	// e.GET("/hello", func(c echo.Context) error {
	// 	var items []HelloEntity
	// 	err := db.Find(&items).Error
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	if len(items) == 0 {
	// 		return c.NoContent(http.StatusNoContent)
	// 	}

	// 	return c.JSON(http.StatusOK, items)
	// })

	// e.POST("/hello", func(c echo.Context) error {
	// 	var hello HelloEntity
	// 	err := c.Bind(&hello)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// 	}

	// 	hello.Id = uuid.New()
	// 	err = db.Create(&hello).Error
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.JSON(http.StatusCreated, echo.Map{
	// 		"id": hello.Id,
	// 	})
	// })

	// e.GET("/hello/:id", func(c echo.Context) error {
	// 	var binder struct {
	// 		HelloId uuid.UUID `param:"id"`
	// 	}

	// 	err := c.Bind(&binder)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// 	}

	// 	var hello HelloEntity
	// 	err = db.First(&hello, binder.HelloId).Error
	// 	if err == gorm.ErrRecordNotFound {
	// 		return echo.NewHTTPError(http.StatusNotFound, "hello entity not found")
	// 	} else if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.JSON(http.StatusOK, hello)
	// })

	// e.PUT("/hello/:id", func(c echo.Context) error {
	// 	var binder struct {
	// 		HelloId uuid.UUID `json:"-" param:"id"`
	// 		Name    string    `json:"name"`
	// 	}

	// 	err := c.Bind(&binder)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// 	}

	// 	var hello HelloEntity

	// 	err = db.Transaction(func(tx *gorm.DB) (err error) {
	// 		// value update
	// 		// business logic

	// 		err = tx.Clauses(clause.Locking{
	// 			Strength: "UPDATE",
	// 		}).First(&hello, binder.HelloId).Error
	// 		if err != nil {
	// 			return
	// 		}
	// 		hello.Name = binder.Name
	// 		err = tx.Save(&hello).Error
	// 		return
	// 	}, &sql.TxOptions{
	// 		Isolation: sql.LevelSerializable,
	// 	})

	// 	if err == gorm.ErrRecordNotFound {
	// 		return echo.NewHTTPError(http.StatusNotFound, "hello entity not found")
	// 	} else if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.NoContent(http.StatusNoContent)
	// })

	// e.DELETE("/hello/:id", func(c echo.Context) error {
	// 	var binder struct {
	// 		HelloId uuid.UUID `json:"-" param:"id"`
	// 	}
	// 	err := c.Bind(&binder)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// 	}

	// 	err = db.Delete(&HelloEntity{}, binder.HelloId).Error
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	return c.NoContent(http.StatusNoContent)
	// })
	e.Start(os.Getenv("SERVE_ADDR"))
}

// func getDB() *gorm.DB {
// 	conn := os.Getenv("DB_CONN")
// 	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = sqlDB.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	sqlDB.SetMaxIdleConns(15)
// 	sqlDB.SetMaxOpenConns(15)
// 	return db
// }

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

// func migrate(db *gorm.DB) {
// 	err := db.AutoMigrate(&HelloEntity{}, &DenyLogEntity{}, &VideoEntity{}, &DenyTagEntity{})
// 	// db.Exec("ALTER TABLE deny_log ADD FOREIGN KEY (content_id) REFERENCES video(content_id);")
// 	if err != nil {
// 		panic(err)
// 	}
// }

// const (
// 	F_DENY string = "DENY"
// )

// func isValid(state VideoEntity) bool {

// 	switch state.State {
// 	case DENY, APPROVE, NONE:
// 		return true
// 	default:
// 		return false
// 	}
// }

// type inputBody struct {
// 	State Videostate `json:"state" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
// }
