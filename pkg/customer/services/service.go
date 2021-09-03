package services

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/types"
	"github.com/jackc/pgx/v4"
)


//Посмотреть список счетов
func  ViewListAccounts(connect *pgx.Conn,phone string) (Accounts []types.Account,err error) {	
	ctx :=context.Background()
	rows,err:=connect.Query(ctx,`SELECT account.id,account.customer_id,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.phone=$1`,phone)
	if err != nil {
		log.Printf("can't open accounts in customer %e",err)
		return Accounts,err
	}
	for rows.Next(){
		account:=types.Account{}
		err:=rows.Scan(&account.ID,&account.Customer_Id,&account.Currency_code,&account.Account_Name,&account.Amount)
		if err != nil {
			log.Printf("can't scan account %e",err)
		}
		Accounts = append(Accounts, account)
		fmt.Println(account)
	}
	if rows.Err() !=nil {
		log.Printf("rows err %e",err)
		return nil,rows.Err()
	}
	
	return Accounts,nil
}
//Авторизация клиента
func CustomerAccount(connect *pgx.Conn,phone string) error{
	var password, pass string
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
	Loop(connect,phone)
	return nil
}
func CustomerAtm(conn *pgx.Conn) (Atms []types.Atm,err error)  {
	ctx:=context.Background()
	sql:=`select *from atm;`
	rows,err:=conn.Query(ctx,sql)
	if err != nil {
		log.Printf("can't open atm %e",err)
		
	}
	for rows.Next(){
	item:=types.Atm{}
	err:=rows.Scan(&item.ID,&item.Numbers,&item.District,&item.Address)
	if err != nil {
		log.Print(err)
		continue
	}
	Atms = append(Atms, item)
	fmt.Println(item)
}
// defer rows.Close()
err=rows.Err()
if err != nil {
	log.Print(err)
	return
}
	
	return Atms,err
}

//Для выбора из список
func Loop(con *pgx.Conn,phone string) {
	var cmd string
	for {
		fmt.Println(types.MenuCustomer)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			//TODO: Добавить пользователя
			ViewListAccounts(con,phone)
			continue
		case "2":
			//TODO: Добавить счет
			continue
		case "3":
			continue
		case "5":
			CustomerAtm(con)
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}