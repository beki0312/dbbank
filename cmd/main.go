package main

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
	"os"
	"github.com/jackc/pgx/v4"
)


func main() {
	fmt.Println("Start server...")
	
	dsn := "postgres://app:pass@localhost:5432/db"
	Connect, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalln("Шумо парол ё логинро нодуруст дохил намудед")
	}
	var number string
	for {
		fmt.Println(types.MenuAuther)
		fmt.Scan(&number)
		switch number {
		case "1":
			service.Auther(Connect)
			continue
		case "2":
			// service.ManagerAddAtm(Connect)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}


	}
	
}
