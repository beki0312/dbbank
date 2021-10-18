package account
import (
	"context"
	"errors"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)

//Сервис - описывает обслуживание клиентов.
type AccountRepository struct {
	connect *pgx.Conn
}
//NewServer - функция-конструктор для создания нового сервера.
func NewAccountRepository(connect *pgx.Conn) *AccountRepository {
	return &AccountRepository{connect: connect}
}

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

//вывод счетов по их Id
func (s *AccountRepository) GetById(id int64) (types.Account, error) {
    var account types.Account
    err:=s.connect.QueryRow(context.Background(),`select id, customer_id,currency_code,account_name, amount from account where id=$1`,id).Scan(&account.ID,&account.Customer_Id,&account.Currency_code,&account.Account_Name,&account.Amount)
	if err != nil {
    utils.ErrCheck(err)
    return account,err
	}
    return account, err
}
//обновление баланс счета по Id
func (s *AccountRepository) SetAmountById(amount,id int64)  error{
	_,err:=s.connect.Exec(context.Background(),`update account set amount = $1 where id = $2`,amount,id)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}
//вывод Id счетов по номеру счета
func(s *AccountRepository) GetByIdAccountName(accountName string,accountId int64) error {
    err:=s.connect.QueryRow(context.Background(),`select id from account where account_name = $1`,accountName).Scan(&accountId)
	if err != nil {
    utils.ErrCheck(err)
    return err
	}
    return err
}
//Перевод по номеру телефона
func(s *AccountRepository) GetByIdCustomerPhone(payerPhone string, payerAccountId int64) error {
	ctx:=context.Background()
	err:=s.connect.QueryRow(ctx,`select account.id from account left join customer on 
	customer.id=account.customer_id where customer.phone=$1`,payerPhone).Scan(&payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
		}
		return err
	
}
//Таблица транзаксия
func (s *AccountRepository) CreateTransactions(payerAccountId,receiverAccountId,amount int64) error {
	ctx:=context.Background()
	item:=types.Transactions{}
	err:=s.connect.QueryRow(ctx, `insert into transactions (debet_account_id,credit_account_id,amount) values ($1,$2,$3) returning id,debet_account_id,credit_account_id,amount,date 
	`,payerAccountId,receiverAccountId,amount).Scan(&item.ID,&item.Debet_account_id,&item.Credit_account_id,&item.Amount,&item.Date)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return err
}
//Список счетов
func (s *AccountRepository) Accounts() ([]*types.Account,error) {
	ctx:=context.Background()
	accounts:=[]*types.Account{}
	rows,err:=s.connect.Query(ctx,`select * from account`)
	if err != nil {
		return nil, ErrInternal
	}
	for rows.Next(){
		account:=&types.Account{}
		err=rows.Scan(&account.ID,&account.Customer_Id,&account.Currency_code,&account.Account_Name,&account.Amount)
		if err != nil {
			log.Println(err)
		}
		accounts=append(accounts,account)
	}
	return accounts,nil
}

//Список счетов по их Id
func (s *AccountRepository) GetAccountById(id int64) (*types.Account,error) {
	accounts:=&types.Account{}
	err:=s.connect.QueryRow(context.Background(),`SELECT id,customer_id,currency_code, account_name,amount FROM account where id=$1`,id).Scan(
		&accounts.ID,&accounts.Customer_Id,&accounts.Currency_code,&accounts.Account_Name,&accounts.Amount)
	if err != nil {
		// utils.ErrCheck(err)
		return nil,err
	}
	return accounts,nil
}
//создание нового клиента
func (s *AccountRepository) CreateAccounts(account *types.Account) (*types.Account,error) {
	ctx:=context.Background()
	item:=&types.Account{}
	err:=s.connect.QueryRow(ctx,`insert into account(id,customer_id,currency_code,account_name,amount) values($1,$2,$3,$4,$5) returning id,customer_id,currency_code, account_name,amount`,
	account.ID,account.Customer_Id,account.Currency_code,account.Account_Name,account.Amount).Scan(&item.ID,&item.Customer_Id,&item.Currency_code,&item.Account_Name,&item.Amount)	
	if err != nil {
		return nil,ErrInternal
	}
	return item,nil
}