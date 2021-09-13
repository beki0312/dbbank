package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

//
func ReadString(lable string) string {
	var input string
	fmt.Print(lable)
	fmt.Scan(&input)
	return input
}
func ReadInt(lable string) int64  {
	var input int64
	fmt.Print(lable)
	fmt.Scan(&input)
	return input
}
//Перевод 
func TransferMoneyByAccountName(connect *pgx.Conn,accountName,accountCustomer string, amount int64) error {
	var amountPayer,amountReceiver int64
	err := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amountPayer)
	if err != nil {
		fmt.Print("can't get Balance")
		return err
	}
	if amount > amountPayer {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	cerr := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountCustomer).Scan(&amountReceiver)
	if cerr != nil {
		return cerr
	}	
	return nil
}
