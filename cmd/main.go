package main

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"os"
	"github.com/jackc/pgx/v4"
)
func main() {
	fmt.Println("Start server...")	
	dsn := "postgres://app:pass@localhost:5432/db"
	connect, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("can't connect to db %e",err)
	}
 customersService:=customers.NewCustomerServicce(connect)
 managerService:=service.NewManagerServicce(connect)

	var phone string
	for {
		num:=utils.ReadString(types.MenuAuther)
		switch num {
		case "1":
			managerService.Auther(phone)
			continue
		case "2":
			customersService.CustomerAtm()
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}

	}
	
}
