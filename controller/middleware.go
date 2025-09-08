package controller

import (
    "context"
    "net/http"

    "database-example/util"
    "github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
            return
        }

        tokenString := authHeader[len("Bearer "):] // ukloni "Bearer " prefiks
        claims := &util.Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            return util.JwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
            return
        }

        // Spremi role i id u context
        ctx := context.WithValue(r.Context(), "role", claims.Role)
        ctx = context.WithValue(ctx, "userID", claims.ID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

        next.ServeHTTP(w, r)
    })
}