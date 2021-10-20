package handler

import (
	"context"
	"encoding/json"
	"log"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
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
func (h *ManagerHandler) ManagerRegistration(w http.ResponseWriter, r *http.Request) {
	var managers *types.Registration
	err := json.NewDecoder(r.Body).Decode(&managers)
	if err != nil {
		return
	}
	_, err = h.managerRepository.Register(r.Context(), managers)
	if err != nil {
		return
	}
	RespondJSON(w, managers)
}

//Авторизация Менеджера
func (h *ManagerHandler) GetManagersTokens(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		log.Print(err)
		return
	}
	token, err := h.managerRepository.Token(r.Context(), auther.Phone, auther.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(err)
		return
	}
	RespondJSON(w, token)
}

//найти токен менеджера идентификатор
func (s *ManagerHandler) IDByTokenManagers(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.connect.QueryRow(ctx, `SELECT manager_id FROM managers_tokens WHERE token =$1`, token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, err
}

//Список Всех Менеджеров
func (h *ManagerHandler) GetAllManagers(w http.ResponseWriter, r *http.Request) {
	managers, err := h.managerRepository.ManagersAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, managers)
}

//Список Менеджеров по их Id
func (h *ManagerHandler) GetManagersById(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Удалиение менеджеров по их Id
func (h *ManagerHandler) GetDeleteManagerById(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Удалиение Токен менеджера по их Id
func (h *ManagerHandler) GetDeleteManagerTokensById(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

func (h *ManagerHandler) PostManagers(ctx context.Context, managers *types.Manager) (*types.Manager, error) {
	if managers.ID <= 0 {
		return nil, ErrInternal
	}
	managers, err := h.managerRepository.CreateManagers(managers)
	if err != nil {
		return nil, ErrInternal
	}
	if managers == nil {
		return nil, ErrNotFound
	}
	return managers, err
}
