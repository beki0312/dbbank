package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"
)

//Сервис - описывает обслуживание клиентов.
type ManagerHandler struct {
	connect           *pgx.Conn
	managerRepository *service.ManagerRepository
}

//NewServer - функция-конструктор для создания нового сервера.
func NewManagerHandler(connect *pgx.Conn, managerRepository *service.ManagerRepository) *ManagerHandler {
	return &ManagerHandler{connect: connect, managerRepository: managerRepository}
}

//Регистрация
func (h *ManagerHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var managers *types.Registration
	err := json.NewDecoder(r.Body).Decode(&managers)
	if err != nil {
		RespondBadRequest(w, "Получен не правильный тип")
		return
	}
	manager, err := h.managerRepository.Register(r.Context(), managers)
	if err != nil {
		RespondBadRequest(w,"Произошла ошибка во время регистрации менеджера")
				return
	}
	RespondJSON(w, manager)
}

//Авторизация Менеджера
func (h *ManagerHandler) ManagerToken(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		RespondBadRequest(w, "Получен не правильный тип")
		return
	}
	token, err := h.managerRepository.Token(r.Context(), auther.Phone, auther.Password)
	if err != nil {
		RespondUnauthorized(w, "Токен или логин неправильно")
		return
	}
	RespondJSON(w, token)
}

//найти токен менеджера идентификатор
func (s *ManagerHandler) TokenManagers(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.connect.QueryRow(ctx, `SELECT manager_id FROM managers_tokens WHERE token =$1`, token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		log.Print("Не удалос найти токена менеджера")
		return 0, err
	}
	return id, err
}

//Список Всех Менеджеров
func (h *ManagerHandler) GetAllManagers(w http.ResponseWriter, r *http.Request) {
	managers, err := h.managerRepository.ManagersAll(r.Context())
	if err != nil {
		RespondNotFound(w, "ошибка при выводе список всех менеджеров")
		return
	}
	RespondJSON(w, managers)
}

//Список Менеджеров по их Id
func (h *ManagerHandler) GetManagerById(w http.ResponseWriter, r *http.Request) {
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
	item, err := h.managerRepository.ManagersById(r.Context(), id)
	if err != nil {
		RespondNotFound(w, "Невозможно получить список менеджера по id")
		return
	}
	RespondJSON(w, item)
}

//Удалиение менеджеров по их Id
func (h *ManagerHandler) DeleteManagerById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := h.managerRepository.ManagersRemoveByID(r.Context(), id)
	if err != nil {
		RespondNotImplemented(w, "Не удалось удалит менеджера")
		return
	}
	RespondJSON(w, item)
}

//Удалиение Токен менеджера по их Id
func (h *ManagerHandler) DeleteTokenById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := h.managerRepository.ManagersTokenRemoveByID(r.Context(), id)
	if err != nil {
		RespondNotImplemented(w, "Не удалось удалит токен менеджера")
		return
	}
	RespondJSON(w, item)
}
