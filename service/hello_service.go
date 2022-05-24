package service

import (
	"github.com/google/uuid"
	"stockcontent-monitor-demo-back/model"
)

func CreateHello(name string) (*model.HelloEntity, error) {
	return model.CreateHello(name)
}

func FetchHello() ([]model.HelloEntity, error) {
	return model.FetchHello()
}

func UpdateHello(id uuid.UUID, name string) error {
	return model.UpdateHello(id, name)
}

func DeleteHello(id uuid.UUID) error {
	return model.DeleteHello(id)
}
