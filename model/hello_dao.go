package model

import (
	"github.com/google/uuid"
	"stockcontent-monitor-demo-back/db"
)

func init() {
	db.Client.AutoMigrate(&HelloEntity{})
}

type HelloEntity struct {
	Id   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"size:300;not null" json:"name"`
}

func (HelloEntity) TableName() string {
	return "tb_hello"
}

func CreateHello(name string) (res *HelloEntity, err error) {
	res = &HelloEntity{
		Id:   uuid.New(),
		Name: name,
	}
	err = db.Client.Create(res).Error
	return
}

func FetchHello() (res []HelloEntity, err error) {
	err = db.Client.Find(&res).Error
	return
}

func UpdateHello(id uuid.UUID, name string) (err error) {
	err = db.Client.Model(&HelloEntity{Id: id}).Update("name", name).Error
	return
}

func DeleteHello(id uuid.UUID) (err error) {
	err = db.Client.Delete(&HelloEntity{Id: id}).Error
	return
}
