package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"
)

//Сервис - описывает обслуживание клиентов.
type AccountHandler struct {
	accountRepository *account.AccountRepository
}

//NewServer - функция-конструктор для создания нового сервера.
func NewAccountHandler(accountRepository *account.AccountRepository) *AccountHandler {
	return &AccountHandler{accountRepository: accountRepository}
}

//Список счетов по Id
func (h *AccountHandler) GetAccountById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println("err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := h.accountRepository.GetAccountById(r.Context(), id)
	if err != nil {
		log.Printf("Не удалось вывести счет по id: %d, ошибка: %s", id, err)
		RespondNotFound(w, "Не удалось получить список счета по Id")
		return
	}
	RespondJSON(w, item)
}

//список счетов
func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	account, err := h.accountRepository.Accounts(r.Context())
	if err != nil {
		log.Printf("Не удалось вывести счет ошибка: %s", err)
		RespondNotFound(w, "Не удалось получить список всех счетов")
		return
	}
	RespondJSON(w, account)
}

//Добавление счет клиента
func (h *AccountHandler) PostNewAccounts(w http.ResponseWriter, r *http.Request) {
	var account *types.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Print(err)
		return
	}
	item, err := h.accountRepository.CreateAccounts(r.Context(), account)
	if err != nil {
		log.Printf("Не удалось добавить счет клиенту ошибка: %s", err)
		RespondBadRequest(w, "Не удалось создать новый счет для клиента")
		return
	}
	RespondJSON(w, item)
}
