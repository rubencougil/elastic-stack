package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Save(user *User) error
}

type UserStore struct {
	Db     *sqlx.DB
	Logger *logrus.Logger
}

func NewUserStore(db *sqlx.DB, logger *logrus.Logger) Store {
	return &UserStore{
		Db:     db,
		Logger: logger,
	}
}

func (s UserStore) Save(user *User) error {

	query := `INSERT INTO users (name, email) VALUES (?, ?)`
	_, err := s.Db.Exec(query, user.Name, user.Email)
	if err != nil {
		return errors.Wrap(err, "error inserting new user into database")
	}

	s.Logger.WithFields(
		logrus.Fields{
			"name":  user.Name,
			"email": user.Email,
		}).Debug("new user saved into the database")

	return nil
}
