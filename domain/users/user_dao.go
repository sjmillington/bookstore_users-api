package users

import (
	"fmt"
	"strings"

	"github.com/sjmillington/bookstore_users-api/datasources/mysql/users_db"
	utils "github.com/sjmillington/bookstore_users-api/utils/date"
	"github.com/sjmillington/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Save() *errors.RestErr {

	user.DateCreated = utils.GetNowString()
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)

	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error while trying to save user %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Error trying to get last insert id: %s", err.Error()))
	}

	user.Id = userId

	return nil
}

func (user *User) Get() *errors.RestErr {

	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to retrieve user %d", user.Id))
	}

	return nil
}
