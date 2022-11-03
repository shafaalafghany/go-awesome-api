package auth

import (
	"awesome-api/api/common"
	apierror "awesome-api/api/error"
	"awesome-api/api/response"
	"awesome-api/jwt"
	"awesome-api/store"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type LogoutRequest struct {
	Token string `json:"token"`
}

func (lr *LogoutRequest) validateRequest() *apierror.UnprocessableEntity {
	var err error
	if err = ValidateToken(lr.Token); err != nil {
		field := apierror.InvalidField{
			Name:    "token",
			Message: err.Error(),
		}
		fieldErr := apierror.ClientInvalidField(field)
		return &fieldErr
	}
	return nil
}

func Signout(
	zlog zerolog.Logger,
	userStore store.UserStore,
	token jwt.JWT,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := LogoutRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apierror.ClientBadRequest())
			return
		}
		if fieldErr := req.validateRequest(); fieldErr != nil {
			response.ValidationError(w, *fieldErr)
			return
		}
		ctx := r.Context()
		wlog := common.WrapperZlog{Logger: &zlog}
		claim, err := token.ExpectRefreshToken(req.Token)
		if err != nil {
			err = fmt.Errorf("token.ExpectRefreshToken: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to expect refresh token")
			response.Error(w, apierror.ServerError())
			return
		}
		user, err := userStore.FindOneById(ctx, claim.UserId)
		if err != nil {
			err = fmt.Errorf("userStore.FindOneById: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to find on by id")
			response.Error(w, apierror.ServerError())
			return
		}
		if user.TokenID.String != claim.TokenId {
			response.Error(w, apierror.ClientUnauthorized())
			return
		}
		err = userStore.DeleteTokenIdById(ctx, user.ID)
		if err != nil {
			err = fmt.Errorf("userStore.DeleteTokenIdById: %w", err)
			wlog.Error(ctx).
				Err(err).Msg("failed to delete token_id by id")
			response.Error(w, apierror.ServerError())
			return
		}
		res := AuthResponse{
			Message: "user has been signout",
		}
		response.GenerateResponse(w, http.StatusOK, res)
	}
}
