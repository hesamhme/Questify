package service

import (
	"Questify/internal/question"
	"Questify/internal/survey"
	"context"
	"log"

	"Questify/config"
	"Questify/internal/user"
	"Questify/pkg/adapters/storage"
	"Questify/pkg/smtp"
	"Questify/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg           config.Config
	dbConn        *gorm.DB
	userService   *UserService
	authService   *AuthService
	smtpClient    *smtp.SMTPClient
	surveyService *SurveyService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage.Migrate(app.dbConn)

	app.setSMTPClient()
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
		user.NewOps(storage.NewUserRepo(gc), a.smtpClient), // Inject SMTPClient
	)
}

func (a *AppContainer) AuthService() *AuthService {
	return a.authService
}

func (a *AppContainer) setUserService() {
	if a.userService != nil {
		return
	}
	a.userService = NewUserService(user.NewOps(storage.NewUserRepo(a.dbConn), a.smtpClient)) // Inject SMTPClient
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

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn), a.smtpClient), // Inject SMTPClient
		[]byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}

func (a *AppContainer) SurveyService() *SurveyService {
	return a.surveyService
}

func (a *AppContainer) setSurveyService() {
	if a.surveyService != nil {
		return
	}

	a.surveyService = NewSurveyService(
		question.NewOps(storage.NewQuestionRepo(a.dbConn)),
		survey.NewOps(storage.NewSurveyRepo(a.dbConn)),
	)
}

func (a *AppContainer) setSMTPClient() {
	if a.smtpClient != nil {
		return
	}

	a.smtpClient = smtp.NewSMTPClient(a.cfg.SMTP)
}

func (a *AppContainer) SMTPClient() *smtp.SMTPClient {
	return a.smtpClient
}

func (s *SurveyService) CreateAnswer(ctx context.Context, answer *question.Answer) error {
	return s.questionOps.CreateAnswer(ctx, answer)
}
