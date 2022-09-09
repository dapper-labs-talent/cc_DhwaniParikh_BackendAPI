package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/internal/api/appuser"
	"server/internal/api/auth"

	"github.com/go-chi/chi"
	jwtauth "server/internal/auth"
	"server/internal/config"
	"server/internal/database"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/router"
	il "server/internal/util/logger"
)

var tokenAuth *jwtauth.JWTAuth

// @title Dapper Labs API
// @version 1.0
// @host localhost:9090
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {

	log.Info().Msg("Initializing server")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load server configuration")
		panic(err)
	}

	logger := il.New(c.Log.ZeroLogLevel)

	// Set up JWT authentication
	tokenAuth = jwtauth.New("HS256", []byte(c.Auth.JWT.Key), nil)

	// If we get here then we have loaded the configuration and can continue
	r := router.New(c, logger.ZeroLog())

	zl := logger.ZeroLog()
	d := database.Connect(c, zl)
	// Set up the services, repositories and controllers

	appUserRepo := repositories.NewAppUserRepo(d, zl)
	tokenService := jwtauth.NewTokenService(c.Auth, zl)

	authController := auth.NewAuthController(c, appUserRepo, tokenService, zl)
	appUserController := appuser.NewAppUserController(appUserRepo)

	// JWT protected routes
	r.Group(func(protected chi.Router) {
		// Seek, verify and validate JWT tokens
		protected.Use(jwtauth.Verifier(tokenAuth))
		protected.Use(jwtauth.Authenticator)
		protected.Use(jwtauth.Verifier(tokenAuth))
		router.AddControllersToGroup("/api/users", protected, []models.Controller{appUserController})
	})

	// Public routes
	r.Group(func(public chi.Router) {
		router.AddControllersToGroup("/auth", public, []models.Controller{authController})
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), r)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
