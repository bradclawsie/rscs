package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/pressly/chi"
	"net/http"
)

// ContextKeyType is the type used to reference values in the Context.
type ContextKeyType int

var contextKey ContextKeyType // == 0

// insertKeyContext places the key into the Context.
func insertKeyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, keyName)
		ctx := context.WithValue(r.Context(), contextKey, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractKeyContext attempts to extract the key as a string from the Context.
func extractKeyContext(r *http.Request) (string, error) {
	key := r.Context().Value(contextKey)
	if key == nil {
		e := fmt.Sprintf("%v not found in Context", contextKey)
		return "", errors.New(e)
	}
	keyStr, ok := key.(string)
	if !ok {
		return "", errors.New("cannot assert to string")
	}
	if keyStr == "" {
		return "", errors.New("empty key string")
	}
	return keyStr, nil
}
