package relationship

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/handler/response/errorWithString"
	"github.com/hxcuber/friends-management/api/internal/handler/response/recipients"
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

func TestHandler_Receivers(t *testing.T) {
	type relaCtrlGetReceivers struct {
		called bool
		count  int
		err    error
	}
	type test struct {
		senderEmail    string
		getReceivers   relaCtrlGetReceivers
		expStatus      int
		expSuccess     bool
		expErrorString string
	}

	for s, tc := range map[string]test{
		"success_empty": {
			senderEmail: "test@test.com",
			getReceivers: relaCtrlGetReceivers{
				called: true,
				count:  0,
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
		},
		"success_non_empty": {
			senderEmail: "test@test.com",
			getReceivers: relaCtrlGetReceivers{
				called: true,
				count:  2,
				err:    nil,
			},
			expStatus:  http.StatusOK,
			expSuccess: true,
		},
		// TODO
		"fail_invalid_email_1": {
			senderEmail: "test",
			getReceivers: relaCtrlGetReceivers{
				called: false,
				count:  0,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "sender field must be an email",
		},
		"fail_invalid_email_2": {
			senderEmail: "test@",
			getReceivers: relaCtrlGetReceivers{
				called: false,
				count:  0,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "sender field must be an email",
		},
		"fail_invalid_email_3": {
			senderEmail: "test@test",
			getReceivers: relaCtrlGetReceivers{
				called: false,
				count:  0,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "sender field must be an email",
		},
		"fail_email_not_found": {
			senderEmail: "test@test.com",
			getReceivers: relaCtrlGetReceivers{
				called: true,
				count:  0,
				err:    user.ErrEmailNotFound,
			},
			expStatus:      http.StatusNotFound,
			expSuccess:     false,
			expErrorString: user.ErrEmailNotFound.Error(),
		},
		"fail_unknown": {
			senderEmail: "test@test.com",
			getReceivers: relaCtrlGetReceivers{
				called: true,
				count:  0,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
	} {
		t.Run(s, func(t *testing.T) {
			var stringBuilder strings.Builder
			var receiverEmails []string
			for i := 1; i <= tc.getReceivers.count; i++ {
				receiverEmail := fmt.Sprintf("receiver%d@test.com", i)
				stringBuilder.WriteString(" " + receiverEmail)

				receiverEmails = append(receiverEmails, receiverEmail)
			}

			r := httptest.NewRequest(http.MethodGet, "/friends",
				bytes.NewReader([]byte(
					fmt.Sprintf("{\"sender\":\"%s\", \"text\":\"%s\"}", tc.senderEmail, stringBuilder.String()))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			relaCtrl := relationship.NewMockController(t)

			if tc.getReceivers.called {

				relaCtrl.On("Receivers", mock.Anything, tc.senderEmail, stringBuilder.String()).Return(
					func(ctx context.Context, sender string, text string) ([]string, error) {
						return receiverEmails, tc.getReceivers.err
					})

			}
			New(relaCtrl).Receivers().ServeHTTP(w, r)

			require.Equal(t, tc.expStatus, w.Result().StatusCode)

			respBody, _ := io.ReadAll(w.Result().Body)
			if !tc.expSuccess {
				var response errorWithString.Response
				json.Unmarshal(respBody, &response)
				require.False(t, response.Success)
				require.Equal(t, tc.expErrorString, response.ErrMessage)
			} else {
				var response recipients.Response
				json.Unmarshal(respBody, &response)
				require.True(t, response.Success)
				require.ElementsMatch(t, receiverEmails, response.Recipients)
			}
		})
	}
}
