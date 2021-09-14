package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mybankcli/pkg/customers/services"
	"mybankcli/pkg/types"
	"os"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

//Auther Авторизация, менеджера и клиента
func Auther(conn *pgx.Conn,phone string)  {
	var numberauther string
	for{
		fmt.Println(types.Auther)
		fmt.Scan(&numberauther)
		switch numberauther{
		case "1":
			//ManagerAccount - Авторизация Менеджера
			ManagerAccount(conn)
			continue
		case "2":
			//CustomerAccount акавнт клиента
			services.CustomerAccount(conn,phone)
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}	
}

//ManagerAccount - Авторизация Менеджера
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
		fmt.Println(err)
		return err
	}
	if password ==pass{
		fmt.Println("Хуш омадед Менедчер")
		println("")
	}else{
		fmt.Println("Шумо логин ё паролро нодуруст дохил намудед!!!")
		fmt.Println(err)
		return err
	}
	managerLoop(connect)
	return nil
}
//ManagerLoop - Меню менеджера
func managerLoop(con *pgx.Conn) {
	var number string
	for {
		fmt.Println(types.MenuManager)
		fmt.Scan(&number)
		switch number {
		case "1":
			// Добавить пользователя
			managerAddCustomer(con)
			continue
		case "2":
			// Добавить счет
			managerAddAccount(con)
			continue
		case "3":
			// Добавить услугу
			managerAddServices(con)
			continue
		case "4":
			// Экспорт список клиента
			exportCustomer(con)
			continue
		case "5":
			// экспорт список счетов
			exportAccounts(con)
			continue
		case "6":
			// экспорт список банкоматов
			exportAtm(con)
			continue
		case "7":
			importCustomer()
			continue
		case "10":
			//Добавить Банкоматов
			managerAddAtm(con)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
//ManagerAddCustomer- добавляет аккаунт клиента
func managerAddCustomer(connect *pgx.Conn,)  {
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
	err:=connect.QueryRow(ctx, `insert into customer (name,surname,phone,password)	values ($1,$2,$3,$4) returning id,name,surname,phone,password,active,created 
	`,name,surname,phone,password).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		return 
	}
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	cerr:=bcrypt.CompareHashAndPassword(hash,[]byte(password))
	if cerr != nil {
		log.Print("Invalid phone or password!")
	}
}
//ManagerAddAccount - добавляет счет для клиента
func managerAddAccount(connect *pgx.Conn)  {
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
//ManagerAddServices - добавляет название услуги
func managerAddServices(connect *pgx.Conn)  {
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
//ManagerAddAtm - Добавляет банкомата
func managerAddAtm(connect *pgx.Conn,)  {
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
	sql:=`insert into atm (numbers,district,address) values ($1,$2,$3) returning id,numbers,district,address`
	err:=connect.QueryRow(ctx,sql,numbers,district,address).Scan(&item.ID,&item.Numbers,&item.District,&item.Address)
	if err != nil {
		fmt.Printf("can't insert %e",err)
	}
}

// ExportCustomer - Экспортирует списка пользователей в json
func exportCustomer(conn *pgx.Conn) (Customers []types.Customer,err error) {
	ctx:=context.Background()
	sql:=`select *from customer`
	rows,err:=conn.Query(ctx,sql)
	CheckError(err)
	for rows.Next(){
		item:=types.Customer{}
		err:=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		CheckError(err)
	Customers = append(Customers, item)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Customers)

	file,err:=os.Create("data/Customer/Customers.json")
	CheckError(err)
	defer file.Close()
	io.Copy(file,buf)
	}
	return Customers,nil
}

// ExportAccounts - экспортирует списка счетов в json
func exportAccounts(conn *pgx.Conn) (Accounts []types.Account,err error) {
	ctx:=context.Background()
	sql:=`select *from account`
	rows,err:=conn.Query(ctx,sql)
	CheckError(err)
	for rows.Next(){
		item:=types.Account{}
		err:=rows.Scan(&item.ID,&item.Customer_Id,&item.Currency_code,&item.Account_Name,&item.Amount)
		CheckError(err)
	Accounts = append(Accounts, item)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Accounts)

	file,err:=os.Create("data/Accounts/Accounts.json")
	CheckError(err)
	defer file.Close()
	io.Copy(file,buf)
	}
	return Accounts,nil
}
// ExportAtm Экспортирует - списка банкоматов в json
func exportAtm(conn *pgx.Conn) (Atms []types.Atm,err error) {
	ctx:=context.Background()
	sql:=`select *from atm`
	rows,err:=conn.Query(ctx,sql)
	CheckError(err)
	for rows.Next(){
		item:=types.Atm{}
		err:=rows.Scan(&item.ID,&item.Numbers,&item.District,&item.Address)
		CheckError(err)
		Atms = append(Atms, item)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Atms)

	file,err:=os.Create("data/Atm/Atm.json")
	CheckError(err)
	defer file.Close()
	io.Copy(file,buf)
	}
	return Atms,nil
}

func importCustomer() ( atm types.Atm,err error)  {
	configFile,err:=os.Open("Bankomat.json")
	// defer configFile.Close()
	if err != nil {
		return atm,err
	}
	jsonParser :=json.NewDecoder(configFile)
	err=jsonParser.Decode(&atm)
	fmt.Println(&atm)
	return atm,err
}
//Ошибка
func CheckError(err error)  {
	if err != nil {
		log.Print(err)
	}
}