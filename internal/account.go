package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/KonstantinGasser/required"
	"github.com/KonstantinGasser/sherlock/security"
)

var (
	ErrInsecurePassword   = fmt.Errorf("provided password is insecure (use --insecure to ignore this message)")
	ErrInvalidAccountName = fmt.Errorf("account name must be a consecutive string")
	ErrMissingValues      = fmt.Errorf("account is missing required values")
)

type Account struct {
	Name      string    `json:"name" required:"yes"`
	Password  string    `json:"password" required:"yes"`
	Tag       string    `json:"tag"`
	CreatedOn time.Time `json:"created_on" required:"yes"`
	UpdatedOn time.Time `json:"updated_on"`
}

// NewAccount creates a new Account and if insecure=false checks the password strength
// returning an err if strength security.Low
func NewAccount(name, password, tag string, insecure bool) (*Account, error) {
	a := Account{
		Name:      name,
		Password:  password,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Tag:       tag,
	}
	if err := a.valid(); err != nil {
		return nil, err
	}

	if insecure {
		return &a, nil
	}
	if level := a.secure(); level == security.Low {
		return nil, ErrInsecurePassword
	}
	return &a, nil
}

func (a Account) valid() error {
	if err := required.Atomic(&a); err != nil {
		return ErrMissingValues
	}
	if set := strings.Split(a.Name, " "); len(set) > 1 {
		return ErrInvalidAccountName
	}
	return nil
}

// secure checks the Accounts on how secure it is
func (a Account) secure() int {
	return security.PasswordStrength(a.Password)
}
