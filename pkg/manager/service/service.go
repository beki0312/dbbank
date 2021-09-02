package service

import (
	"context"
	"fmt"
	"mybankcli/pkg/customer/services"
	"mybankcli/pkg/types"
	"os"
	"github.com/jackc/pgx/v4"
)
func Auther(conn *pgx.Conn)  {
	var numberauther string
	for{
		fmt.Println(types.Auther)
		fmt.Scan(&numberauther)
		switch numberauther{
		case "1":
			ManagerAccount(conn)
			continue
		case "2":
			services.CustomerAccount(conn)
			continue
		case "q":
			os.Exit(0)
		}
	}	
}
func ManagerAccount(connect *pgx.Conn) error {
	var phone,password, pass string 
	fmt.Print("Введите Лог: ")
	fmt.Scan(&phone)
	fmt.Print("Введите парол: ")
	fmt.Scan(&password)
	println("")
	ctx:=context.Background()
	err:=connect.QueryRow(ctx, `select password from managers where phone=$1`,phone).Scan(&pass)
	if err != nil {
		fmt.Printf("can't get password %e",err)
		return err
	}
	if password ==pass{
		fmt.Println("Хуш омадед Менедчер")
		println("")
	}else{
		fmt.Println("Шумо паролро нодуруст дохил намудед!!!")
		fmt.Println(err)
		return err
	}
	Loop(connect)
	return nil
}
func Loop(con *pgx.Conn) {
	var cmd string
	for {
		fmt.Println(types.MenuManager)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			//TODO: Добавить пользователя
			ManagerAddCustomer(con)
			continue
		case "2":
			//TODO: Добавить счет
			ManagerAddAccount(con)
			continue
		case "3":
			//TODO: Добавить услугу
			ManagerAddServices(con)
			continue
		case "10":
			//TODO: Добавить Банкоматов
			ManagerAddAtm(con)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}
func ManagerAddCustomer(connect *pgx.Conn,)  {
	var name,surname,phone,password string 
			fmt.Print("Введите Имя: ")
			fmt.Scan(&name)
			fmt.Print("Введите Фамилия: ")
			fmt.Scan(&surname)
			fmt.Print("Введите Лог: ")
			fmt.Scan(&phone)
			fmt.Print("Введите парол: ")
			fmt.Scan(&password)
			println("")
	fmt.Println("Добалили клиент: Имя ",name, " фамиля ",surname," Логин ",phone," Парол ",password)
	println("")
	ctx:=context.Background()
	item:=types.Customer{}
	err:=connect.QueryRow(ctx, `insert into customer (name,surname,phone,password)
	values ($1,$2,$3,$4) returning id,name,surname,phone,password,active,created 
	`,name,surname,phone,password).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddAccount(connect *pgx.Conn)  {
	fmt.Println("Добавить Счеты ")
	var customerId,amount int64
	var accountname, currency string 
			fmt.Print("Введите id клиента: ")
			fmt.Scan(&customerId)
			fmt.Print("Ввведите код валюти TJS, RUB,USD,EUR: ")
			fmt.Scan(&currency)
			fmt.Print("Введите Счет: ")			
			fmt.Scan(&accountname)
			fmt.Print("Введите Баланс: ")			
			fmt.Scan(&amount)
			println("")
	fmt.Println("Добавили счет клиента id-клиента: ",customerId," код валюта: ",currency," номер счет: ",accountname," Баланс: ",amount)
	println("")
	ctx:=context.Background()
	item:=types.Account{}
	err:=connect.QueryRow(ctx, `insert into account (customer_id,currency_code,account_name,amount) values ($1,$2,$3,$4) returning id,customer_id,currency_code,account_name,amount 
	`,customerId,currency,accountname,amount).Scan(&item.ID,&item.Customer_Id,&item.Currency_code,&item.Account_Name,&item.Amount)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddServices(connect *pgx.Conn)  {
	fmt.Println("Добавить услуги ")
	var name string 
			fmt.Print("Введите название услуги: ")			
			fmt.Scan(&name)
			println("")
	fmt.Println("Добавили услуги : ",name)
	println("")
	ctx:=context.Background()
	item:=types.Services{}
	err:=connect.QueryRow(ctx, `insert into services (name) values ($1) returning id,name 
	`,name).Scan(&item.ID,&item.Name)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddAtm(connect *pgx.Conn,)  {
	var numbers int64
	var district, address string 
			fmt.Print("Введите № Банкомата: ")
			fmt.Scan(&numbers)
			fmt.Print("ВВедите район: ")
			fmt.Scan(&district)
			fmt.Print("Введите адрес Банкомата: ")
			fmt.Scan(&address)
			println("")
	fmt.Println("Добалили список Банкомат:  № ",numbers,", Район: ",district,", Адресс: ",address)
	println("")
	ctx:=context.Background()
	item:=types.Atm{}
	err:=connect.QueryRow(ctx, `insert into atm (numbers,district,address)
	values ($1,$2,$3) returning id,numbers,district,address,active,created 
	`,numbers,district,address).Scan(&item.ID,&item.Numbers,&item.District,&item.Address,&item.Active,&item.Created)
	if err != nil {
		fmt.Printf("can't insert %e",err)
	}
}
