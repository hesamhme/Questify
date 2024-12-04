package service

import (
	"Questify/internal/question"
	"context"
	"log"

	"Questify/config"
	"Questify/internal/user"
	"Questify/pkg/adapters/storage"
	"Questify/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg           config.Config
	dbConn        *gorm.DB
	userService   *UserService
	authService   *AuthService
	surveyService *SurveyService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage.Migrate(app.dbConn)

	app.setUserService()
	app.setAuthService()
	app.setSurveyService()

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

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db

	err = storage.AddExtension(a.dbConn)
	if err != nil {
		log.Fatal("Create extension failed: ", err)
	}

	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn)), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}

func (a *AppContainer) setSurveyService() {
	if a.surveyService != nil {
		return
	}

	a.surveyService = NewSurveyService(question.NewOps(storage.NewQuestionRepo(a.dbConn)))
}

func (a *AppContainer) SurveyService() *SurveyService {
	return a.surveyService
}


func (s *SurveyService) CreateAnswer(ctx context.Context, answer *question.Answer) error {
	return s.questionOps.CreateAnswer(ctx, answer)
}
