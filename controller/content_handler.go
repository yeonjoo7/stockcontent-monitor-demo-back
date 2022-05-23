package controller

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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

func ContentDetail(c echo.Context) error {
	contentId := c.Param("id")
	// db.First()

	return c.String(http.StatusOK, contentId)
}

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
	ContentId    uuid.UUID  `gorm:"type:varchar(36);primaryKey;not null;" json:"contentId"`
	State        Videostate `gorm:"type:varchar(30);not null;default:NONE;"  json:"state" validate:"eq=APPORVE|eq=DENY|eq=NONE"`
	MonitorExp   int64      `gorm:"autoUpdateTime:milli;" json:"monitorExp"`
	Subject      string     `gorm:"type:varchar(60);not null" json:"subject"`
	Description  string     `gorm:"type:varchar(300);not null" json:"description"`
	ThumbnailImg string     `gorm:"not null" json:"thumbnailImg"`
	VideoUrl     string     `gorm:"not null" json:"videoUrl"`
	Tags         Tags       `gorm:"type:json" json:"tags"`
	UploadedAt   time.Time  `gorm:"type:datetime(6);not null;"`
}
