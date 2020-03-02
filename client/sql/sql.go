package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/hokonco/gdk/errors"
	"github.com/hokonco/gdk/pool"
	"github.com/hokonco/gdk/timeout"
)

// Instance ...
type Instance interface {
	Close() error
	Conn(context.Context) (*sql.Conn, error)
	Do(context.Context, func() error) error
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

func (x *impl) Conn(ctx context.Context) (*sql.Conn, error) {
	return x.DB.Conn(ctx)
}

func (x *impl) Do(ctx context.Context, fn func() error) error {
	if fn == nil {
		return errors.New("sql: empty callback func")
	}
	x.Init()
	defer x.Done()

	return timeout.WithContext(ctx,
		func() error {
			<-time.After(time.Minute)
			return context.DeadlineExceeded
		},
		func() error {
			return fn()
		},
	)
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
