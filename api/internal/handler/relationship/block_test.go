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
	"testing"
)

func TestHandler_Block(t *testing.T) {
	type relaCtrlBlock struct {
		called bool
		err    error
	}
	type test struct {
		requestorEmail string
		targetEmail    string
		block          relaCtrlBlock
		expStatus      int
		expSuccess     bool
		expErrorString string
	}

	for s, tc := range map[string]test{
		"fail_bad_request": {
			requestorEmail: "",
			targetEmail:    "",
			block: relaCtrlBlock{
				called: false,
				err:    nil,
			},
			expStatus:      http.StatusBadRequest,
			expSuccess:     false,
			expErrorString: "requestor is a required field",
		},
		"fail_email_not_found": {
			requestorEmail: "request@test.com",
			targetEmail:    "target@test.com",
			block: relaCtrlBlock{
				called: true,
				err:    user.ErrEmailNotFound,
			},
			expStatus:      http.StatusNotFound,
			expSuccess:     false,
			expErrorString: user.ErrEmailNotFound.Error(),
		},
		"fail_already_created": {
			requestorEmail: "request@test.com",
			targetEmail:    "target@test.com",
			block: relaCtrlBlock{
				called: true,
				err:    relationship.ErrAlreadyCreated,
			},
			expStatus:      http.StatusConflict,
			expSuccess:     false,
			expErrorString: errors.Wrap(relationship.ErrAlreadyCreated, "block").Error(),
		},
		"fail_blocked": {
			requestorEmail: "request@test.com",
			targetEmail:    "target@test.com",
			block: relaCtrlBlock{
				called: true,
				err:    relationship.ErrBlocked,
			},
			expStatus:      http.StatusConflict,
			expSuccess:     false,
			expErrorString: errors.Wrap(relationship.ErrBlocked, "block").Error(),
		},
		"fail_unknown": {
			requestorEmail: "request@test.com",
			targetEmail:    "target@test.com",
			block: relaCtrlBlock{
				called: true,
				err:    errors.New("unknown"),
			},
			expStatus:      http.StatusInternalServerError,
			expSuccess:     false,
			expErrorString: "Something went wrong",
		},
		"success": {
			requestorEmail: "request@test.com",
			targetEmail:    "target@test.com",
			block: relaCtrlBlock{
				called: true,
				err:    nil,
			},
			expStatus:      http.StatusOK,
			expSuccess:     true,
			expErrorString: "",
		},
	} {
		t.Run(s, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/befriend",
				bytes.NewReader([]byte(fmt.Sprintf("{\"requestor\":\"%s\", \"target\":\"%s\"}", tc.requestorEmail, tc.targetEmail))))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			relaCtrl := relationship.NewMockController(t)

			if tc.block.called {
				relaCtrl.On("Block", mock.Anything, tc.requestorEmail, tc.targetEmail).Return(
					func(ctx context.Context, requestorEmail string, targetEmail string) error {
						return tc.block.err
					})
			}

			New(relaCtrl).Block().ServeHTTP(w, r)

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
