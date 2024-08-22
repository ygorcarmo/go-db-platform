package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygorcarmo/db-platform/views/home"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) error {
	ctx := context.WithValue(r.Context(), "isLoggedIn", true)
	_, err := s.store.GetUserById("bla")
	if err != nil {
		fmt.Println("User not found")
	}
	// fmt.Println("user found: ", user.Id)

	return home.Index().Render(ctx, w)
	// return Render(w, ctx, home.Index())
}
