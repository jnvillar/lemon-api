package usersrepository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	usermodel "lemonapp/domain/user/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	dummyUser := usermodel.NewUser("juan", "noli", "alias", "email")
	expectedQuery := regexp.QuoteMeta(
		`INSERT INTO user (alias,date_created,email,firstname,id,lastname) 
            VALUES (?,?,?,?,?,?)`)
	expectedArgs := []driver.Value{
		dummyUser.Alias,
		dummyUser.DateCreated,
		dummyUser.Email,
		dummyUser.FirstName,
		dummyUser.ID,
		dummyUser.LastName,
	}
	cases := []struct {
		name             string
		initializeDBMock func(dbMock sqlmock.Sqlmock)
		user             *usermodel.User
		expectedUser     *usermodel.User
		expectedError    error
	}{
		{
			name: "Should return an internal error when database returns an  error",
			initializeDBMock: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectExec(expectedQuery).WithArgs(expectedArgs...).WillReturnError(fmt.Errorf("database error"))
			},
			user:          dummyUser,
			expectedUser:  nil,
			expectedError: fmt.Errorf("error creating user: database error"),
		},
		{
			name: "ok",
			initializeDBMock: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectExec(expectedQuery).WithArgs(expectedArgs...).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user:          dummyUser,
			expectedUser:  dummyUser,
			expectedError: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			convey.Convey("Given a repository", t, func() {
				db, dbMock, _ := sqlmock.New()
				c.initializeDBMock(dbMock)
				repository := NewRepository(db)
				convey.Convey("When we create an user", func() {
					ctx := context.Background()
					user, err := repository.Create(ctx, c.user)
					convey.Convey("Then we have to get the user and the error", func() {
						convey.So(user, assertions.ShouldResemble, c.expectedUser)
						convey.So(err, assertions.ShouldResemble, c.expectedError)
						convey.So(dbMock.ExpectationsWereMet(), assertions.ShouldBeNil)
					})
				})
			})
		})
	}
}
