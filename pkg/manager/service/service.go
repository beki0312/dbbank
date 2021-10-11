package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"os"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)
type ManagerService struct {
	connect *pgx.Conn
}
func NewManagerServicce(connect *pgx.Conn) *ManagerService{
	return &ManagerService{connect: connect}
}
//Auther Авторизация, менеджера и клиента
func (s *ManagerService) Auther(phone string)  {
	customerService:=customers.NewCustomerServicce(s.connect)

	var numberauther string
	for{
		fmt.Println(types.Auther)
		fmt.Scan(&numberauther)
		switch numberauther{
		case "1":
			//ManagerAccount - Авторизация Менеджера
			s.ManagerAccount()
			continue
		case "2":
			//CustomerAccount акавнт клиента
			customerService.CustomerAccount(phone)
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
func (s *ManagerService) ManagerAccount() error {
	var pass string 
	phone:=utils.ReadString("Введите Лог: ")
	password:=utils.ReadString("Введите парол: ")
	println("")
	ctx:=context.Background()
	err:=s.connect.QueryRow(ctx, `select password from managers where phone=$1`,phone).Scan(&pass)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if password ==pass{
		fmt.Println("Хуш омадед Менедчер")
		println("")
	}else{
	fmt.Println("Шумо логин ё паролро нодуруст дохил намудед!!!")
	return err
	}
	s.managerLoop()
	return nil
}
//ManagerLoop - Меню менеджера
func (s *ManagerService) managerLoop() {
	var number string
	for {
		fmt.Println(types.MenuManager)
		fmt.Scan(&number)
		switch number {
		case "1":
			// Добавить пользователя
			s.managerAddCustomer()
			continue
		case "2":
			// Добавить счет
			s.managerAddAccount()
			continue
		case "3":
			// Добавить услугу
			s.managerAddServices()
			continue
		case "4":
			// Экспорт список клиента
			s.exportCustomer()
			continue
		case "5":
			// экспорт список счетов
			s.exportAccounts()
			continue
		case "6":
			// экспорт список банкоматов
			s.exportAtm()
			continue
		case "7":
			continue
		case "10":
			//Добавить Банкоматов
			s.managerAddAtm()
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
func HashPassword(password string) (string,error)  {
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes),err
}


//ManagerAddCustomer- добавляет аккаунт клиента
func (s *ManagerService) managerAddCustomer() error {
	name:=utils.ReadString("Введите имя: ")
	surName:=utils.ReadString("Введите Фамилия: ")
	phone:=utils.ReadString("Введите лог: ")
	password:=utils.ReadString("Введите парол: ")
	pass,_:=HashPassword(password)
	// if err != nil {
	// 	log.Print(err)
	// 	return err
	// }
	// PassString:=string(pass)
	println("")
	fmt.Println("Добалили клиент: Имя ",name, " фамиля ",surName," Логин ",phone," Парол ",password)
	println("")
	ctx:=context.Background()
	item:=types.Customer{}
	err:=s.connect.QueryRow(ctx, `insert into customer (name,surname,phone,password) values ($1,$2,$3,$4) returning id,name,surname,phone,password,active,created 
	`,name,surName,phone,pass).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}	
	return nil
}
//ManagerAddAccount - добавляет счет для клиента
func (s *ManagerService) managerAddAccount() error {
	fmt.Println("Добавить Счеты ")
	customerId:=utils.ReadInt("Введите id клиента: ")
	currency:=utils.ReadString("Ввведите код валюти TJS, RUB,USD,EUR: ")
	accountname:=utils.ReadString("Введите Счет: ")
	amount:=utils.ReadInt("Введите Баланс: ")
	println("")
	fmt.Println("Добавили счет клиента id-клиента: ",customerId," код валюта: ",currency," номер счет: ",accountname," Баланс: ",amount)
	println("")
	ctx:=context.Background()
	item:=types.Account{}
	err:=s.connect.QueryRow(ctx, `insert into account (customer_id,currency_code,account_name,amount) values ($1,$2,$3,$4) returning id,customer_id,currency_code,account_name,amount 
	`,customerId,currency,accountname,amount).Scan(&item.ID,&item.Customer_Id,&item.Currency_code,&item.Account_Name,&item.Amount)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}	
	return nil

}
//ManagerAddServices - добавляет название услуги
func (s *ManagerService) managerAddServices() error {
			fmt.Println("Добавить услуги ")
			name:=utils.ReadString("Введите название услуги: ")			
			println("")
			fmt.Println("Добавили услуги : ",name)
			println("")
	ctx:=context.Background()
	item:=types.Services{}
	err:=s.connect.QueryRow(ctx, `insert into services (name) values ($1) returning id,name 
	`,name).Scan(&item.ID,&item.Name)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}
//ManagerAddAtm - Добавляет банкомата
func (s *ManagerService) managerAddAtm() error {
	numbers:=utils.ReadInt("Введите № Банкомата: ")
	district:=utils.ReadString("ВВедите район: ")
	address:=utils.ReadString("Введите адрес Банкомата: ")
	println("")
	fmt.Println("Добалили список Банкомат:  № ",numbers,", Район: ",district,", Адресс: ",address)
	println("")
	ctx:=context.Background()
	item:=types.Atm{}
	sql:=`insert into atm (numbers,district,address) values ($1,$2,$3) returning id,numbers,district,address`
	err:=s.connect.QueryRow(ctx,sql,numbers,district,address).Scan(&item.ID,&item.Numbers,&item.District,&item.Address)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}

// ExportCustomer - Экспортирует списка пользователей в json
func (s ManagerService) exportCustomer() (Customers []types.Customer,err error) {
	ctx:=context.Background()
	sql:=`select *from customer`
	rows,err:=s.connect.Query(ctx,sql)
	if err != nil {
		utils.ErrCheck(err)
		return Customers,err
	}
	for rows.Next(){
		item:=types.Customer{}
		err:=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			utils.ErrCheck(err)
			return Customers,err
		}
			Customers = append(Customers, item)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Customers)

	file,err:=os.Create("data/Customer/Customers.json")
	utils.ErrCheck(err)
	defer file.Close()
	io.Copy(file,buf)
	}
	return Customers,nil
}

// ExportAccounts - экспортирует списка счетов в json
func (s *ManagerService) exportAccounts() (Accounts []types.Account,err error) {
	ctx:=context.Background()
	sql:=`select *from account`
	rows,err:=s.connect.Query(ctx,sql)
	if err != nil {
		utils.ErrCheck(err)
		return Accounts,err
	}
	for rows.Next(){
		item:=types.Account{}
		err:=rows.Scan(&item.ID,&item.Customer_Id,&item.Currency_code,&item.Account_Name,&item.Amount)
		if err != nil {
			utils.ErrCheck(err)
			return Accounts,err
		}
		Accounts = append(Accounts, item)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Accounts)

	file,err:=os.Create("data/Accounts/Accounts.json")
	utils.ErrCheck(err)
	defer file.Close()
	io.Copy(file,buf)
	}
	return Accounts,nil
}
// ExportAtm Экспортирует - списка банкоматов в json
func (s *ManagerService) exportAtm() (Atms []types.Atm,err error) {
	ctx:=context.Background()
	sql:=`select *from atm`
	rows,err:=s.connect.Query(ctx,sql)
	if err != nil {
		utils.ErrCheck(err)
		return Atms,err
	}
	for rows.Next(){
		item:=types.Atm{}
		err:=rows.Scan(&item.ID,&item.Numbers,&item.District,&item.Address)
		if err != nil {
			utils.ErrCheck(err)
			return Atms,err
		}
		Atms = append(Atms, item)
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(Atms)

	file,err:=os.Create("data/Atm/Atm.json")
	if err != nil {
		utils.ErrCheck(err)
		return Atms,err
	}
	defer file.Close()
	io.Copy(file,buf)
	}
	return Atms,nil
}

