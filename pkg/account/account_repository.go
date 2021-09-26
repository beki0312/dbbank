package account

import (
	"context"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)
type AccountRepository struct {
	connect *pgx.Conn
}

func (s *AccountRepository) GetById(id int64) (types.Account, error) {
    var account types.Account
    err:=s.connect.QueryRow(context.Background(),`select id, customer_id,currency_code,account_name, amount from account where id=$1`,id).Scan(&account.ID,&account.Customer_Id,&account.Currency_code,&account.Account_Name,&account.Amount)
	if err != nil {
    utils.ErrCheck(err)
    return account,err
	}
    return account, err
}
func (s *AccountRepository) SetAmountById(amount,id int64)  error{
	_,err:=s.connect.Exec(context.Background(),`update account set amount = $1 where id = $2`,amount,id)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}
