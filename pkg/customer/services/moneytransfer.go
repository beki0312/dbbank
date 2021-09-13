package services

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// CustomerPerevod - Перевести деньги другому клиенту
func CustomerPerevod(conn *pgx.Conn) {
	// var number string
	for {
		fmt.Print("Переводы")
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//Перевод по номер счета
			CustomerPerevodAccount(conn)
			continue
		case "2":
			//перевод по номеру телефона
			PhoneTransaction(conn)
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
	var amountPayer, amountReceiver int64
	accountName:=utils.ReadString("введите номер счета для снятия денег: ")
	amount:=utils.ReadInt("Введите сумму: ")
	accountCustomer:=utils.ReadString("Введите номер счета получателя: ")
	// err := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amountPayer)
	// if err != nil {
	// 	fmt.Print("can't get Balance")
	// 	return err
	// }
	// cerr := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountCustomer).Scan(&amountReceiver)
	// if cerr != nil {
	// 	return err
	// }
	// if amount > amountPayer {
	// 	err = errors.New("Not enough amount on your balance")
	// 	fmt.Println(err)
	// 	return err
	// }
	err:=utils.TransferMoneyByAccountName(connect,accountName,accountCustomer,amount)
	туц
	if err != nil {
		log.Printf("can't %e",err)
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
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amountPayer-amount, accountName)
	ErrCheck(err)
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amountReceiver+amount, accountCustomer)
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
		return
	}
}
// CustomerPerevodPhone - перевод по номеру телефона
func PhoneTransaction(connect *pgx.Conn ) error {
	var payerAmount,receiverAmount int64
	fmt.Println("Перевод по номеру телефона")
	payerPhone:=utils.ReadString("Input payerPhone: ")
	amount:=utils.ReadInt("Input amount: ")
	receiverPhone:=utils.ReadString("Input receiverPhone: ")
	ctx:=context.Background()
	err:=connect.QueryRow(ctx,`select account.amount from account left join customer on customer.id=account.customer_id where customer.phone=$1`,payerPhone).Scan(&payerAmount)
	cerr:=connect.QueryRow(ctx,`select account.amount from account left join customer on customer.id=account.customer_id where customer.phone=$1`,receiverPhone).Scan(&receiverAmount)
	ErrCheck(err)
	if cerr != nil {
		return err
	}
	if amount>payerAmount {
		return errors.New("Не достаточно средств")
	}
	tx, err := connect.Begin(ctx)
	if err != nil {
		fmt.Printf("can't open transaction %e", err)
		return err
	}
	defer func() {
		if cerr != nil {
			cerr = tx.Rollback(ctx)
		}
		gerr := tx.Commit(ctx)
		if gerr != nil {
			fmt.Println(err)
		}
	}()
	_, err = tx.Exec(ctx, `  update account a set amount =$1 	from customer c where c.id=a.customer_id and c.phone=$2`, payerAmount-amount, payerPhone)
	ErrCheck(err)
	_, err = tx.Exec(ctx, ` update account a set amount =$1 	from customer c where c.id=a.customer_id and c.phone=$2`, receiverAmount+amount, receiverPhone)
	if err != nil {
		return err
	} else {
		fmt.Println("Перевод Успешно отправлено!!!")
		fmt.Println("")
	}
	return nil
}
