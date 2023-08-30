package relationship

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/handler/response/errorWithString"
	"github.com/hxcuber/friends-management/api/internal/handler/response/listWithCount"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Friends(t *testing.T) {
	type relaCtrlGetFriends struct {
		called bool
		out    []string
		err    error
	}
	type test struct {
		email          string
		getFriends     relaCtrlGetFriends
		expStatus      int
		expSuccess     bool
		expErrorString string
	}

	for s, tc := range map[string]test{
		"fail_bad_request": {
			email: "test",
			getFriends: relaCtrlGetFriends{
				called: false,
				out:    nil,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "email field must be an email",
		},
		"fail_email_not_found": {
			email: "test@test.com",
			getFriends: relaCtrlGetFriends{
				called: true,
				out:    nil,
				err:    user.ErrEmailNotFound,
			},
			expStatus:      http.StatusNotFound,
			expSuccess:     false,
			expErrorString: user.ErrEmailNotFound.Error(),
		},
		"fail_unknown": {
			email: "test@test.com",
			getFriends: relaCtrlGetFriends{
				called: true,
				out:    nil,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
		"success_empty": {
			email: "test@test.com",
			getFriends: relaCtrlGetFriends{
				called: true,
				out:    []string{},
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
		},
		"success_non_empty": {
			email: "test@test.com",
			getFriends: relaCtrlGetFriends{
				called: true,
				out:    []string{"f1@test.com", "f2@test.com"},
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
		},
	} {
		t.Run(s, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/friends",
				bytes.NewReader([]byte(fmt.Sprintf("{\"email\":\"%s\"}", tc.email))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			relaCtrl := relationship.NewMockController(t)

			if tc.getFriends.called {
				relaCtrl.On("Friends", mock.Anything, tc.email).Return(
					func(ctx context.Context, email string) ([]string, error) {
						return tc.getFriends.out, tc.getFriends.err
					})

			}
			New(relaCtrl).Friends().ServeHTTP(w, r)

			require.Equal(t, tc.expStatus, w.Result().StatusCode)

			respBody, _ := io.ReadAll(w.Result().Body)
			if !tc.expSuccess {
				var response errorWithString.Response
				json.Unmarshal(respBody, &response)
				require.False(t, response.Success)
				require.Equal(t, tc.expErrorString, response.ErrMessage)
			} else {
				var response listWithCount.Response
				json.Unmarshal(respBody, &response)
				require.True(t, response.Success)
				require.ElementsMatch(t, tc.getFriends.out, response.Friends)
				require.Equal(t, len(tc.getFriends.out), response.Count)
			}
		})
	}
}
