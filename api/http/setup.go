package http

import (
	"Questify/api/http/handlers"
	middlewares "Questify/api/http/middlerwares"
	"Questify/config"
	"Questify/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()

	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	// register global routes
	registerGlobalRoutes(api, app)
	secret := []byte(cfg.Server.TokenSecret)
	fmt.Println(secret)
	//registerQuestionRoutes(api, app, secret, createGroupLogger("boards"))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer) {
	//router.Use(loggerMiddleWare)
	router.Post("/register", handlers.Register(app.AuthService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/test-email", handlers.SendTestEmail(app))
	router.Get("/refresh", handlers.RefreshToken(app.AuthService()))
}

// func userRoleChecker() fiber.Handler {
// 	return middlewares.RoleChecker("user")
// }

// func registerBoardRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
// 	router = router.Group("/boards")
// 	router.Use(loggerMiddleWare)

// 	router.Post("",
// 		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
// 		middlewares.Auth(secret),
// 		handlers.CreateUserBoard(app.BoardServiceFromCtx),
// 	)
// 	router.Get("/my-boards",
// 		middlewares.Auth(secret),
// 		handlers.GetUserBoards(app.BoardService()),
// 	)
// 	router.Get("/publics",
// 		middlewares.Auth(secret),
// 		middlewares.SetupCacheMiddleware(5),
// 		handlers.GetPublicBoards(app.BoardService()),
// 	)
// 	router.Get("/:boardID",
// 		middlewares.Auth(secret),
// 		handlers.GetFullBoardByID(app.BoardService()),
// 	)

// 	router.Delete("/:boardID",
// 		middlewares.Auth(secret),
// 		handlers.DeleteBoard(app.BoardService()),
// 	)

// 	router.Post("/invite", middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
// 		middlewares.Auth(secret),
// 		handlers.InviteToBoard(app.BoardServiceFromCtx))
// }
