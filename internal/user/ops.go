package user

import (
	"Questify/pkg/smtp"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
	smtp *smtp.SMTPClient
}

func NewOps(repo Repo, smtpClient *smtp.SMTPClient) *Ops {
	return &Ops{
		repo: repo,
		smtp: smtpClient,
	}
}

func (o *Ops) Create(ctx context.Context, user *User) error {
	// Validate National Code
	if !IsValidIranianNationalCode(user.NationalCode) {
		return ErrInvalidNationalCode
	}

	// Check for Duplicate National Code
	existingUserByNationalCode, err := o.repo.GetByNationalCode(ctx, user.NationalCode)
	if err != nil {
		return fmt.Errorf("failed to check existing national code: %w", err)
	}
	if existingUserByNationalCode != nil {
		return ErrDuplicatedUserNID
	}

	// Validate Password Strength
	if err := ValidatePasswordWithFeedback(user.Password); err != nil {
		return err
	}

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

	// Check for Duplicate Email
	existingUserByEmail, err := o.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUserByEmail != nil {
		return ErrEmailAlreadyExists
	}

	// Save User to Database
	err = o.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (o *Ops) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {

	user, err :=  o.repo.GetByID(ctx, id)
	if err!=nil{
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (o *Ops) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	email = LowerCaseEmail(email)

	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := CheckPasswordHash(password, user.Password); err != nil {
		return nil, ErrInvalidAuthentication
	}

	return user, nil
}

func (o *Ops) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	email = LowerCaseEmail(email)

	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (o *Ops) UpdateUser(ctx context.Context, user *User) error {
	return o.repo.UpdateUser(ctx, user)
}

func (o *Ops) Send2FACodeEmail(ctx context.Context, email, code string) error {
	subject := "Your 2FA Code"
	body := fmt.Sprintf("Your 2FA code is: %s. It will expire in 2 minutes.", code)

	return o.smtp.SendEmail(email, subject, body)
}

func (o *Ops) GetUsers(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, ErrInvalidPagination
	}

	users, totalCount, err := o.repo.GetUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, ErrFailedToGetUsers
	}

	return users, totalCount, nil
}

func (o *Ops) GetUserByNationalCode(ctx context.Context, nationalCode string) (*User, error) {
	// Ensure the National Code is valid
	if !IsValidIranianNationalCode(nationalCode) {
		return nil, ErrInvalidNationalCode
	}

	// Retrieve the user by National Code
	user, err := o.repo.GetByNationalCode(ctx, nationalCode)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
