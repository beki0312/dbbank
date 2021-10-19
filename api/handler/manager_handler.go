package handler

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"
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

//Регистрация Менеджера
func (h *ManagerHandler) RegistersManagers(ctx context.Context, item *types.Registration) (*types.Manager, error) {
	manager := &types.Manager{}
	registration, err := h.managerRepository.Register(item)
	if err != nil {
		return nil, err
	}
	if registration == nil {
		return nil, ErrNotFound
	}
	return manager, err
}

//Авторизация Менеджера
func (h *ManagerHandler) GetManagersToken(ctx context.Context, item *types.Authers) (token string, err error) {
	token, err = h.managerRepository.Token(item.Phone, item.Password)
	if err != nil {
		return "", err
	}
	return token, err
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

//Получить весь менеджер
func (h *ManagerHandler) GetManagersAll(ctx context.Context) ([]*types.Manager, error) {
	managers, err := h.managerRepository.ManagersAll()
	if err != nil {
		return nil, ErrInternal
	}

	return managers, nil
}

//Получить все активный менеджер
func (h *ManagerHandler) GetManagersAllActive(ctx context.Context) ([]*types.Manager, error) {
	managers, err := h.managerRepository.ManagersAllActive()
	if err != nil {
		return nil, ErrInternal
	}
	return managers, nil
}

//Получить менеджеров по Id
func (h *ManagerHandler) GetManagersById(ctx context.Context, id int64) (*types.Manager, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	managers, err := h.managerRepository.ManagersById(id)
	if err != nil {
		return nil, ErrInternal
	}
	if managers == nil {
		return nil, ErrNotFound
	}
	return managers, nil
}

// Удалеит Менеджера по их Id
func (h *ManagerHandler) GetManagersRemoveByID(ctx context.Context, id int64) (*types.Manager, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	managers, err := h.managerRepository.ManagersRemoveByID(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if managers == nil {
		return nil, ErrNotFound
	}
	return managers, nil
}

// Удалеит Токена по их Id
func (h *ManagerHandler) GetManagersTokensRemoveByID(ctx context.Context, id int64) (*types.Tokens, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	tokens, err := h.managerRepository.ManagersTokenRemoveByID(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if tokens == nil {
		return nil, ErrNotFound
	}
	return tokens, nil
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
