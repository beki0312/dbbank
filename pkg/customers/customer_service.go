package customers
import (
	"context"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"os"
	"fmt"
	"github.com/jackc/pgx/v4"
)


type CustomerService struct {
	customerRepository *CustomerRepository
}

func NewCustomerServicce(connect *pgx.Conn) *CustomerService{
	return &CustomerService{customerRepository: &CustomerRepository{connect: connect}}
}
//ServiceLoop - Для выбора из список
func (s *CustomerService) ServiceLoop(phone string) {
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
			s.customerRepository.CustomerTransfer()
			continue
		case "3":
			s.CustomerService()
			continue
		case "4":
			s.customerRepository.PayService()
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
func (s *CustomerService) CustomerAccount(phone string) error{
	var password, pass string
	phone=utils.ReadString("Введите Лог: ")
	password=utils.ReadString("Введите парол: ")
	println("")
	ctx := context.Background()
	err := s.customerRepository.connect.QueryRow(ctx, `select password from customer where phone=$1`, phone).Scan(&pass)
	if err != nil {
		utils.ErrCheck(err)
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
func (s *CustomerService) ViewListAccounts(phone string) (Accounts []types.Account,err error) {	
	ctx :=context.Background()
	rows,err:=s.customerRepository.connect.Query(ctx,`SELECT account.id,account.customer_id,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.phone=$1`,phone)
	if err != nil {
		utils.ErrCheck(err)
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
func (s *CustomerService) CustomerAtm() (Atms []types.Atm,err error)  {
	ctx:=context.Background()
	sql:=`select *from atm;`
	rows,err:=s.customerRepository.connect.Query(ctx,sql)
	if err != nil {
		utils.ErrCheck(err)
		return Atms ,err
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
func (s *CustomerService) CustomerService() (Atms []types.Services,err error)  {
	ctx:=context.Background()
	sql:=`select *from services;`
	rows,err:=s.customerRepository.connect.Query(ctx,sql)
	if err != nil {
		utils.ErrCheck(err)
		return Atms,err
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