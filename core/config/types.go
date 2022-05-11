package config

import (
	"fmt"
	"net/url"
	"stockcontent-monitor-demo-back/build"
	"time"
)

type Config interface {
	IsDebug() bool
	DBConnection() string
	ServerAddress() string
	UseCaseTimeout() time.Duration
}

var _ Config = &configImpl{}

type configImpl struct {
	DB struct {
		User   string     `json:"user"`
		Pass   string     `json:"pass"`
		Host   string     `json:"host"`
		Port   uint16     `json:"port"`
		Name   string     `json:"name"`
		Values url.Values `json:"query_values"`
	} `json:"db"`

	ServeAddr string `json:"serve_addr"`

	UseCaseTimeoutStr string `json:"use_case_timeout"`
}

func (c *configImpl) UseCaseTimeout() time.Duration {
	d, _ := time.ParseDuration(c.UseCaseTimeoutStr)
	return d
}

func (c *configImpl) ServerAddress() string {
	return c.ServeAddr
}

func (c *configImpl) DBConnection() string {
	db := c.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		db.User, db.Pass, db.Host, db.Port, db.Name, db.Values.Encode())
}

func (c *configImpl) IsDebug() bool {
	return build.RuntimeProfile == build.RuntimeLocal
}
