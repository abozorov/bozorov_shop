package errs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ServerError
	ErrSomethingWentWrong = errors.New("something went wrong")

	// BadRequest
	ErrBadRequest        = errors.New("bad request")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrBadRequestBody    = errors.New("bad request body")
	ErrBadRequestQuery   = errors.New("bad request query")
	ErrInvalidUserId     = errors.New("invalid user id")

	// Too Many Requests
	ErrTooManyRequests = errors.New("too many requests")

	// Timeout exceeded
	ErrTimeoutExceeded = errors.New("timeout exceeded")

	// Verify
	ErrVerifyingFailed     = errors.New("Verifying failed")
	ErrUserNotBeenVerified = errors.New("email has not been verified yet")

	// Unautorized
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
	ErrIncorrectPassword        = errors.New("incorrect password")

	// NotFound
	ErrNotFound       = errors.New("not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserIDNotFound = errors.New("user id not Found")
)

// pkg.errs -> http.Error
func ErrsToHttp(w http.ResponseWriter, err error) {

	switch {
	// http.StatusNotFound
	case errors.Is(err, ErrNotFound):
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, ErrUserNotFound):
		http.Error(w, ErrUserNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, ErrUserIDNotFound):
		http.Error(w, ErrUserIDNotFound.Error(), http.StatusNotFound)

	// http.StatusBadRequest
	case errors.Is(err, ErrBadRequest):
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrUserAlreadyExists):
		http.Error(w, ErrUserAlreadyExists.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrBadRequestBody):
		http.Error(w, ErrBadRequestBody.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrBadRequestQuery):
		http.Error(w, ErrBadRequestQuery.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidUserId):
		http.Error(w, ErrInvalidUserId.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrVerifyingFailed):
		http.Error(w, ErrVerifyingFailed.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrUserNotBeenVerified):
		http.Error(w, ErrUserNotBeenVerified.Error(), http.StatusBadRequest)

	// http.StatusTooManyRequests
	case errors.Is(err, ErrTooManyRequests):
		http.Error(w, ErrTooManyRequests.Error(), http.StatusTooManyRequests)

	// http.StatusGatewayTimeout
	case errors.Is(err, ErrTimeoutExceeded):
		http.Error(w, ErrTimeoutExceeded.Error(), http.StatusGatewayTimeout)

	// http.StatusUnauthorized
	case errors.Is(err, ErrIncorrectLoginOrPassword):
		http.Error(w, ErrIncorrectLoginOrPassword.Error(), http.StatusUnauthorized)
	case errors.Is(err, ErrIncorrectPassword):
		http.Error(w, ErrIncorrectPassword.Error(), http.StatusUnauthorized)

	// http.StatusInternalServerError
	default:
		http.Error(w, ErrSomethingWentWrong.Error(), http.StatusInternalServerError)

	}
}

// GORM -> pkg.errs
func PostgresToErrs(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows): // no rows
		return fmt.Errorf("error %w: %w", err, ErrNotFound)

	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("error %w: %w", err, ErrTimeoutExceeded)

	case errors.Is(err, context.Canceled):
		return fmt.Errorf("error %w: %w", err, ErrTimeoutExceeded)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("error %w: %w", err, ErrBadRequest)

		case "23503": // foreign_key_violation
			return fmt.Errorf("error %w: %w", err, ErrBadRequest)

		case "23502": // not_null_violation
			return fmt.Errorf("error %w: %w", err, ErrBadRequestBody)

		case "23514": // check_violation
			return fmt.Errorf("error %w: %w", err, ErrBadRequestBody)

		case "22001": // string_data_right_truncation
			return fmt.Errorf("error %w: %w", err, ErrBadRequestBody)

		case "22003": // numeric_value_out_of_range
			return fmt.Errorf("error %w: %w", err, ErrBadRequestBody)

		case "57014": // query_canceled
			return fmt.Errorf("error %w: %w", err, ErrTimeoutExceeded)

		case "22P02": // invalid_text_representation
			return fmt.Errorf("error %w: %w", err, ErrBadRequest)
		}
	}

	return fmt.Errorf("error %w: %w", err, ErrSomethingWentWrong)
}
