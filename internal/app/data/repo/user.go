package repo

import (
	"github.com/dstgo/maxwell/ent"
	"github.com/dstgo/maxwell/ent/user"
	"golang.org/x/net/context"
)

func NewUserRepo(client *ent.Client) *UserRepo {
	return &UserRepo{Ent: client}
}

type UserRepo struct {
	Ent *ent.Client
}

// FindByNameOrMail returns a User matching the given name or email
func (u *UserRepo) FindByNameOrMail(ctx context.Context, name string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.Or(
				user.UsernameEQ(name),
				user.EmailEQ(name),
			),
		).Only(ctx)
}

// FindByName returns a user matching the given name
func (u *UserRepo) FindByName(ctx context.Context, name string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.UsernameEQ(name),
		).Only(ctx)
}

// FindByEmail returns a User matching the given email
func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.EmailEQ(email),
		).Only(ctx)
}

// CreateNewUser creates a new user with the minimum information
func (u *UserRepo) CreateNewUser(ctx context.Context, username string, email string, password string) (*ent.User, error) {
	return u.Ent.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
}

// UpdateOnePassword updates the user password with specified email
func (u *UserRepo) UpdateOnePassword(ctx context.Context, id int, password string) (*ent.User, error) {
	return u.Ent.User.UpdateOneID(id).
		SetPassword(password).
		Save(ctx)
}
