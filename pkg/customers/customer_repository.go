package customers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)
var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")


type CustomerRepository struct {
	connect *pgx.Conn
}
func NewCustomerRepository(connect *pgx.Conn) *CustomerRepository {
	return &CustomerRepository{connect: connect}
}
// CustomerPerevod - Перевести деньги другому клиенту
func(s *CustomerRepository) CustomerTransfer() {
	// var number string
	for {
		fmt.Print("Переводы")
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//Перевод по номер счета
			s.CustomerTransferAccount()
			continue
		case "2":
			//перевод по номеру телефона
			s.PhoneTransaction()
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
// CustomerPerevodAccount - перевод по номеру счета
func (s *CustomerRepository) CustomerTransferAccount() error {
	var payerAccountId, receiverAccountId int64
	accountService:=account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру счета")
	payerAccount:=utils.ReadString("введите номер счета для снятия денег: ")
	amount:=utils.ReadInt("Введите сумму: ")
	receiverAccount:=utils.ReadString("Введите номер счета получателя: ")
	err := s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, payerAccount).Scan(&payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err = s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, receiverAccount).Scan(&receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return 	accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}
// CustomerPerevodPhone - перевод по номеру телефона
func (s *CustomerRepository) PhoneTransaction() error {
	var payerAccountId,receiverAccountId int64
	accountService:=account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру телефона")
	payerPhone:=utils.ReadString("Input payerPhone: ")
	amount:=utils.ReadInt("Input amount: ")
	receiverPhone:=utils.ReadString("Input receiverPhone: ")
	ctx:=context.Background()
	selectSql:=`select account.id from account left join customer on customer.id=account.customer_id where customer.phone=$1`
	err:=s.connect.QueryRow(ctx,selectSql,payerPhone).Scan(&payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.connect.QueryRow(ctx,selectSql,receiverPhone).Scan(&receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}
// PayService - Меню для оплата услуг
func (s *CustomerRepository) PayService() {
	// var number string
	for {
		fmt.Print("Оплатить услуги")
		fmt.Print(types.ServiceAdd)
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//попплнение баланса телефон
			s.PayServicePhone()
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
// PayServicePhone - попплнение баланса телефон
func (s *CustomerRepository) PayServicePhone() error {
	fmt.Println("услуга для пополнение номер телефона")
	var amuntaccount, amount int64
	var accountName, phone string
	accountName=utils.ReadString("Введите номер счета для снятия денег: ")
	amount=utils.ReadInt("Введите сумму: ")
	fmt.Print("Введите номер телефона: ")
	fmt.Scan(&phone)
	err := s.connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntaccount)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount > amuntaccount {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	_, err = s.connect.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntaccount-amount, accountName)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Успешно!!!")
	
	return nil
}
//Customers -для вывода список всех клиентов
func(s *CustomerRepository) Customers() ([]*types.Customer,error) {
	ctx:=context.Background()
	customers:=[]*types.Customer{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM customer`)
	if err != nil {
		return nil, ErrInternal
	}
	//
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, item)
	}
	return customers,nil
}

func (s *CustomerRepository) AllActiveCustomers() ([]*types.Customer,error) {
	ctx:=context.Background()
	customers:=[]*types.Customer{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM customer where active=true`)
	if err != nil {
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, item)
	}
	return customers,nil
}

func (s *CustomerRepository) CustomerById(id int64) (*types.Customer,error) {
	ctx:=context.Background()
	customers:=&types.Customer{}
	err:=s.connect.QueryRow(ctx,`select id,name,surname,phone,password,active,created from customer where id=$1`,
	id).Scan(&customers.ID,&customers.Name,&customers.SurName,&customers.Phone,&customers.Password,&customers.Active,&customers.Created)
	if err != nil {
		log.Println(err)
		return nil,ErrInternal
	}
	return customers,nil	
}
func (s *CustomerRepository) CustomersDeleteById(id int64) (*types.Customer,error) {
	ctx:=context.Background()
	cust := &types.Customer{}
	err := s.connect.QueryRow(ctx, `DELETE FROM customer WHERE id = $1`, 
	id).Scan(&cust.ID, &cust.Name, &cust.SurName,&cust.Phone,&cust.Password,&cust.Active, &cust.Created)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return cust, nil	
}
func HashPassword(password string) (string,error)  {
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes),err
}
func (s *CustomerRepository) CreateCustomers(customer *types.Customer) (*types.Customer,error) {
	ctx:=context.Background()
	item:=&types.Customer{}
	pass,_:=HashPassword(customer.Password)
	err:=s.connect.QueryRow(ctx,`insert into customer(id,name,surname,phone,password) values($1,$2,$3,$4,$5) returning id,name,surname,phone,password,active,created`,
	customer.ID,customer.Name,customer.SurName,customer.Phone,pass).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)	
	if err != nil {
		return nil,ErrInternal
	}
	return item,nil
}
//Block and Unblock customer by his id
func (s *CustomerRepository) CustomerBlockAndUnblockById(id int64,active bool) (*types.Customer,error) {
	ctx:=context.Background()
	customers:=&types.Customer{}
	err:=s.connect.QueryRow(ctx,`update customer set active =$1 where id=$2`,active,id).Scan(&customers.ID,&customers.Name,&customers.SurName,&customers.Phone,&customers.Password,&customers.Active,&customers.Created)
		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		return customers,nil
}

