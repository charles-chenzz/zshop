package data

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"zshop/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewMysql, NewGreeterRepo, NewOrderRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *sqlx.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *sqlx.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewMysql(c *conf.Data) (*sqlx.DB, error) {
	// abstract into a package todo
	db, err := sqlx.Connect("mysql", c.Database.Source)
	if err != nil {
		fmt.Printf("error:%v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error:%v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return db, nil
}
