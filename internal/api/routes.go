package api

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func (api *Api) Routes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	CSRF := csrf.Protect(
		[]byte(os.Getenv("CSRF_KEY")),
		csrf.Secure(false), // HTTP connection for DEV
	)
	api.Router.Use(CSRF)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrftoken", api.HandleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.HandleSignUpUser)
				r.Post("/login", api.handleLoginUser)

				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})

			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)

					r.Post("/", api.HandleCreateProduct)
				})
			})
		})
	})
}
