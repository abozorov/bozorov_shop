package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type updateRole struct {
	Role string `json:"role"`
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
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format(time.RFC822Z),
	}
}

type responseOrder struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Product   string  `json:"product"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

func newResponseOrder(o models.Order) *responseOrder {
	return &responseOrder{
		ID:        o.ID,
		UserID:    o.UserID,
		Product:   o.Product,
		Price:     o.Price,
		Status:    o.Status,
		CreatedAt: o.CreatedAt.Format(time.RFC822Z),
	}
}

type responseProfile struct {
	*responseUser
	UserOrders []responseOrder `json:"orders"`
}

func newResponseProfile(prof models.Profile) *responseProfile {
	orders := make([]responseOrder, 0, 1000)
	for _, v := range prof.UserOrders {
		orders = append(orders, *newResponseOrder(v))
	}
	return &responseProfile{
		responseUser: newResponseUser(*prof.User),
		UserOrders:   orders,
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
	w.Write([]byte("User Registered"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	tokens, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		h.logger.Error("user_handler.Login: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrSomethingWentWrong)
	}
}

// принимаем код, после проверяяем и создаем запрос для похранения в БД
func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req models.Verification
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("user_handler.Verify: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// sending code
	err = h.service.Verification(r.Context(), req)
	if err != nil {
		h.logger.Error("user_handler.Verify: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("send for verrification"))
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get tokens
	rToken := models.Tokens{}
	err := json.NewDecoder(r.Body).Decode(&rToken)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// try refresh tokens
	tokens, err := h.service.RefreshTokens(r.Context(), rToken.Refresh)
	if err != nil {
		h.logger.Error("user_handler.RefreshToken: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		h.logger.Error("user_handler.RefreshToken: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get tokens
	tokens := models.Tokens{}
	err := json.NewDecoder(r.Body).Decode(&tokens)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// try logout
	err = h.service.Logout(r.Context(), tokens)
	if err != nil {
		h.logger.Error("user_handler.Logout: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	w.Write([]byte("logout success"))
}

// Admin
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("user_handler.GetAllUsers: ", zap.String("error", err.Error()))
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
		h.logger.Error("user_handler.GetAllUsers: ", zap.String("error", err.Error()))
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
		h.logger.Error("user_handler.GetMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := *newResponseUser(*usr)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get by id
	prof, err := h.service.GetProfile(r.Context(), id)
	if err != nil {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := *newResponseProfile(*prof)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.UpdateMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
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
		h.logger.Error("user_handler.UpdateMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("User updated"))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get id
	id, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.UpdateMe: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// get user
	pass := models.UpdatePassword{}
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}
	pass.UserID = id

	// creating & transform models.User -> user
	err = h.service.UpdatePassword(r.Context(), pass)
	if err != nil {
		h.logger.Error("user_handler.UpdateMe: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("User password updated"))
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("user_handler.UpdateUserRole: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get my id
	myID, ok := r.Context().Value(mycontext.UserIDKey).(int)
	if !ok {
		h.logger.Error("user_handler.UpdateUserRole: ", zap.String("error", errs.ErrIncorrectLoginOrPassword.Error()))
		errs.ErrsToHttp(w, errs.ErrIncorrectLoginOrPassword)
		return
	}

	// chek my
	if myID == userID {
		h.logger.Error("user_handler.UpdateUserRole: ", zap.String("error", "you can't change your role"))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get role
	role := updateRole{}
	err = json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// update role
	err = h.service.UpdateUserRole(r.Context(), models.User{
		ID:   userID,
		Role: role.Role,
	})
	if err != nil {
		h.logger.Error("user_handler.UpdateUserRole: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte(fmt.Sprintf("Role for user %d updated", userID)))
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
	w.Write([]byte(fmt.Sprintf("user %d Delete", id)))
}
