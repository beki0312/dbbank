package account

// import (
// 	"context"
// 	"mybankcli/pkg/types"
// 	"mybankcli/pkg/utils"
// 	"github.com/jackc/pgx/v4"
// )

// type TransactionRepository struct {
// 	connect *pgx.Conn
// }
// //Таблица транзаксия
// func (s *AccountService) CreateTransactions(payerAccountId,receiverAccountId,amount int64) error {
// 	ctx:=context.Background()
// 	item:=types.Transactions{}
// 	err:=s.transactionRepository.connect.QueryRow(ctx, `insert into transactions (debet_account_id,credit_account_id,amount) values ($1,$2,$3) returning id,debet_account_id,credit_account_id,amount,date 
// 	`,payerAccountId,receiverAccountId,amount).Scan(&item.ID,&item.Debet_account_id,&item.Credit_account_id,&item.Amount,&item.Date)
// 	if err != nil {
// 		utils.ErrCheck(err)
// 		return err
// 	}
// 	return err
// }
