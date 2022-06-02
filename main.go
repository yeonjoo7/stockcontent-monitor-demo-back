package main

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"database/sql"
	"net/http"
	"os"
	"stockcontent-monitor-demo-back/model"
	"stockcontent-monitor-demo-back/util/sqlx"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

const (
	MonitroEpxiresValue = 1
	MonitorExpiresTime  = MonitroEpxiresValue * time.Second
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
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "hello world",
		})
	})
	e.GET("/deny-tag", func(c echo.Context) error {

		var denyTag []DenyTagEntity
		err := db.Find(&denyTag).Error

		if err != nil {
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
			logId := DenyLogEntity{
				ContentId: ChangeDenyEntity.Content_id,
				Reason:    ChangeDenyEntity.Reason,
				DeniedAt:  time.Now(),
			}

			tx.Model(&video).Where("content_id = ?", ChangeDenyEntity.Content_id).Update("state_label", "DENY")
			tx.Create(&logId)

			for i := 0; i < len(ChangeDenyEntity.Tag_id); i++ {
				tx.Exec("INSERT INTO stock_content_deny_tag VALUES ( ?, ? );", logId.LogId, ChangeDenyEntity.Tag_id[i])
			}

			return
		}, &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
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

		// var video VideoEntity
		db.Model(&VideoEntity{}).Where("content_id = ?", ChangeApproceEntity.Content_id).Update("state_label", "APPROVE")

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
		content.GET("/:id/monitoring", func(c echo.Context) error {
			var binder struct {
				ContentId uuid.UUID `param:"id"`
			}

			err := c.Bind(&binder)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			var Video VideoEntity

			err = db.First(&Video, binder.ContentId).Error
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "No record")
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			// content-label 변환

			if Video.StateLabel == "NONE" {
				if Video.MonitorExp-time.Now().UnixMilli() > 0 {
					Video.StateLabel = "CHECK"
				} else {
					Video.StateLabel = "WAIT"
				}
			}

			// content의 stateLabel만 반환

			return c.JSON(http.StatusOK, echo.Map{"stateLabel": Video.StateLabel})
		})

		content.GET("/:id",
			func(c echo.Context) error {
				var binder struct {
					ContentId uuid.UUID `param:"id"`
				}

				err := c.Bind(&binder)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}

				var Video VideoEntity

				err = db.First(&Video, binder.ContentId).Error
				if err == gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusNotFound, "No record")
				} else if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				// content-label 변환

				if Video.StateLabel == "NONE" {
					if Video.MonitorExp-time.Now().UnixMilli() > 0 {
						Video.StateLabel = "CHECK"
					} else {
						Video.StateLabel = "WAIT"
					}
				}

				// content_id별로 deny-log 조회

				var DenyLogResult []struct {
					LogId         int64     `gorm:"column:log_id"`
					ContentId     uuid.UUID `gorm:"column:content_id"`
					DenyTagEntity string    `gorm:"column:deny_tag"` // "asadf,qwe"
					Reason        string    `gorm:"column:reason"`
					DeniedAt      time.Time `gorm:"column:denied_at"`
				}

				db.Raw(`SELECT dl.log_id, dl.content_id, dl.reason, dl.denied_at, 
				group_concat(dt.content) AS deny_tag 
				FROM deny_log dl 
				LEFT JOIN stock_content_deny_tag sc ON dl.log_id=sc.deny_log_entity_log_id 
				LEFT JOIN deny_tag dt	ON sc.deny_tag_entity_tag_id=dt.tag_id 
				WHERE dl.content_id = ? GROUP BY dl.log_id`, binder.ContentId).Scan(&DenyLogResult)

				// content의 키 값으로 deny-log를 추가해서 반환

				var Result struct {
					ContentId   uuid.UUID  `json:"contentId"`
					StateLabel  Videostate `json:"stateLabel"`
					MonitorExp  int64      `json:"monitorExp"`
					Subject     string     `json:"subject"`
					Description string     `json:"description"`
					// Thumb         string          `json:"thumb"`
					SampleContent string          `json:"sampleContent"`
					Tags          TagList         `gorm:"type:json" json:"tags"`
					UploadedAt    time.Time       `json:"uploadedAt"`
					DenyLogs      []denyLogResult `json:"denyLogs"`
				}
				Result.ContentId = Video.ContentId
				Result.StateLabel = Video.StateLabel
				Result.MonitorExp = Video.MonitorExp
				Result.Subject = Video.Subject
				Result.Description = Video.Description
				// Result.Thumb = Video.Thumb
				Result.SampleContent = Video.SampleContent
				Result.Tags = Video.Tags
				Result.UploadedAt = Video.UploadedAt

				Result.DenyLogs = make([]denyLogResult, len(DenyLogResult))
				for i, v := range DenyLogResult {
					Result.DenyLogs[i] = denyLogResult{
						LogId: v.LogId,
						// ContentId:     v.ContentId,
						DenyTagEntity: strings.Split(v.DenyTagEntity, ","), // ** todo refactor
						Reason:        v.Reason,
						DeniedAt:      v.DeniedAt,
					}
				}
				// Result.DenyLogs = DenyLogResult

				return c.JSON(http.StatusOK, Result)
			})

		content.GET("/", func(c echo.Context) error {

			// 필수 값인 state, lim이 들어왔는지 확인

			if c.QueryParam("state") == "" || c.QueryParam("lim") == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "missing queryparams : state or lim")
			}
			state := c.QueryParam("state")
			limit, err := strconv.Atoi(c.QueryParam("lim"))

			start, _ := strconv.Atoi(c.QueryParam("start"))
			// contentId := c.QueryParam("contentId")

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			var items []VideoEntity
			err = db.Where("state_label = ?", state).
				Offset(start).
				Limit(limit).
				Order("`uploaded_at` ASC").
				Find(&items).Error

			var total int64
			db.Model(&VideoEntity{}).Where("state_label = ?", state).Count(&total)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			// state 변환

			for i := 0; i < len(items); i++ {
				if items[i].StateLabel == "NONE" {
					if items[i].MonitorExp-time.Now().UnixMilli() > 0 {
						items[i].StateLabel = "CHECK"
					} else {
						items[i].StateLabel = "WAIT"
					}
				}
			}

			// 리스트 페이지 수를 추가해서 반환하기 -> 전체 아이템 갯수를 반환하는 방법!

			type Contents struct {
				Items      []VideoEntity `json:"items"`
				TotalItems int           `json:"totalItems"`
			}

			// totalPages := int(math.Ceil(float64(total) / float64(limit)))
			result := Contents{Items: items, TotalItems: int(total)}
			return c.JSON(http.StatusOK, result)
		})

		// PUT
		content.PUT("/:id/monitoring", func(c echo.Context) error {
			var binder struct {
				ContentId uuid.UUID `param:"id"`
			}

			err := c.Bind(&binder)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			var Video VideoEntity

			err = db.First(&Video, binder.ContentId).Error
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "No record")
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			Video.MonitorExp = time.Now().Add(MonitorExpiresTime).UnixMilli()
			db.Save(&Video)

			return c.JSON(http.StatusOK, echo.Map{
				"message": fmt.Sprintf("Time to monitor ends in %d minutes", MonitroEpxiresValue),
			})
		})
	}

	e.Start(os.Getenv("SERVE_ADDR"))
}

type denyLogResult struct {
	LogId         int64     `json:"logId"`
	ContentId     uuid.UUID `json:"contentId"`
	DenyTagEntity []string  `json:"deny_tag"`
	Reason        string    `json:"reason"`
	DeniedAt      time.Time `json:"denied_at"`
}

// video entity

type TagList []string

type VideoEntity struct {
	ContentId     uuid.UUID  `gorm:"type:varchar(36);primaryKey;not null;" json:"contentId"`
	StateLabel    Videostate `gorm:"type:varchar(30);not null;default:NONE;"  json:"stateLabel" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
	MonitorExp    int64      `gorm:"not null;" json:"monitorExp"`
	Subject       string     `gorm:"type:varchar(60);not null" json:"subject"`
	Description   string     `gorm:"type:varchar(300);not null" json:"description"`
	Thumb         string     `gorm:"not null" json:"thumb"`
	SampleContent string     `gorm:"not null" json:"sampleContent"`
	Tags          TagList    `gorm:"type:json" json:"tags"`
	// Tags       datatypes.JSON `gorm:"type:json" json:"tags"`
	UploadedAt time.Time `gorm:"type:datetime(6);not null;" json:"uploadedAt"`

	DenyLogs []DenyLogEntity `gorm:"foreignKey:ContentId" json:"denyLog"`
}

func (t TagList) Value() (driver.Value, error) {
	return sqlx.JsonValue(t)
}

func (t *TagList) Scan(src interface{}) error {
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
	LogId         int64           `gorm:"primaryKey;auto_increment" json:"logId"`
	ContentId     uuid.UUID       `gorm:"type:varchar(36);not null" json:"contentId"`
	DenyTagEntity []DenyTagEntity `gorm:"many2many:stock_content_deny_tag"`
	Reason        string          `gorm:"type:varchar(500);" json:"reason"`
	DeniedAt      time.Time       `gorm:"type:datetime(6);not null"`
}

func (DenyLogEntity) TableName() string {
	return "deny_log"
}

// deny Tag

type DenyTagEntity struct {
	TagId   int16  `gorm:"primaryKey;auto_increment" json:"tagId"`
	Content string `gorm:"type:varchar(100);not null" json:"content"`
}

func (DenyTagEntity) TableName() string {
	return "deny_tag"
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
