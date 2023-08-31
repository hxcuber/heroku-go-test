package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/user"
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
	"github.com/hxcuber/friends-management/api/internal/handler/response/error_with_string"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateUserByEmail(t *testing.T) {
	type userRepoCreateUser struct {
		called bool
		err    error
	}
	type test struct {
		email          string
		createUser     userRepoCreateUser
		expStatus      int
		expSuccess     bool
		expErrorString string
	}
	for s, tc := range map[string]test{
		"fail_bad_request": {
			email: "test",
			createUser: userRepoCreateUser{
				called: false,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "email field must be an email",
		},
		"fail_already_created": {
			email: "test@test.com",
			createUser: userRepoCreateUser{
				called: true,
				err:    user.ErrAlreadyCreated,
			},
			expStatus:      http.StatusConflict,
			expSuccess:     false,
			expErrorString: user.ErrAlreadyCreated.Error(),
		},
		"fail_unknown": {
			email: "test@test.com",
			createUser: userRepoCreateUser{
				called: true,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
		"success": {
			email: "test@test.com",
			createUser: userRepoCreateUser{
				called: true,
				err:    nil,
			},
			expStatus:  http.StatusCreated,
			expSuccess: true,
		},
	} {
		t.Run(s, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/friends",
				bytes.NewReader([]byte(fmt.Sprintf("{\"email\":\"%s\"}", tc.email))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userCtrl := user.NewMockController(t)

			if tc.createUser.called {
				userCtrl.On("CreateUserByEmail", mock.Anything, tc.email).Return(
					func(ctx context.Context, email string) error {
						return tc.createUser.err
					})
			}

			New(userCtrl).CreateUserByEmail().ServeHTTP(w, r)

			require.Equal(t, tc.expStatus, w.Result().StatusCode)

			respBody, _ := io.ReadAll(w.Result().Body)
			if !tc.expSuccess {
				var response error_with_string.Response
				json.Unmarshal(respBody, &response)
				require.False(t, response.Success)
				require.Equal(t, tc.expErrorString, response.ErrMessage)
			} else {
				var response basic_success.Response
				json.Unmarshal(respBody, &response)
				require.True(t, response.Success)
			}
		})
	}
}
