package account

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v4"
)
type AccountService struct {
	connect *pgx.Conn
}
func NewAccountServicce(connect *pgx.Conn) *AccountService{
	return &AccountService{connect: connect}
}
//Перевод 
func (s *AccountService) TransferMoneyByAccountId(payerAccountId,receiverAccountId int64, amount int64) error {
	var payerAmount,receiverAmount int64
	selectAmoundId:=`select amount from account where id = $1`
	err := s.connect.QueryRow(context.Background(), selectAmoundId, payerAccountId).Scan(&payerAmount)
	if err != nil {
		fmt.Print("can't get Balance")
		return err
	}
	if amount > payerAmount {
		log.Print("не достаточно баланс")
		return err
	}
	cerr := s.connect.QueryRow(context.Background(), selectAmoundId, receiverAccountId).Scan(&receiverAmount)
	if cerr != nil {
		return cerr
	}	

	newPayerAmount:=payerAmount-amount
	newreceiverAmount:=receiverAmount+amount
	updateAmount_AccountId:=`update account set amount = $1 where id = $2`
	_,err = s.connect.Exec(context.Background(),updateAmount_AccountId,newPayerAmount,payerAccountId)
if err != nil {
	return  err
}	
_, err = s.connect.Exec(context.Background(),updateAmount_AccountId,newreceiverAmount,receiverAccountId)
	if err != nil {
		return err
		} 
fmt.Println("Перевод Успешно отправлено!!!")
	
	return nil
}
