package main

import (
	"context"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/hokonco/gdk/client"
	"github.com/hokonco/gdk/service"
	"github.com/hokonco/gdk/service/web"
	"github.com/hokonco/gdk/signal"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	v, err := client.Redis{
		Address: "localhost:6379",
	}.New().Pipeline(
		client.RedisArgument("SET", "KEY", "VAL", "EX", 60, "NX"),
		client.RedisArgument("GET", "KEY"),
	)
	fmt.Println(">>>", v, err)

	sql := client.SQL{
		DriverName: "postgres",
		DSN:        "host=0.0.0.0 port=5432 user=postgres password=password sslmode=disable",
	}.New()

	c, err := sql.Conn(ctx)
	if err != nil {
		panic(err)
	}
	err = sql.Do(ctx, func() error {
		// row := c.QueryRowContext(ctx, "")
		// err := row.Scan()
		// if err != nil {
		// 	return err
		// }
		rows, err := c.QueryContext(ctx, "")
		if err != nil {
			return err
		}
		for rows.Next() {
			err = rows.Err()
			if err != nil {
				return err
			}
			err = rows.Scan()
			if err != nil {
				return err
			}
		}
		fmt.Println(">>>", err)
		return err
	})
	fmt.Println(">>>", err)

	if true {
		panic(1)
	}
	service.New(service.Config{
		ListenSignals: signal.Wrap(
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTSTP,
		),
		Web: web.Config{
			Addr:           ":8080",
			OnShutdown:     func() { fmt.Println("web: OnShutdown") },
			GlobalPoolSize: 1,
			PanicHandler: func(rcv interface{}) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("PANIC!! [" + fmt.Sprint(rcv) + "]"))
				}
			},
			HandlerFuncs: web.HandlerFuncs{}.
				Add(0, http.MethodGet, "/", func(http.ResponseWriter, *http.Request) {
					<-time.After(5 * time.Second)
					msg := "yeye"
					fmt.Println(msg)
					panic(msg)
				}),
		},
	})
}
