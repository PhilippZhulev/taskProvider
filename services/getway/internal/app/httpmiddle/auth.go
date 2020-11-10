package httpmiddle

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
	"gitlab.com/taskProvider/services/getway/internal/app/helpers"
)

var (
	errNoSession = errors.New("session does not exist")
)

// Service ...
type Service struct {
	respond helpers.Respond
}

// Authenticator ...
// Проверка jwt токена
// Промежуточное программное обеспеченеи
func (s Service) Authenticator(sesStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получить токен
			token, cl, err := jwtauth.FromContext(r.Context())
			if err != nil {
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Получить сессию
			session, _ := sesStore.Get(r, cl["uuid"].(string))
			// Если какая либо ошибка
			if err != nil {
				s.respond.ClearSession(session, w, r)
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Если ошибка валидации
			if token == nil || !token.Valid {
				s.respond.ClearSession(session, w, r)
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Проверить сессию
			uuid := session.Values["uuid"]
			if uuid == nil {
				s.respond.Error(w, r, http.StatusUnauthorized, errNoSession)
				return
			}
			// Продолжить обработку
			next.ServeHTTP(w, r)
		})
	}
}
