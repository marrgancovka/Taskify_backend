package middleware

import (
	"TaskTracker/internal/pkg/constans"
	"TaskTracker/internal/pkg/tokenizer"
	"context"
	"go.uber.org/fx"
	"log"
	"net/http"
)

type AuthMiddlewareParams struct {
	fx.In

	Tokenizer *tokenizer.Tokenizer
}

type AuthMiddleware struct {
	tokenizer *tokenizer.Tokenizer
}

func NewAuthMiddleware(tokenizer *tokenizer.Tokenizer) *AuthMiddleware {
	return &AuthMiddleware{tokenizer: tokenizer}
}

func (md *AuthMiddleware) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(constans.CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		tokenPayload, err := md.tokenizer.ValidateJWT(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constans.ContextValue, tokenPayload.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
