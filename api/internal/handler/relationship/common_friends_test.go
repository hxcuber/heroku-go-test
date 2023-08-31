package relationship

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/handler/response/error_with_string"
	"github.com/hxcuber/friends-management/api/internal/handler/response/list_with_count"
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

func TestHandler_CommonFriends(t *testing.T) {
	type relaCtrlCommonFriends struct {
		called bool
		out    []string
		err    error
	}
	type test struct {
		emails         []string
		getCommon      relaCtrlCommonFriends
		expStatus      int
		expSuccess     bool
		expErrorString string
	}

	for s, tc := range map[string]test{
		"fail_bad_request": {
			emails: nil,
			getCommon: relaCtrlCommonFriends{
				called: false,
				out:    nil,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "friends is a required field",
		},
		"fail_email_not_found": {
			emails: []string{"test1@test.com", "test2@test.com"},
			getCommon: relaCtrlCommonFriends{
				called: true,
				out:    nil,
				err:    user.ErrEmailNotFound,
			},
			expStatus:      http.StatusNotFound,
			expSuccess:     false,
			expErrorString: user.ErrEmailNotFound.Error(),
		},
		"fail_unknown": {
			emails: []string{"test1@test.com", "test2@test.com"},
			getCommon: relaCtrlCommonFriends{
				called: true,
				out:    nil,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
		"success_empty": {
			emails: []string{"test1@test.com", "test2@test.com"},
			getCommon: relaCtrlCommonFriends{
				called: true,
				out:    []string{},
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
		},
		"success_non_empty": {
			emails: []string{"test1@test.com", "test2@test.com"},
			getCommon: relaCtrlCommonFriends{
				called: true,
				out:    []string{"f1@test.com", "f2@test.com"},
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
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

			r := httptest.NewRequest(http.MethodGet, "/friends",
				bytes.NewReader([]byte(fmt.Sprintf("{\"friends\":%s}", stringBuilder.String()))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			relaCtrl := relationship.NewMockController(t)

			if tc.getCommon.called {
				relaCtrl.On("CommonFriends", mock.Anything, tc.emails[0], tc.emails[1]).Return(
					func(ctx context.Context, email1 string, email2 string) ([]string, error) {
						return tc.getCommon.out, tc.getCommon.err
					})

			}
			New(relaCtrl).CommonFriends().ServeHTTP(w, r)

			require.Equal(t, tc.expStatus, w.Result().StatusCode)

			respBody, _ := io.ReadAll(w.Result().Body)
			if !tc.expSuccess {
				var response error_with_string.Response
				json.Unmarshal(respBody, &response)
				require.False(t, response.Success)
				require.Equal(t, tc.expErrorString, response.ErrMessage)
			} else {
				var response list_with_count.Response
				json.Unmarshal(respBody, &response)
				require.True(t, response.Success)
				require.ElementsMatch(t, tc.getCommon.out, response.Friends)
				require.Equal(t, len(tc.getCommon.out), response.Count)
			}
		})
	}
}
