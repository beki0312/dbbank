package handler

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
		return
	}
	_, err = h.managerRepository.Register(r.Context(), managers)
	if err != nil {
		log.Print("Ошибка при регистрация менеджера")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	RespondJSON(w, managers)
}

//Авторизация Менеджера
func (h *ManagerHandler) ManagerToken(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		log.Print(err)
		return
	}
	token, err := h.managerRepository.Token(r.Context(), auther.Phone, auther.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print("Токен или логин неправильно")
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
		log.Println("ошибка при выводе список всех менеджеров")
		w.WriteHeader(http.StatusNotFound)
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
		log.Println("Ошибка при выводе список менеджера по Id")
		w.WriteHeader(http.StatusNotFound)
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
		log.Println("Не удалось удалит менеджера")
		w.WriteHeader(http.StatusNotImplemented)
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
		log.Println("Не удалось удалить токен менеджера")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	RespondJSON(w, item)
}
