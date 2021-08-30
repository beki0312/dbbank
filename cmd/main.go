package main

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/manager/service"
	"github.com/jackc/pgx/v4"
)


func main() {
	fmt.Println("Start server...")
	
	dsn := "postgres://app:pass@localhost:5432/db"

	Connect, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalln("Шумо парол ё логинро нодуруст дохил намудед")
	}
	err=service.ManagerAccount(Connect)
	if err != nil {
		log.Print(err)
	}
	

}




