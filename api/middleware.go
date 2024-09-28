package api

import (
	"context"
	"net/http"

	"db-platform/models"

	"db-platform/utils"
)

func (s *Server) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := r.Cookie("token")

		if err != nil || jwtToken.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userId, err := utils.DecodeToken(jwtToken.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/"})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := s.store.GetUserById(userId)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), models.UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) adminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.UserCtx).(*models.AppUser)

		if !user.IsAdmin {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) addHttpHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Security-Policy", "frame-ancestors 'none'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("Referrer-Policy", "no-referrer")
		w.Header().Add("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), cross-origin-isolated=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), keyboard-map=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("Strict-Transport-Security", "max-age=5")
		w.Header().Add("Cache-control", "no-store")

		next.ServeHTTP(w, r)

	})
}
