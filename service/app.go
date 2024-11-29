package service

import (
	"context"
	"log"

	"github.com/hesamhme/Questify/config"
	"github.com/hesamhme/Questify/internal/user"
	"github.com/hesamhme/Questify/pkg/adapters/storage"
	"github.com/hesamhme/Questify/pkg/valuecontext"
	"gorm.io/gorm"
)

type AppContainer struct {
	cfg         config.Config
	dbConn      *gorm.DB
	userService *UserService
	authService *AuthService
	//questionService TODO
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage.Migrate(app.dbConn)

	app.setUserService()
	app.setAuthService()

	return app, nil
}

func (a *AppContainer) RawRBConnection() *gorm.DB {
	return a.dbConn
}

func (a *AppContainer) UserService() *UserService {
	return a.userService
}

func (a *AppContainer) UserServiceFromCtx(ctx context.Context) *UserService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.userService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.userService
	}

	return NewUserService(

		user.NewOps(storage.NewUserRepo(gc)),
	)
}

func (a *AppContainer) AuthService() *AuthService {
	return a.authService
}

func (a *AppContainer) setUserService() {
	if a.userService != nil {
		return
	}
	a.userService = NewUserService(user.NewOps(storage.NewUserRepo(a.dbConn)))
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewMysqlGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn)), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}
