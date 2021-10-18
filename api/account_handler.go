package api

import (
	"context"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
)
//Сервис - описывает обслуживание клиентов.
type AccountHandler struct {
	accountRepository 	*account.AccountRepository
}
//NewServer - функция-конструктор для создания нового сервера.
func NewAccountHandler(accountRepository *account.AccountRepository) *AccountHandler {
	return &AccountHandler{accountRepository: accountRepository}
}
//Вывод всех список счетов
func (h *AccountHandler) GetAccountsAll(ctx context.Context) ([]*types.Account,error) {
	accounts,err:=h.accountRepository.Accounts()
	if err != nil {
		return nil, ErrInternal
	}
	return accounts,nil
}
//вывод список счетов по Id
func (h *AccountHandler) GetCustomerAccountById(ctx context.Context,id int64) (*types.Account,error) {
	if (id<=0) {
		return nil,ErrInternal
	}
	account,err:=h.accountRepository.GetAccountById(id)
	if err != nil {
		return nil,ErrNotFound
	}
	if account==nil{

		return nil,ErrNotFound
	}
	return account,nil
}
//добавление счет клиента
func (h *CustomerHandler) PostAccounts(ctx context.Context, account *types.Account) (*types.Account,error) {
	if (account.ID<=0) {
		return nil,ErrInternal
	}
	accounts,err:=h.accountRepository.CreateAccounts(account)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if accounts==nil {
		return nil,ErrNotFound
	}
	return accounts,nil
}