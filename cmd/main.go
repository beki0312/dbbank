package main

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg//customer/services"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"

	"os"

	"github.com/jackc/pgx/v4"
)
func main() {
	fmt.Println("Start server...")	
	dsn := "postgres://app:pass@localhost:5432/db"
	Connect, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("can't connect to db %e",err)
	}


	var phone string
	for {
		num:=utils.ReadString(types.MenuAuther)
		switch num {
		case "1":
			service.Auther(Connect,phone)
			continue
		case "2":
			services.CustomerAtm(Connect)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}

	}
	
}
