package services

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"os"
)

//ServiceLoop - Для выбора из список
func (s *MoneyService) ServiceLoop(phone string) {
	var number string
	for {
		fmt.Println(types.MenuCustomer)
		fmt.Scan(&number)
		switch number {
		case "1":
			//TODO: список счетов пользователя
			s.ViewListAccounts(phone)
			continue
		case "2":
			//TODO: Перевести деньги другому клиенту
			s.CustomerPerevod()
			continue
		case "3":
			s.CustomerService()
			continue
		case "4":
			s.PayService()
			continue
		case "5":
			s.CustomerAtm()
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}
//CustomerAccount - Авторизация клиента
func (s *MoneyService) CustomerAccount(phone string) error{
	var password, pass string
	phone=utils.ReadString("Введите Лог: ")
	password=utils.ReadString("Введите парол: ")
	println("")
	ctx := context.Background()
	err := s.connect.QueryRow(ctx, `select password from customer where phone=$1`, phone).Scan(&pass)
	if err != nil {
		// fmt.Println("can't get login or password Customer")
		return err
	}
	if password == pass {
		fmt.Println("Хуш омадед Мизоч!!!")
		println("")
	} else {
		fmt.Println("Шумо логин ё паролро нодуруст дохил намудед!!!")
		fmt.Println(err)
		return err
	}
	s.ServiceLoop(phone)
	return nil
}

//ViewListAccounts - Посмотреть список счетов
func (s *MoneyService) ViewListAccounts(phone string) (Accounts []types.Account,err error) {	
	ctx :=context.Background()
	rows,err:=s.connect.Query(ctx,`SELECT account.id,account.customer_id,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.phone=$1`,phone)
	if err != nil {
		// log.Println("can't open accounts in customer %e",err)
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

// CustomerAtm - список банкомат
func (s *MoneyService) CustomerAtm() (Atms []types.Atm,err error)  {
	ctx:=context.Background()
	sql:=`select *from atm;`
	rows,err:=s.connect.Query(ctx,sql)
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
if rows.Err() !=nil {
	log.Print(err)
}
	return Atms,err
}

// CustomerService - список Услуг
func (s *MoneyService) CustomerService() (Atms []types.Services,err error)  {
	ctx:=context.Background()
	sql:=`select *from services;`
	rows,err:=s.connect.Query(ctx,sql)
	if err != nil {
		log.Printf("can't open service %e",err)
	}
	for rows.Next(){
	item:=types.Services{}
	err:=rows.Scan(&item.ID,&item.Name)
	if err != nil {
		log.Print(err)
		continue
	}
	Atms = append(Atms, item)
	fmt.Println(item)
}
// defer rows.Close()
if rows.Err() !=nil {
	log.Print(err)
}
	return Atms,err
}