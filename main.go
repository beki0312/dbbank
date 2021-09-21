package main


import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)
// 1: account_repository.go
type AccountRepository struct {
	connect *pgx.Conn
}

type Account struct {

}

// 2: amount > Accoun 
// accountPayer, err := s.accountRepository.GetById(id)

// 
// func (s *AccountRepository) GetById(id int64) (Account, error) {
//     var account Account
//     err:=s.connect.QueryRow(context.Background(),`select id, name, customer_id, amount from account where id=$1`,id).Scan(&amount)
// 	if err != nil {
//     utils.ErrCheck(err)
//     return amount,err
// 	}
//     return amount, err
// }

// 3: func (s *AR) GetByCustomerId(customerId int64) ([]Account, error)
// 4: внедрить задание 3 в этот метод func (s *MoneyService) ViewListAccounts(phone string) (Accounts []types.Account,err error) {	
// 5: pkg.customer.service > pkg.customer 
// 6: MoneyService > CustomerService
// 7: func (s *MoneyService) Transactions(payerAccountId,receiverAccountId,amount int64) {
// перенести в pkg.account, создать TransactionRepository
// назвать метод Create(payerAccountId,receiverAccountId,amount int64) error
// прочитать чем отличается метод от функции
// Экспортируемые методы getById(), GetById()



func (s *AccountRepository) GetAmountById(id int64) (int64, error) {
    var amount int64
    err:=s.connect.QueryRow(context.Background(),`select amount from account where id=$1`,id).Scan(&amount)
	if err != nil {
    utils.ErrCheck(err)
    return amount,err
	}
    return amount, err
}


func (s *AccountRepository) SetAmountById(amount,id int64)  error{
	_,err:=s.connect.Exec(context.Background(),`update account set amount = $1 where id = $2`,amount,id)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}

func NewAccountServicce(connect *pgx.Conn) *AccountService{
  return &AccountService{accountRepository: &AccountRepository{connect: connect}}
}

type AccountService struct {
	accountRepository *AccountRepository
//   transactionRepository *....
}

//Перевод 
func (s *AccountService) TransferMoneyByAccountId(payerAccountId,receiverAccountId int64, amount int64) error {
	payerAmount,err:=s.accountRepository.GetAmountById(payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount > payerAmount {
		log.Printf("не достаточно баланс")
		return err
	}
	receiverAmount,err:=s.accountRepository.GetAmountById(receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	newPayerAmount:=payerAmount-amount
	newreceiverAmount:=receiverAmount+amount
  
//   err = s.transactionRepository.Create(payerAccountId,receiverAccountId,amount)
// 	if err != nil {
// 		utils.ErrCheck(err)
// 		return err
// 	}
  
	err=s.accountRepository.SetAmountById(newPayerAmount,payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.accountRepository.SetAmountById(newreceiverAmount,receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Перевод Успешно отправлено!!!")
	
	return nil
}