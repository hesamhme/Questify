package middlewares

import (
	"Questify/pkg/valuecontext"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

func SetUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctxValue := &valuecontext.ContextValue{
			Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		}

		c.SetUserContext(valuecontext.NewValueContext(c.UserContext(), ctxValue))

		return c.Next()
	}
}

func SetTransaction(committer valuecontext.Committer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cm := committer.Begin()
		valuecontext.SetTx(c.UserContext(), cm)

		logger := valuecontext.GetLogger(c.UserContext())

		logger.Info("starting transaction")
		if err := c.Next(); err != nil {
			logger.Info("rollback on error", "error", err.Error())
			cm.Rollback()
			return err
		}

		err, ok := c.Locals(valuecontext.IsTxError).(error)
		if ok && err != nil {
			logger.Info("rollback on not ok response", "error", err.Error())
			cm.Rollback()
			return nil
		}

		if err := cm.Commit(); err != nil {
			logger.Info("commit error", "err", err.Error())
			cm.Rollback()
			return err
		}

		logger.Info("ending transaction")
		return nil
	}
}
