package services

import (
	"context"
	// "errors"
	"fmt"
	"log"
	"mybankcli/pkg/types"
	"os"
	"github.com/jackc/pgx/v4"
)



//Посмотреть список счетов
func  ViewListAccounts(connect *pgx.Conn) error {
	// items:=types.Customer{}
	// item:=types.Customer{}
	ctx :=context.Background()
	err:=connect.QueryRow(ctx,`SELECT customer.name,account.currency_code, account.account_name,account.amount 
	FROM account 
	JOIN customer ON account.customer_id = customer.id 
	where account.customer_id=customer.id`)
	
	if err != nil {
		log.Printf("can't open account %e",err)
	}
	return nil
	
}
//Авторизация клиента
func CustomerAccount(connect *pgx.Conn) error{
	var phone,password, pass string
	fmt.Print("Введите Лог: ")
	fmt.Scan(&phone)
	fmt.Print("Введите парол: ")
	fmt.Scan(&password)
	println("")
	ctx := context.Background()
	err := connect.QueryRow(ctx, `select password from customer where phone=$1`, phone).Scan(&pass)
	
	if err != nil {
		fmt.Printf("can't get password Customer %e", err)
		return err
	}
	if password == pass {
		fmt.Println("Хуш омадед Мизоч!!!")
		println("")
	} else {
		fmt.Println("Шумо паролро нодуруст дохил намудед!!!")
		fmt.Println(err)
		return err
	}
	Loop(connect)
	return nil
}
//Для выбора из список
func Loop(con *pgx.Conn) {
	var cmd string
	for {
		fmt.Println(types.MenuCustomer)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			//TODO: Добавить пользователя

			continue
		case "2":
			//TODO: Добавить счет
			ViewListAccounts(con)
			continue
		case "3":
			//TODO: Добавить услугу
			// ManagerAddServices(con)
			continue
		case "10":
			//TODO: Добавить Банкоматов
			// ManagerAddAtm(con)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}
