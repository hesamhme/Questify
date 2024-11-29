package user

import (
	"context"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) Create(ctx context.Context, user *User) error {
	// TODO validate national code and email
	err := validateUserRegistration(user)
	if err != nil {
		return err
	}
	// TODO
	// hashedPass, err := HashPassword(user.Password)
	// if err != nil {
	// 	return nil, err
	// }
	// user.SetPassword(hashedPass)

	// lowercase email
	//TOTO
	// user.Email = LowerCaseEmail(user.Email)
	err = o.repo.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (o *Ops) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *Ops) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	// email = LowerCaseEmail(email) TODO
	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// if err := CheckPasswordHash(password, user.Password); err != nil {
	// 	return nil, ErrInvalidAuthentication
	// } TODO

	return user, nil
}

func (o *Ops) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// email = LowerCaseEmail(email) TODO
	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func validateUserRegistration(user *User) error {
	// err := ValidateEmail(user.Email) TODO
	// if err != nil {
	// 	return err
	// }

	// if err := ValidatePasswordWithFeedback(user.Password); err != nil {
	// 	return err
	// } // TODO
	return nil
}
