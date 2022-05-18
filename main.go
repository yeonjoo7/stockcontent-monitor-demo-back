package main

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	db := getDB()
	migrate(db)

	e := echo.New()
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

	e.GET("/hello", func(c echo.Context) error {
		var items []HelloEntity
		err := db.Find(&items).Error
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if len(items) == 0 {
			return c.NoContent(http.StatusNoContent)
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": items,
		})
	})

	e.GET("/hello/:id", func(c echo.Context) error {
		var binder struct {
			HelloId uuid.UUID `param:"id"`
		}

		err := c.Bind(&binder)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var hello HelloEntity
		err = db.First(&hello, binder.HelloId).Error
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "hello entity not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, hello)
	})

	e.PUT("/hello/:id", func(c echo.Context) error {
		var binder struct {
			HelloId uuid.UUID `json:"-" param:"id"`
			Name    string    `json:"name"`
		}

		err := c.Bind(&binder)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var hello HelloEntity

		err = db.Transaction(func(tx *gorm.DB) (err error) {
			// value update
			// business logic

			err = tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).First(&hello, binder.HelloId).Error
			if err != nil {
				return
			}
			hello.Name = binder.Name
			err = tx.Save(&hello).Error
			return
		}, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "hello entity not found")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/hello/:id", func(c echo.Context) error {
		var binder struct {
			HelloId uuid.UUID `json:"-" param:"id"`
		}
		err := c.Bind(&binder)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = db.Delete(&HelloEntity{}, binder.HelloId).Error
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	})
	e.Start(os.Getenv("SERVE_ADDR"))
}

func getDB() *gorm.DB {
	conn := os.Getenv("DB_CONN")
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(15)
	return db
}

type HelloEntity struct {
	Id   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"size:300;not null" json:"name"`
}

func (HelloEntity) TableName() string {
	return "tb_hello"
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&HelloEntity{})
	if err != nil {
		panic(err)
	}
}
