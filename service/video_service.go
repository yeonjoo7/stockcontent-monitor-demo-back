package service

import (
	"github.com/google/uuid"
	// "github.com/labstack/echo/v4"
	// "gorm.io/gorm"
)

func VideoService(video *VideoEntity) {
	video.Id = uuid.New()

	// db := main.db
	// err = db.Create(&video).Error
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }
}

type VideoEntity struct {
	Id   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"size:300;not null" json:"name"`
}
