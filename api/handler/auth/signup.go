package auth

import (
	"awesome-api/api/common"
	apierror "awesome-api/api/error"
	"awesome-api/api/response"
	mailer "awesome-api/mail"
	"awesome-api/store"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Fullname          string `json:"fullname"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	IsVerified        bool   `json:"is_verified"`
	TokenVerification string `json:"token_verification"`
	TokenExpiration   string `json:"token_expiration"`
}

const accountStatus = false

func (sr *SignUpRequest) validateRequest() *apierror.UnprocessableEntity {
	var err error
	if err = ValidateEmail(sr.Email); err != nil {
		field := apierror.InvalidField{
			Name:    "email",
			Message: err.Error(),
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	if err = ValidatePassword(sr.Password); err != nil {
		field := apierror.InvalidField{
			Name:    "password",
			Message: err.Error(),
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	if err = ValidateName(sr.Fullname); err != nil {
		field := apierror.InvalidField{
			Name:    "fullname",
			Message: err.Error(),
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	return nil
}

func Signup(
	zlog zerolog.Logger,
	userStore store.UserStore,
	tokenExpiration time.Duration,
	mailer mailer.EmailSender,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := SignUpRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apierror.ClientBadRequest())
			return
		}
		if fieldErr := req.validateRequest(); fieldErr != nil {
			response.ValidationError(w, *fieldErr)
		}
		ctx := r.Context()
		wlog := common.WrapperZlog{Logger: &zlog}
		_, err := userStore.FindOneByEmail(ctx, req.Email)
		if err == nil {
			response.Error(w, apierror.ClientAlreadyExists())
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			wlog.Error(ctx).
				Err(err).Msg("failed to generate bcrypt")
			response.Error(w, apierror.ServerError())
			return
		}
		req.Password = string(hashedPassword)
		req.IsVerified = accountStatus
		req.TokenVerification, err = RandString(128)
		if err != nil {
			wlog.Error(ctx).
				Err(err).Msg("failed to generate random string")
			response.Error(w, apierror.ServerError())
			return
		}
		req.TokenExpiration = strconv.Itoa(int(time.Now().Add(tokenExpiration * time.Minute).Unix()))
		err = registerNewUser(ctx, userStore, req)
		if err != nil {
			wlog.Error(ctx).
				Err(err).Msg("failed to insert new user")
			response.Error(w, apierror.ServerError())
			return
		}
		user, err := userStore.FindOneByEmail(ctx, req.Email)
		if err != nil {
			err = fmt.Errorf("userStore.FindOneByEmail: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to find one by email")
			response.Error(w, apierror.ServerError())
			return
		}
		go mailer.SendActivationLink(user.ID, req.Email, user.TokenVerification.String)
		res := AuthResponse{
			Message: "email activation has been sent, please check your email",
		}
		response.GenerateResponse(w, http.StatusCreated, res)
	}
}

func registerNewUser(
	ctx context.Context,
	userStore store.UserStore,
	user SignUpRequest,
) error {
	usr := &store.UserRegister{
		Email:             user.Email,
		Password:          user.Password,
		Fullname:          user.Fullname,
		IsVerified:        user.IsVerified,
		TokenVerification: user.TokenVerification,
		TokenExpiration:   user.TokenExpiration,
	}
	err := userStore.Insert(ctx, usr)
	if err != nil {
		return fmt.Errorf("userStore.Insert: %w", err)
	}
	return nil
}
