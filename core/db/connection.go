package db

import (
	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"stockcontent-monitor-demo-back/core/config"
	"stockcontent-monitor-demo-back/ent"
)

func New(cfg config.Config) *ent.Client {
	drv, err := sql.Open("mysql", cfg.DBConnection())
	if err != nil {
		panic(err)
	}

	db := drv.DB()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(15)
	db.SetMaxOpenConns(15)

	options := []ent.Option{
		ent.Driver(drv),
	}

	if cfg.IsDebug() {
		options = append(options, ent.Debug())
	}

	return ent.NewClient(options...)
}
