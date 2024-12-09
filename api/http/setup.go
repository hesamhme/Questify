package http

import (
	"Questify/api/http/handlers"
	middlewares "Questify/api/http/middlerwares"
	"Questify/config"
	"Questify/pkg/adapters"
	"Questify/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()

	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	// register global routes
	secret := []byte(cfg.Server.TokenSecret)
	registerGlobalRoutes(api, app)

	registerUserRoutes(api, app, secret)
	//registerQuestionRoutes(api, app, secret, createGroupLogger("boards"))
	// Register survey routes
	registerSurveyRoutes(cfg, api, app)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer) {
	//router.Use(loggerMiddleWare)
	router.Post("/register", handlers.Register(app.AuthService()))
	router.Post("/sign-up",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawRBConnection())),
		handlers.SignUpUser(app.AuthServiceFromCtx),
	)
	router.Post("/confirm-tfa", handlers.ConfirmTFA(app.AuthService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/refresh", handlers.RefreshToken(app.AuthService()))

	// Role and permissions management
	router.Post("/roles", handlers.CreateRole(app.RoleService()))
	router.Get("/roles", handlers.GetAllRoles(app.RoleService()))
	router.Delete("roles/delete", handlers.DeleteRoles(app.RoleService()))
}

func registerUserRoutes(router fiber.Router, app *service.AppContainer, secret []byte) {
	router = router.Group("/users")

	router.Get("",
		middlewares.Auth(secret),
		handlers.GetAllVerifiedUsers(app.UserService()),
	)
}

func registerSurveyRoutes(cfg config.Config, router fiber.Router, app *service.AppContainer) {
	router.Use(middlewares.Auth([]byte(cfg.Server.TokenSecret)))
	router = router.Group("/survey")
	router.Post("", handlers.CreateSurvey(app.SurveyService()))
	router.Post("/:surveyId", handlers.GetSurvey(app.SurveyService()))
	router.Post("/:surveyId/question", handlers.CreateQuestion(app.SurveyService()))
	router.Get("/:surveyId/question/next", handlers.GetNextQuestion(app.SurveyService()))
	router.Get("/:surveyId/question/previous", handlers.GetPreviousQuestion(app.SurveyService()))
	router.Get("/:surveyId/question/:questionId", handlers.GetQuestion(app.SurveyService()))
	router.Put("/:surveyId/question/:questionId", handlers.UpdateQuestion(app.SurveyService()))
	router.Post("/question/:questionId/answer", handlers.CreateAnswer(app.SurveyService()))
	// Survey-specific role management
	router.Post("/:surveyId/roles/assign", middlewares.Auth([]byte(cfg.Server.TokenSecret)), handlers.AssignRoleToSurveyUser(app.RoleService()))
	router.Get("/:surveyId/roles/check-permission", handlers.CheckSurveyPermission(app.RoleService()))
}
