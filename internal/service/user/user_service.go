package userservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/abozorov/bozorov_shop/internal/repo"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/abozorov/bozorov_shop/pkg/jwt"
	mailsender "github.com/abozorov/bozorov_shop/pkg/mail_sender"
	"github.com/abozorov/bozorov_shop/pkg/password"
	refreshtoken "github.com/abozorov/bozorov_shop/pkg/refresh_token"
	"github.com/patrickmn/go-cache"
)

type UserService struct {
	userR            repo.UserRepo
	orderR           repo.OrderRepo
	refreshTokenRepo repo.RefreshTokenRepo
	jwt              *jwt.JWTSecret
	memCache         *cache.Cache
	mailSender       *mailsender.MailSender
}

func NewUserService(
	userR repo.UserRepo,
	orderR repo.OrderRepo,
	refreshTokenRepo repo.RefreshTokenRepo,
	jwt *jwt.JWTSecret,
	memCache *cache.Cache,
	mailsender *mailsender.MailSender) *UserService {

	return &UserService{
		userR:            userR,
		orderR:           orderR,
		refreshTokenRepo: refreshTokenRepo,
		jwt:              jwt,
		memCache:         memCache,
		mailSender:       mailsender,
	}
}

type sendOtp struct {
	code       int
	user       *models.User
	attemptOTP *int
}

func (u *UserService) Verification(ctx context.Context, req models.Verification) error {
	// check mem cash for exist
	exists, err := u.userR.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("user_service.Verification: %w", err)
	}
	if exists {
		return fmt.Errorf("user_service.Verification: %w", errs.ErrUserAlreadyExists)
	}

	// check mem cash for exist
	user, ok := u.memCache.Get(req.Email)
	if !ok {
		return fmt.Errorf("user_service.Verification: %w", errs.ErrVerifyingFailed)
	}
	defer func() {
		*user.(sendOtp).attemptOTP++
	}()

	if *user.(sendOtp).attemptOTP > 2 {
		u.memCache.Delete(req.Email)
		return fmt.Errorf("user_service.Verification: %w", errs.ErrToManyAttempt)
	}

	if user.(sendOtp).code != req.Code {
		return fmt.Errorf("user_service.Verification: %w", errs.ErrIncorrectOTPCode)
	}

	err = u.userR.Add(ctx, *user.(sendOtp).user)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}
	u.memCache.Delete(req.Email)

	return nil
}

func (u *UserService) Register(ctx context.Context, request models.RegisterRequest) error {
	err := request.Validate()
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	// check for exist in db
	exists, err := u.userR.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}
	if exists {
		return fmt.Errorf("user_service.Register: %w", errs.ErrUserAlreadyExists)
	}

	// check for exist in memcache
	_, ok := u.memCache.Get(request.Email)
	if ok {
		return fmt.Errorf("user_service.Register: %w", errs.ErrUserNotBeenVerified)
	}

	//
	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
	}

	// saving user in memCach & waiting user for verification
	// generate otp code & send user
	otpCode := rand.Int()%899999 + 100000
	attempt := 0

	// save in cache
	u.memCache.Set(user.Email, sendOtp{
		code:       otpCode,
		user:       &user,
		attemptOTP: &attempt,
	}, cache.DefaultExpiration)

	// sending email
	err = u.mailSender.SendMail(user.Email, strconv.Itoa(otpCode))
	if err != nil {
		u.memCache.Delete(user.Email)
		return fmt.Errorf("user_service.Register: %w", err)
	}

	//
	return nil
}

func (u *UserService) Login(ctx context.Context, request models.LoginRequest) (*models.Tokens, error) {
	err := request.Validate()
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// get user by email
	user, err := u.userR.GetByEmail(ctx, request.Email)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// check for delete
	if !user.DeletedAt.IsZero() {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", errs.ErrUserNotFound)
	}

	// compare password
	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// generate tokens
	jwtToken, err := u.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}
	refreshToken := refreshtoken.Generate()
	if exist, _ := u.refreshTokenRepo.ExistByUserID(ctx, user.ID); exist {
		err = u.refreshTokenRepo.DeleteByUserID(ctx, user.ID)
		if err != nil {
			return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
		}
	}
	err = u.refreshTokenRepo.Create(ctx, models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshtoken.HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// return tokens
	return &models.Tokens{
		Refresh: refreshToken,
		JWT:     jwtToken,
	}, nil
}

func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	// hash refresh token
	refreshToken = refreshtoken.HashRefreshToken(refreshToken)

	// get token
	rToken, err := u.refreshTokenRepo.GetByTokenHash(ctx, refreshToken)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w : %w", errs.ErrInvalidToken, err)
	}

	// check token expiration date
	if rToken.ExpiresAt.UnixMilli() < time.Now().UnixMilli() {
		log.Println(rToken.ExpiresAt.UnixMilli(), time.Now().UnixMilli())
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w", errs.ErrInvalidToken)
	}

	// get user
	user, err := u.userR.GetByID(ctx, rToken.UserID)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w", err)
	}

	// check user for delete
	if !user.DeletedAt.IsZero() {
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w", errs.ErrUserNotFound)
	}

	// create tokens
	tokens := &models.Tokens{}
	tokens.Refresh = refreshtoken.Generate()
	tokens.JWT, err = u.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w", err)
	}
	rToken.TokenHash = refreshtoken.HashRefreshToken(tokens.Refresh)
	rToken.CreatedAt = time.Now()
	rToken.ExpiresAt = rToken.CreatedAt.Add(time.Hour * 24 * 7)

	// update refresh token
	err = u.refreshTokenRepo.Update(ctx, *rToken)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.RefreshTokens: %w", err)
	}

	// send tokens
	return tokens, nil
}

func (u *UserService) Logout(ctx context.Context, tokens models.Tokens) error {
	// hash refresh token
	tokens.Refresh = refreshtoken.HashRefreshToken(tokens.Refresh)

	// look for existense
	if exist, _ := u.refreshTokenRepo.ExistByToken(ctx, tokens.Refresh); !exist {
		return fmt.Errorf("user_service.Logout: %w", errs.ErrInvalidToken)
	}

	// delete token
	err := u.refreshTokenRepo.DeleteByToken(ctx, tokens.Refresh)
	if err != nil {
		return fmt.Errorf("user_service.Logout: %w", err)
	}
	return nil
}

func (u *UserService) Create(ctx context.Context, user models.User) error {
	// validation
	if !user.Validate(true) {
		return fmt.Errorf("user_service.Create: %w", errs.ErrUserNotFound)
	}

	passwordHash, err := password.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = passwordHash

	// creating
	err = u.userR.Add(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Create: %w", err)
	}

	return nil
}

func (u *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	// get all users
	allUsers, err := u.userR.GetAll(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("user_service.GetAll: %w", err)
	}

	return allUsers, nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	// get all users
	user, err := u.userR.GetByID(ctx, id)
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", err)
	}

	// get active users
	if !user.DeletedAt.IsZero() {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", errs.ErrUserNotFound)
	}
	return user, nil
}

func (u *UserService) GetProfile(ctx context.Context, id int) (*models.Profile, error) {
	// get all users
	user, err := u.userR.GetByID(ctx, id)
	if err != nil {
		return &models.Profile{}, fmt.Errorf("user_service.GetProfile: %w", err)
	}

	// get active users
	if !user.DeletedAt.IsZero() {
		return &models.Profile{}, fmt.Errorf("user_service.GetProfile: %w", errs.ErrUserNotFound)
	}

	// get orders
	orders, err := u.orderR.GetAllByUserID(ctx, id)
	if err != nil {
		return &models.Profile{}, fmt.Errorf("user_service.GetProfile: %w", err)
	}
	// return profile
	prof := models.NewProfile()
	prof.User = user
	prof.UserOrders = orders
	return prof, nil

}

func (u *UserService) Update(ctx context.Context, user models.User) error {
	// validation
	user.Name = strings.TrimSpace(user.Name)
	user.Phone = strings.TrimSpace(user.Phone)
	if user.Name == "" || user.Phone == "" {
		return fmt.Errorf("user_service.Update: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.userR.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	return nil
}

func (u *UserService) UpdatePassword(ctx context.Context, pass models.UpdatePassword) error {
	// get all users
	user, err := u.userR.GetByID(ctx, pass.UserID)
	if err != nil {
		return fmt.Errorf("user_service.UpdatePassword: %w", err)
	}

	// check old password
	err = password.Compare(user.Password, pass.OldPassword)
	if err != nil {
		return fmt.Errorf("user_service.UpdatePassword: %w", err)
	}

	// hash password
	user.Password, err = password.Hash(pass.NewPassword)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	// update password
	err = u.userR.UpdatePassword(ctx, *user)
	if err != nil {
		return fmt.Errorf("user_service.UpdatePassword: %w", err)
	}
	return nil
}

func (u *UserService) UpdateUserRole(ctx context.Context, user models.User) error {
	// validation
	user.Role = strings.TrimSpace(user.Role)
	if user.Role == "" {
		return fmt.Errorf("user_service.UpdateUserRole: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.userR.UpdateUserRole(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.UpdateUserRole: %w", err)
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int) error {
	// delete user
	err := u.userR.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user_service.DeleteUser: %w", err)
	}

	// delete user orders
	err = u.orderR.DeleteByUserID(ctx, id)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("user_service.DeleteUser: %w", err)
	}

	return nil
}
