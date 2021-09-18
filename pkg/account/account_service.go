package account

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/utils"

	"github.com/jackc/pgx/v4"
)

type AccountService struct {
	connect *pgx.Conn
}
func NewAccountServicce(connect *pgx.Conn) *AccountService{
	return &AccountService{connect: connect}
}

func (s *AccountService) GetAmountById(id int64) (int64, error) {
    var amount int64
    err:=s.connect.QueryRow(context.Background(),`select amount from account where id=$1`,id).Scan(&amount)
	if err != nil {
		utils.ErrCheck(err)
		return amount,err
	}
    return amount, err
}
func (s *AccountService) SetAmountById(amount,id int64)  error{
	_,err:=s.connect.Exec(context.Background(),`update account set amount = $1 where id = $2`,amount,id)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}

//Перевод 
func (s *AccountService) TransferMoneyByAccountId(payerAccountId,receiverAccountId int64, amount int64) error {
	payerAmount,err:=s.GetAmountById(payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount > payerAmount {
		log.Printf("не достаточно баланс")
		return err
	}
	receiverAmount,err:=s.GetAmountById(receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	newPayerAmount:=payerAmount-amount
	newreceiverAmount:=receiverAmount+amount

	err=s.SetAmountById(newPayerAmount,payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.SetAmountById(newreceiverAmount,receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Перевод Успешно отправлено!!!")
	
	return nil
}
