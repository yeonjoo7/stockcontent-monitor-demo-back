package controller

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func ContentDetail(c echo.Context, db *gorm.DB) error {
	var binder struct {
		ContentId uuid.UUID `param:"id"`
	}

	err := c.Bind(&binder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var Result struct {
		Video    VideoEntity     `json:"video"`
		DenyLogs []DenyLogEntity `json:"denyLogs"` // omitempty 속성 : 값이 없을 때는 무시합니다.
	}

	err = db.First(&Result.Video, binder.ContentId).Error
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "video entity not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if Result.Video.StateLabel == "NONE" {
		if Result.Video.MonitorExp-time.Now().Unix() > 0 {
			Result.Video.StateLabel = "CHECK"
		} else {
			Result.Video.StateLabel = "WAIT"
		}
	}

	err = nil
	err = db.Find(&Result.DenyLogs, binder.ContentId).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Result)
}

func ContentList(c echo.Context) error {

	// State := c.QueryParam("state")
	// ContentId := c.QueryParam("contentid")
	// Start := c.QueryParam("start")
	// Limit := c.QueryParam("limit")

	var video VideoEntity
	err := c.Bind(&video)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// service.VideoService(&video)

	return c.JSON(http.StatusCreated, echo.Map{
		"contentId": video.ContentId,
	})
}

func CheckMonotoring(c echo.Context) error {
	contentId := c.Param("id")

	return c.String(http.StatusOK, contentId)
}

// func ContentDetail(c echo.Context) error {
// 	contentId := c.Param("id")
// 	// db.First()

// 	return c.String(http.StatusOK, contentId)
// }

func MarkMonotoring(c echo.Context) error {
	contentId := c.Param("id")

	return c.String(http.StatusOK, contentId)
}

// video entity

type Videostate string

const (
	NONE    Videostate = "NONE"
	DENY    Videostate = "DENY"
	APPROVE Videostate = "APPROVE"
)

type Tags []string

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
