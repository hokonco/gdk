package sql

import (
	"context"
	"database/sql"

	"github.com/hokonco/gdk/errors"
	"github.com/hokonco/gdk/pool"
)

// Instance ...
type Instance interface {
	Close() error
	Do(context.Context, func(*sql.Conn) error) error
	Stats() sql.DBStats
}
type impl struct {
	*sql.DB
	pool.Instance
}

// Config ...
type Config struct {
	DriverName string
	DSN        string

	PoolSize int
}

// New ...
func New(conf Config) Instance {
	DB, err := sql.Open(conf.DriverName, conf.DSN)
	if err != nil {
		panic(err)
	}
	if DB == nil {
		panic(errors.New("sql: [%+v]", DB))
	}

	return &impl{DB, pool.New(conf.PoolSize)}
}

func (x *impl) Close() error {
	return x.DB.Close()
}

func (x *impl) Do(ctx context.Context, fn func(*sql.Conn) error) error {
	if fn == nil {
		return errors.New("")
	}
	conn, err := x.DB.Conn(ctx)
	if err != nil {
		return err
	}

	x.Init()
	defer x.Done()
	defer conn.Close()
	return fn(conn)
}

func (x *impl) Stats() sql.DBStats {
	return x.DB.Stats()
}

// func a() {
// 	var ctx = context.Background()
// 	New(Config{}).Do(ctx, func(c *sql.Conn) error {
// 		tx, err := c.BeginTx(ctx, &sql.TxOptions{})
// 		err = tx.Commit()
// 		return err
// 	})
// }
