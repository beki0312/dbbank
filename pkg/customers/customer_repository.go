package customers

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")
var ErrNoSuchUser = errors.New("no such user")
var ErrPhoneUsed = errors.New("phone already registered")
var ErrInvalidPassword = errors.New("invalid password")
var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExpired = errors.New("token expired")

//Сервис - описывает обслуживание клиентов.
type CustomerRepository struct {
	connect *pgx.Conn
}

//NewServer - функция-конструктор для создания нового сервера.
func NewCustomerRepository(connect *pgx.Conn) *CustomerRepository {
	return &CustomerRepository{connect: connect}
}

//Регистрация клиента
func (s *CustomerRepository) Register(ctx context.Context, reg *types.Registration) (*types.Customer, error) {
	item := &types.Customer{}
	hash, err := utils.HashPassword(reg.Password)
	if err != nil {
		log.Print("Невозможно хешировать пароль")
		return nil, ErrInternal
	}
	err = s.connect.QueryRow(ctx, `INSERT INTO customer (name,surname,phone, password)
	VALUES ($1,$2,$3,$4) ON CONFLICT (phone) DO NOTHING RETURNING id,name,surname,phone,password,active, created`, reg.FirstName, reg.LastName, reg.Phone, hash).Scan(
		&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrInternal
	}
	if err != nil {
		log.Print("ошибка при регистрация клиента")
		return nil, ErrInternal
	}
	return item, err
}

// method for generating a token
func (s *CustomerRepository) Token(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	err = s.connect.QueryRow(ctx, `SELECT id,password FROM customer WHERE phone =$1`, phone).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		log.Print("Неправильно ввели номер телефона")
		return "", ErrInternal
	}
	err = utils.CheckPasswordHass(password, hash)
	if err != nil {
		log.Print("неправильно")
		return "", ErrInvalidPassword
	}
	token, _ = utils.HashPassword(password)
	_, err = s.connect.Exec(context.Background(), `INSERT INTO customers_tokens(token,customer_id) VALUES($1,$2)`, token, id)
	if err != nil {
		log.Print("не удалось вставить токен и id клиента")
		return "", ErrInternal
	}
	return token, err
}

// CustomerPerevod - Перевести деньги другому клиенту
func (s *CustomerRepository) CustomerTransfer() {
	// var number string
	for {
		fmt.Print("Переводы")
		num := utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//Перевод по номер счета
			s.CustomerTransferAccount()
			continue
		case "2":
			//перевод по номеру телефона
			s.PhoneTransaction()
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
func (s *CustomerRepository) CustomerTransferAccount() error {
	var payerAccountId, receiverAccountId int64
	accountService := account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру счета")
	payerAccount := utils.ReadString("введите номер счета для снятия денег: ")
	amount := utils.ReadInt("Введите сумму: ")
	receiverAccount := utils.ReadString("Введите номер счета получателя: ")
	err := s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, payerAccount).Scan(&payerAccountId)
	if err != nil {
		log.Print("неправильно ввели номер счета отправителья или номер счет отправитель не сушествует")
		utils.ErrCheck(err)
		return err
	}
	err = s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, receiverAccount).Scan(&receiverAccountId)
	if err != nil {
		log.Print("неправильно ввели номер счета получателья или номер счет получателья не сушествует")
		utils.ErrCheck(err)
		return err
	}
	return accountService.TransferMoneyByAccountId(payerAccountId, receiverAccountId, amount)
}

// CustomerPerevodPhone - перевод по номеру телефона
func (s *CustomerRepository) PhoneTransaction() error {
	var payerAccountId, receiverAccountId int64
	accountService := account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру телефона")
	payerPhone := utils.ReadString("Input payerPhone: ")
	amount := utils.ReadInt("Input amount: ")
	receiverPhone := utils.ReadString("Input receiverPhone: ")
	ctx := context.Background()
	selectSql := `select account.id from account left join customer on customer.id=account.customer_id where customer.phone=$1`
	err := s.connect.QueryRow(ctx, selectSql, payerPhone).Scan(&payerAccountId)
	if err != nil {
		log.Print("неправильно ввели номер телефона или номер телефона не сушествует")
		utils.ErrCheck(err)
		return err
	}
	err = s.connect.QueryRow(ctx, selectSql, receiverPhone).Scan(&receiverAccountId)
	if err != nil {
		log.Print("неправильно ввели номер телефона получателя или номер телефона не сушествует")
		utils.ErrCheck(err)
		return err
	}
	return accountService.TransferMoneyByAccountId(payerAccountId, receiverAccountId, amount)
}

// PayService - Меню для оплата услуг
func (s *CustomerRepository) PayService() {
	// var number string
	for {
		fmt.Print("Оплатить услуги")
		fmt.Print(types.ServiceAdd)
		num := utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//попплнение баланса телефон
			s.PayServicePhone()
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}

// PayServicePhone - попплнение баланса телефон
func (s *CustomerRepository) PayServicePhone() error {
	fmt.Println("услуга для пополнение номер телефона")
	var amuntaccount, amount int64
	var accountName, phone string
	accountName = utils.ReadString("Введите номер счета для снятия денег: ")
	amount = utils.ReadInt("Введите сумму: ")
	fmt.Print("Введите номер телефона: ")
	fmt.Scan(&phone)
	err := s.connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntaccount)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount <= 0 {
		log.Print("Надо ввести сумму больше 0")
		return err
	}
	if amount > amuntaccount {
		fmt.Println("Недостаточно суммы на вашем балансе")
		return err
	}
	_, err = s.connect.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntaccount-amount, accountName)
	if err != nil {
		log.Print("перевод не успешно")
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Успешно!!!")

	return nil
}

//Customers -для вывода список всех клиентов
func (s *CustomerRepository) Customers(ctx context.Context) ([]*types.Customer, error) {
	customers := []*types.Customer{}
	rows, err := s.connect.Query(ctx, `SELECT *FROM customer`)
	if err != nil {
		log.Print("Не получилось вывести список клиентов")
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next() {
		item := &types.Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
		if err != nil {
			log.Println("Ошибка при скане customer")
		}
		customers = append(customers, item)
	}
	return customers, nil
}

//Список всех активный клиента
func (s *CustomerRepository) AllActiveCustomers() ([]*types.Customer, error) {
	ctx := context.Background()
	customers := []*types.Customer{}
	rows, err := s.connect.Query(ctx, `SELECT *FROM customer where active=true`)
	if err != nil {
		log.Print("Не удалось вывести список всех актывный клиентов")
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next() {
		item := &types.Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
		if err != nil {
			log.Println("ошибка при rows.scan ")
		}
		customers = append(customers, item)
	}
	return customers, nil
}

//Клиент по Id
func (s *CustomerRepository) CustomerById(ctx context.Context, id int64) (*types.Customer, error) {
	customers := &types.Customer{}
	err := s.connect.QueryRow(ctx, `select id,name,surname,phone,password,active,created from customer where id=$1`,
		id).Scan(&customers.ID, &customers.Name, &customers.SurName, &customers.Phone, &customers.Password, &customers.Active, &customers.Created)
	if err != nil {
		log.Println("ID клиента не существуеть или неправильно ввели Id клиента")
		return nil, ErrInternal
	}
	return customers, nil
}

//Удалит клиента по Id
func (s *CustomerRepository) CustomersDeleteById(ctx context.Context, id int64) (*types.Customer, error) {
	cust := &types.Customer{}
	err := s.connect.QueryRow(ctx, `DELETE FROM customer WHERE id = $1`,
		id).Scan(&cust.ID, &cust.Name, &cust.SurName, &cust.Phone, &cust.Password, &cust.Active, &cust.Created)
	if err != nil {
		log.Print("Не удалось удалить клиента")
		return nil, ErrInternal
	}
	return cust, nil
}

//Удалит счета по Id клиента
func (s *CustomerRepository) AccountsDeleteById(ctx context.Context, id int64) (*types.Account, error) {
	cust := &types.Account{}
	err := s.connect.QueryRow(ctx, `DELETE FROM account WHERE customer_id = $1`,
		id).Scan(&cust.ID, &cust.Customer_Id, &cust.Currency_code, &cust.Account_Name, &cust.Amount)
	if err != nil {
		log.Print("Не удалось удалить счеть клиента по id клиента ")
		return nil, ErrInternal
	}
	return cust, nil
}

//Регистрация нового клиента
func (s *CustomerRepository) CreateCustomers(customer *types.Customer) (*types.Customer, error) {
	ctx := context.Background()
	item := &types.Customer{}
	pass, _ := utils.HashPassword(customer.Password)
	err := s.connect.QueryRow(ctx, `insert into customer(id,name,surname,phone,password) values($1,$2,$3,$4,$5) returning id,name,surname,phone,password,active,created`,
		customer.ID, customer.Name, customer.SurName, customer.Phone, pass).Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
	if err != nil {
		log.Print("ошибка при регистрация клиента")
		return nil, ErrInternal
	}
	return item, nil
}

func (s *CustomerRepository) CustomerAtm(ctx context.Context) (Atms []*types.Atm, err error) {
	rows, err := s.connect.Query(ctx, `select *from atm`)
	if err != nil {
		log.Print("ошибка при выводе список банкоматов")
		return nil, ErrInternal
	}
	for rows.Next() {
		item := &types.Atm{}
		err := rows.Scan(&item.ID, &item.Numbers, &item.District, &item.Address)
		if err != nil {
			log.Print("Ошибка при rows.scan atm")
			continue
		}
		Atms = append(Atms, item)
		fmt.Println(item)
	}
	return Atms, err
}

//Добавить адресс банкомат
func (s *CustomerRepository) CreateAtms(ctx context.Context, atm *types.Atm) (*types.Atm, error) {
	item := &types.Atm{}
	err := s.connect.QueryRow(ctx, `insert into atm (id,numbers,district,address) values($1,$2,$3,$4) returning id,numbers,district,address`,
		atm.ID, atm.Numbers, atm.District, atm.Address).Scan(&item.ID, &item.Numbers, &item.District, &item.Address)
	if err != nil {
		log.Print("Не удалось добавить адрес банкомата")
		return nil, ErrInternal
	}
	return item, err
}

//Вывод всех список транзакция
func (s *CustomerRepository) HistoryTansfer(ctx context.Context) ([]*types.Transactions, error) {
	accounts := []*types.Transactions{}
	rows, err := s.connect.Query(ctx, `select *from transactions;`)
	if err != nil {
		log.Print("can't open transactions table")
		return nil, ErrInternal
	}
	for rows.Next() {
		item := &types.Transactions{}
		err = rows.Scan(&item.ID, &item.Debet_account_id, &item.Credit_account_id, &item.Amount, &item.Date)
		if err != nil {
			log.Println("error in rows.scan transactions")
		}
		accounts = append(accounts, item)
	}
	return accounts, err
}

//Удаление токен клиента по Id
func (s *CustomerRepository) CustomersTokenRemoveByID(ctx context.Context, id int64) (*types.Tokens, error) {
	tokens := &types.Tokens{}
	err := s.connect.QueryRow(ctx, `delete from customers_tokens where customer_id=$1`,
		id).Scan(&tokens.Id)
	if err != nil {
		log.Print("can't delete customers_tokens ")
		return nil, err
	}
	if id <= 0 {
		log.Print("номер id дольжен быть больше 0")
		return nil, ErrInternal
	}
	return tokens, err
}
