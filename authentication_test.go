package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zopsmart/smart-quiz/services/user"

	"github.com/zopsmart/smart-quiz/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuthentication(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := user.NewMockRepository(ctrl)

	testCases := []struct {
		description string
		value       string
		userData    *models.Users
		expErr      error
	}{
		{
			description: "Success",
			value:       "0 ecf4f9887ed11eb",
			userData:    &models.Users{EmailID: "abc@gmail.com", IsAdmin: false, LoggedAt: time.Now().Format("2006-01-02T15:04:05.000Z")},
		},
		{
			description: "Failed",
			value:       "0ecf4f9887ed11eb",
			userData:    &models.Users{EmailID: "abc@gmail.com", IsAdmin: false, LoggedAt: time.Now().Format("2006-01-02T15:04:05.000Z")},
		},
		{
			description: "Failed : User Details are not present",
			value:       "0 ecf4f9887ed11eb",
			expErr:      errors.New("Random error from store"),
		},
		{
			description: "Failed : TimeOut",
			value:       "0 ecf4f9887ed11eb",
			userData:    &models.Users{EmailID: "abc@gmail.com", IsAdmin: false, LoggedAt: ""},
		},
	}
	for testCase, tc := range testCases {
		authenticateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase == 0 {
				val := r.Context().Value("userEmail")
				if val == nil {
					t.Error("userEmail not present")
				}
				val = r.Context().Value("isAdmin")
				if val == nil {
					t.Error("isAdmin not present")
				}
			}
		})

		if testCase != 1 {
			uuid, _ := uuid.FromBytes([]byte(tc.value[2:]))
			mockStore.EXPECT().GetByUUID(uuid).Return(tc.userData, tc.expErr)
		}
		auth := Authentication(mockStore)

		req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
		req.Header.Add("Authorization", tc.value)
		if testCase == 0 {
			req = req.WithContext(context.WithValue(req.Context(), "userEmail", tc.userData.EmailID))
			req = req.WithContext(context.WithValue(req.Context(), "isAdmin", tc.userData.IsAdmin))
		}
		auth(authenticateHandler).ServeHTTP(httptest.NewRecorder(), req)

	}
}
