package userhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abozorov/bozorov_shop/internal/models"
	mycontext "github.com/abozorov/bozorov_shop/internal/my_context"
	userservice "github.com/abozorov/bozorov_shop/internal/service/user"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/abozorov/bozorov_shop/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *userservice.UserService
	logger  *logger.Logger
}

func NewUserHandler(service *userservice.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

type updateUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type responseUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func newResponseUser(u models.User) *responseUser {
	return &responseUser{
		CreatedAt: u.CreatedAt.Format(time.RFC822Z),
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	if err = h.service.Register(r.Context(), req); err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Created"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		h.logger.Error("user_handler.Register: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrSomethingWentWrong)
	}
}

// Admin
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*2)
	defer cancle()
	users, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := make([]responseUser, 0, len(users))
	for _, v := range users {
		resp = append(resp, *newResponseUser(v))
	}

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get by id
	usr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("user_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := *newResponseUser(*usr)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get user
	usr := updateUser{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Update(r.Context(), models.User{
		ID:    id,
		Name:  usr.Name,
		Phone: usr.Phone,
	})
	if err != nil {
		h.logger.Error("user_handler.Update: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("User updated"))
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.DeleteMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get by id
	err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		h.logger.Error("user_handler.DeleteMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	w.Write([]byte(fmt.Sprintf("user %d deleted", id)))
}
