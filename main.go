package main

import (
	"context"
	"fmt"
	"log"
	"mybankcli/api"
	"mybankcli/api/handlers"
	"mybankcli/pkg/account"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/manager/service"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"go.uber.org/dig"
)

func main() {
	handlers.LogInit()
	fmt.Println("Start server....")
	host := "0.0.0.0"
	port := "7778"
	dsn := "postgres://app:pass@localhost:5432/db"
	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func execute(host, port, dsn string) (err error) {
	deps := []interface{}{
		api.NewServer,
		mux.NewRouter,
		func() (*pgx.Conn, error) {
			connCtx, err := context.WithTimeout(context.Background(), time.Second*5)
			if err != nil {
				log.Print(err)
			}
			return pgx.Connect(connCtx, dsn)
		},
		service.NewManagerRepository,
		account.NewAccountRepository,
		customers.NewCustomerRepository,
		handlers.NewCustomerHandler,
		handlers.NewManagerHandler,
		handlers.NewAccountHandler,
		func(server *api.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}
	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			return err
		}
	}
	err = container.Invoke(func(server *api.Server) {
		server.Init()
	})
	if err != nil {
		return err
	}
	return container.Invoke(func(s *http.Server) error {
		return s.ListenAndServe()
	})
}
