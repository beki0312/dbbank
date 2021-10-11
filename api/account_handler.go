package api

import (
	"context"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
)



type AccountHandler struct {
	accountRepository 	*account.AccountRepository
}

func NewAccountHandler(accountRepository *account.AccountRepository) *AccountHandler {
	return &AccountHandler{accountRepository: accountRepository}
}

func (h *AccountHandler) GetAccountsAll(ctx context.Context) ([]*types.Account,error) {
	accounts,err:=h.accountRepository.Accounts()
	if err != nil {
		return nil, ErrInternal
	}
	return accounts,nil
}

func (h *AccountHandler) GetCustomerAccountById(ctx context.Context,Id int64) (*types.Account,error) {
	if (Id<=0) {
		return nil,ErrInternal
	}
	account,err:=h.accountRepository.GetAccountByCustomerPhone(Id)
	if err != nil {
		log.Println(err)
		return nil,ErrNotFound
	}
	return account,nil
}