package services

import (
	"context"
	"fmt"
	"mybankcli/pkg/types"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// CustomerPerevod - Перевести деньги другому клиенту
func CustomerPerevod(conn *pgx.Conn) {
	var number string
	for {
		fmt.Print("Переводы")
		fmt.Print(types.MenuMoneyTransfer)
		fmt.Scan(&number)
		switch number {
		case "1":
			//Перевод по номер счета
			CustomerPerevodAccount(conn)
			continue
		case "2":
			//перевод по номеру телефона
			CustomerPerevodPhone(conn)
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
func CustomerPerevodAccount(connect *pgx.Conn) error {
	fmt.Println("Перевод по номеру счета")
	var amuntName, amountCustomer, amount int64
	var accountName, accountCustomer string
	fmt.Print("Введите номер счета для снятия денег: ")
	fmt.Scan(&accountName)
	fmt.Print("Введите сумму: ")
	fmt.Scan(&amount)
	fmt.Print("Введите номер счета получателя: ")
	fmt.Scan(&accountCustomer)

	err := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntName)
	err = connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountCustomer).Scan(&amountCustomer)

	if err != nil {
		fmt.Printf("can't get Balance %e", err)
		return err
	}
	if amount > amuntName {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	tx, err := connect.Begin(context.Background())
	if err != nil {
		fmt.Printf("can't open transaction %e", err)
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())

		}
		err = tx.Commit(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntName-amount, accountName)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amountCustomer+amount, accountCustomer)
	if err != nil {
		return err
	} else {
		fmt.Println("Перевод Успешно отправлено!!!")
		fmt.Println("")
	}

	return nil
}
// CustomerPerevodPhone - перевод по номеру телефона
func CustomerPerevodPhone(connect *pgx.Conn) error {
	fmt.Println("Перевод по номеру телефона")
	var amuntName, amountCustomer, amount int64
	var accountPhone, accountCustomer string
	fmt.Print("Введите номер телефона для снятия денег: ")
	fmt.Scan(&accountPhone)
	fmt.Print("Введите сумму для снятия: ")
	fmt.Scan(&amount)
	fmt.Print("Введите номер телефон получателя: ")
	fmt.Scan(&accountCustomer)
	ctx := context.Background()
	err := connect.QueryRow(ctx, `SELECT customer.name,customer.phone,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.phone=$1`, accountPhone).Scan(&amuntName)
	err = connect.QueryRow(ctx, `SELECT account.id,account.customer_id,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.phone=$1`, accountCustomer).Scan(&amountCustomer)
	if err != nil {
		fmt.Printf("can't get Balance %e", err)
		return err
	}
	if amount > amuntName {
		err = errors.New("Not enough amount on your balance")
		return err
	}
	tx, err := connect.Begin(context.Background())
	if err != nil {
		fmt.Printf("can't open transaction %e", err)
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
		}
		err = tx.Commit(context.Background())
		ErrCheck(err)
	}()
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where phone = $2`, amuntName-amount, accountPhone)
	ErrCheck(err)
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where phone = $2`, amountCustomer+amount, accountCustomer)
	if err != nil {
		return err
	} else {
		fmt.Println("Перевод Успешно отправлено!!!")
		fmt.Println("")
	}
	return nil
}
// ошибка
func ErrCheck(err error)  {
	if err != nil {
		fmt.Print(err)
	}
}