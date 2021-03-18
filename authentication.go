package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/zopsmart/smart-quiz/services/user"
)

const TokenExpiryTime = 30

// Authentication method is used to authenticate the api calls by uuid , if uuid expires user needs to again login
// nolint -  cognitive complexity 13 of func `Authentication` is high (> 10)
func Authentication(repository user.Repository) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestWithEmail *http.Request
			var requestWithEmailAndIsAdmin *http.Request
			authHeader := strings.Split(r.Header.Get("Authorization"), " ")
			if len(authHeader) != 2 {
				inner.ServeHTTP(w, r)
			} else {
				fmt.Println(authHeader[1])
				uid, _ := uuid.FromBytes([]byte(authHeader[1]))
				userDetails, err := repository.GetByUUID(uid)
				if err != nil {
					fmt.Println("Non existing token")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Non existing Token"))
					inner.ServeHTTP(w, r)
					return
				}
				time1, _ := time.Parse("2006-01-02T15:04:05.000Z", userDetails.LoggedAt)
				diff := time.Since(time1).Seconds()
				if (diff / 60) > TokenExpiryTime {
					fmt.Println("Expired token")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Expired Token"))
				} else {
					requestWithEmail = r.WithContext(context.WithValue(r.Context(), "userEmail", userDetails.EmailID))
					requestWithEmailAndIsAdmin = requestWithEmail.WithContext(context.WithValue(requestWithEmail.Context(), "isAdmin", userDetails.IsAdmin))
				}
				inner.ServeHTTP(w, requestWithEmailAndIsAdmin)
			}
		})
	}
}
