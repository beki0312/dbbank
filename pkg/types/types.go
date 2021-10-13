package types

import "time"


const MenuAuther = ` выберите из списка
1 - Авторизация
2 - Список банкоматов
q - Завершить работу `


const MenuManager = `
Выберите из списка
1. Добавить пользователя 
2. Добавить счет пользователю
3. Добавить услуги
4. Экспорт списка пользователей (Json XMl file) 
5. Экспорт списка счетов (Json XMl file)
6. Экспорт списка банкоматов (Json XMl file)
7. Импорт списка пользователей (Json XMl file) 
8. Импорт списка счетов (Json XMl file)
9. Импорт списка банкоматов (Json XMl file)
10. Добавить Банкомат
q.  Выйти из приложения`

const MenuCustomer=`
1. Посмотреть список счетов
2. Перевести деньги другому клиенту:
3. Список услуг
4. Оплатить услугу
5. Список банкоматов
q. Выйти (разлогинться)`

const MenuMoneyTransfer =`
1. по номеру счёта
2. по номеру телефона
q. Назад
`
const Auther=`
1. Менеджер
2. Клиент
q. назад
`

const ServiceAdd=`
1. Сотовые операторы
q. назад
`

//Manager представляет информацию о покупателе
type Manager struct {
	ID       	int64     `json:"id"`
	Name     	string    `json:"name"`
	SurName  	string    `json:"surname"`
	Phone    	string    `json:"phone"`
	Password 	string    `json:"password"`
	Active   	bool      `json:"active"`
	Created  	time.Time `json:"created"`
}

type Customer struct {
	ID       	int64     `json:"id"`
	Name     	string    `json:"name"`
	SurName  	string    `json:"surname"`
	Phone    	string    `json:"phone"`
	Password 	string    `json:"password"`
	Active   	bool      `json:"active"`
	Created  	time.Time `json:"created"`
}

type Services struct{
	ID 		int64		`json:"id"`
	Name	string		`json:"name"`
}
type Account struct{
	ID 						int64		`json:"id"`
	Customer_Id				int64		`json:"customer_id"`
	Currency_code			string		`json:"currency_code"`
	Account_Name			string		`json:"account_name"`
	Amount 					int64		`json:"amount"`
}
type AccountTransfer struct{
	Payer_Accont 		string		`json:"payerId"`
	Receiver_Account	string		`json:"receiverId"`
	Amount				int64		`json:"amount"`
}

type AccountPhoneTransactions struct{
	Payer_phone			string		`json:"payerPhone"`
	Receiver_Phone		string		`json:"receiverPhone"`
	Amount				int64		`json:"amount"`
}

type Atm struct{
	ID 			int64		`json:"id"`
	Numbers		int64		`json:"numbering"`
	District 	string		`json:"district"`
	Address		string		`json:"address"`	
}
type Transactions struct{
	ID					int64			`json:"id"`
	Debet_account_id	int64			`json:"payer_account_id"`
	Credit_account_id	int64			`json:"receiver_account_id"`
	Amount				int64			`json:"amount"`	
	Date 				time.Time		`json:"date"`
}
type Registration struct {
	FirstName     	string 	`json:"firs_name"`
	LastName		string	`json:"last_name"`
	Phone    	  	string 	`json:"phone"`
	Password 		string `json:"password"`
}

type Authers struct {
	Phone    	  	string 	`json:"phone"`
	Password 		string 	`json:"password"`
}

type Token struct {
	Token string `json:"token"`
}