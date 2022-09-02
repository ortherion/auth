package rest

import (
	"auth/internal/ports"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(logger *logrus.Logger, service ports.AuthService, help *helpers) *chi.Mux {
	mux := chi.NewMux()
	mux.Route("/", func(r chi.Router) {
		r.Use(Recover(logger))
		r.Use(Tracer)
		r.With(RequestID)
		r.Mount("/auth", AuthRouter(logger, service, help))
		r.With(Validate(help, service))
		r.With(Logger(logger))
	})
	mux.Mount("/debug", PprofHandler())
	return mux
}

func AuthRouter(logger *logrus.Logger, authService ports.AuthService, helpers *helpers) http.Handler {
	handlers := NewHandler(authService, logger, helpers)

	r := chi.NewRouter()
	r.Post("/login", handlers.Login)
	r.Post("/logout", handlers.Logout)
	r.Post("/i", handlers.Info)
	//r.Post("/signup", handlers.SignUp)

	return r
}

func SwaggerRouter(basePath string) http.Handler {
	r := chi.NewRouter()

	httpSwagger.UIConfig(map[string]string{
		"showExtensions":        "true",
		"onComplete":            `() => { window.ui.setBasePath('v3'); }`,
		"defaultModelRendering": `"model"`,
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", basePath))))

	return r
}
