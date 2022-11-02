package auth

import (
	"awesome-api/api/common"
	apierror "awesome-api/api/error"
	"awesome-api/api/response"
	"awesome-api/jwt"
	"awesome-api/store"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	TokenName string    `json:"token_name"`
	TokenType string    `json:"token_type"`
	Token     string    `json:"token"`
	ExpireAt  time.Time `json:"expire_at"`
	Scheme    string    `json:"scheme"`
}

type LoginResponse struct {
	Tokens []Token `json:"token"`
}

const (
	accessTokenName  = "access_token"
	accessTokenType  = "access"
	refreshTokenName = "refresh_token"
	refreshTokenType = "refresh"
)

func (sr *signInRequest) validateRequest() *apierror.UnprocessableEntity {
	var err error
	if err = ValidateEmail(sr.Email); err != nil {
		field := apierror.InvalidField{
			Name:    "email",
			Message: err.Error(),
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	if sr.Password == "" {
		field := apierror.InvalidField{
			Name:    "password",
			Message: "password cannot be empty",
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	return nil
}

func Signin(
	zlog zerolog.Logger,
	userStore store.UserStore,
	token jwt.JWT,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := signInRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apierror.ClientBadRequest())
			return
		}
		if fieldErr := req.validateRequest(); fieldErr != nil {
			response.ValidationError(w, *fieldErr)
		}
		ctx := r.Context()
		wlog := common.WrapperZlog{Logger: &zlog}
		usr, err := userStore.FindOneCredentialByEmail(ctx, req.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Error(w, apierror.ClientUnauthorized())
				return
			}
			err = fmt.Errorf("userStore.FindOneCredentialByEmail: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to find one credential by email")
			response.Error(w, apierror.ServerError())
			return
		} else if !usr.IsVerified {
			response.Error(w, apierror.ClientInactiveUser())
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password.String), []byte(req.Password))
		if err != nil {
			response.Error(w, apierror.ClientInvalidCredential())
			return
		}
		claim := jwt.Claim{
			UserId: usr.ID,
		}
		tokenId := make([]byte, 12)
		if _, err = rand.Read(tokenId); err != nil {
			wlog.Error(ctx).
				Err(err).Msg("failed to generate tokenId")
			response.Error(w, apierror.ServerError())
			return
		}
		jti := base64.RawStdEncoding.EncodeToString(tokenId)
		if err = userStore.UpdateTokenIdById(ctx, jti, usr.ID); err != nil {
			err = fmt.Errorf("userStore.UpdateTokenIdById: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to update token_id by id")
			response.Error(w, apierror.ServerError())
			return
		}
		accessToken, err := token.CreateAccessToken(claim)
		if err != nil {
			err = fmt.Errorf("token.CreateAccessToken(claim): %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to generate access token")
			response.Error(w, apierror.ServerError())
			return
		}
		claim.TokenId = jti
		refreshToken, err := token.CreateRefreshToken(claim)
		if err != nil {
			err = fmt.Errorf("token.CreateRefreshToken: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to generate refresh token")
			response.Error(w, apierror.ServerError())
			return
		}
		token := []Token{
			{
				TokenName: accessTokenName,
				TokenType: accessTokenType,
				Token:     accessToken.Token,
				ExpireAt:  accessToken.ExpireAt,
				Scheme:    accessToken.Scheme,
			},
			{
				TokenName: refreshTokenName,
				TokenType: refreshTokenType,
				Token:     refreshToken.Token,
				ExpireAt:  refreshToken.ExpireAt,
				Scheme:    refreshToken.Scheme,
			},
		}
		res := LoginResponse{
			Tokens: token,
		}
		response.GenerateResponse(w, http.StatusOK, res)
	}
}
