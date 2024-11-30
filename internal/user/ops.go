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
	// Validate National Code
	if !IsValidIranianNationalCode(user.NationalCode) {
		return ErrInvalidNationalCode
	}
	err := validateUserRegistration(user)
	if err != nil {
		return err
	}
	
	//hash password
	// Hash Password
	err = user.SetPassword(user.Password)
	if err != nil {
		return err
	}

	// Validate Email Format
	if !IsValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	// Normalize Email to Lowercase
	user.Email = LowerCaseEmail(user.Email)


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
	// Normalize email to lowercase
	email = LowerCaseEmail(email)

	// Retrieve user by email
	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check password validity
	if err := CheckPasswordHash(password, user.Password); err != nil {
		return nil, ErrInvalidAuthentication
	}

	return user, nil
}


func (o *Ops) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Normalize email to lowercase
	email = LowerCaseEmail(email)

	// Retrieve user by email
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
	// Validate email format
	if !IsValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	// Validate password strength
	if err := ValidatePasswordWithFeedback(user.Password); err != nil {
		return err
	}

	return nil
}