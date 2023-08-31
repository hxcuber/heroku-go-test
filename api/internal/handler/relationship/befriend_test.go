package relationship

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/handler/response/basic_success"
	"github.com/hxcuber/friends-management/api/internal/handler/response/error_with_string"
	"github.com/hxcuber/friends-management/api/internal/repository/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Befriend(t *testing.T) {
	type relaCtrlBefriend struct {
		called bool
		err    error
	}
	type test struct {
		emails         []string
		befriend       relaCtrlBefriend
		expStatus      int
		expSuccess     bool
		expErrorString string
	}

	for s, tc := range map[string]test{
		"fail_bad_request": {
			emails: []string{},
			befriend: relaCtrlBefriend{
				called: false,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "friends is a required field",
		},
		"fail_email_not_found": {
			emails: []string{"test1@test.com", "test2@test.com"},
			befriend: relaCtrlBefriend{
				called: true,
				err:    user.ErrEmailNotFound,
			},
			expStatus:      http.StatusNotFound,
			expSuccess:     false,
			expErrorString: user.ErrEmailNotFound.Error(),
		},
		"fail_already_created": {
			emails: []string{"test1@test.com", "test2@test.com"},
			befriend: relaCtrlBefriend{
				called: true,
				err:    relationship.ErrAlreadyCreated,
			},
			expStatus:      http.StatusConflict,
			expSuccess:     false,
			expErrorString: errors.Wrap(relationship.ErrAlreadyCreated, "friendship").Error(),
		},
		"fail_blocked": {
			emails: []string{"test1@test.com", "test2@test.com"},
			befriend: relaCtrlBefriend{
				called: true,
				err:    relationship.ErrBlocked,
			},
			expStatus:      http.StatusConflict,
			expSuccess:     false,
			expErrorString: errors.Wrap(relationship.ErrBlocked, "friendship").Error(),
		},
		"fail_unknown": {
			emails: []string{"test1@test.com", "test2@test.com"},
			befriend: relaCtrlBefriend{
				called: true,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
		"success": {
			emails: []string{"test1@test.com", "test2@test.com"},
			befriend: relaCtrlBefriend{
				called: true,
				err:    nil,
			},
			expStatus:      http.StatusCreated,
			expSuccess:     true,
			expErrorString: "",
		},
	} {
		t.Run(s, func(t *testing.T) {
			var stringBuilder strings.Builder
			stringBuilder.WriteString("[")
			for i, e := range tc.emails {
				if i == 0 {
					stringBuilder.WriteString(fmt.Sprintf("\"%s\"", e))
				} else {
					stringBuilder.WriteString(fmt.Sprintf(", \"%s\"", e))
				}
			}
			stringBuilder.WriteString("]")

			r := httptest.NewRequest(http.MethodGet, "/befriend",
				bytes.NewReader([]byte(fmt.Sprintf("{\"friends\":%s}", stringBuilder.String()))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			relaCtrl := relationship.NewMockController(t)

			if tc.befriend.called {
				relaCtrl.On("Befriend", mock.Anything, tc.emails[0], tc.emails[1]).Return(
					func(ctx context.Context, email1 string, email2 string) error {
						return tc.befriend.err
					})
			}

			New(relaCtrl).Befriend().ServeHTTP(w, r)

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
