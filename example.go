package main

import (
	"database/sql/driver"
	"stockcontent-monitor-demo-back/util/sqlx"
)

type Tags []string

type ExampleJsonEntity struct {
	Tags Tags `gorm:"type:json"`
}

func (t Tags) Value() (driver.Value, error) {
	return sqlx.JsonValue(t)
}

func (t *Tags) Scan(src interface{}) error {
	return sqlx.JsonScan(t, src)
}
