package usersrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	usersmodel "lemonapp/domain/user/model"
	"lemonapp/errors"
	"lemonapp/storage/mysql"

	sq "github.com/Masterminds/squirrel"
	mysql2 "github.com/go-sql-driver/mysql"
)

const duplicateKeyErrorCode = 1062

const (
	userTable         = "user"
	idColumn          = "id"
	firstnameColumn   = "firstname"
	lastnameColumn    = "lastname"
	aliasColumn       = "alias"
	emailColumn       = "email"
	dateCreatedColumn = "date_created"
)

var userColumns = []string{
	idColumn,
	firstnameColumn,
	lastnameColumn,
	aliasColumn,
	emailColumn,
	dateCreatedColumn,
}

type sqlRepository struct {
	sqlBuilder sq.StatementBuilderType
	db         *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db:         db,
		sqlBuilder: sq.StatementBuilder,
	}
}

func (m *sqlRepository) List(ctx context.Context) ([]*usersmodel.User, error) {
	builder := m.sqlBuilder.
		Select(userColumns...).
		From(userTable).
		OrderBy(fmt.Sprintf("%s DESC", dateCreatedColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list users query")
	}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing users %v", err)
	}

	users := make([]*usersmodel.User, 0)
	for rows.Next() {
		event, err := m.scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, event)
	}
	return users, nil
}

func (m *sqlRepository) Create(ctx context.Context, user *usersmodel.User) (*usersmodel.User, error) {
	valuesMap := map[string]interface{}{
		idColumn:          user.ID,
		firstnameColumn:   user.FirstName,
		lastnameColumn:    user.LastName,
		aliasColumn:       user.Alias,
		emailColumn:       user.Email,
		dateCreatedColumn: user.DateCreated,
	}

	query, args, builderErr := m.sqlBuilder.
		Insert(userTable).
		SetMap(valuesMap).
		ToSql()
	if builderErr != nil {
		return nil, fmt.Errorf("error building user query: %v", builderErr)
	}

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		me, ok := err.(*mysql2.MySQLError)
		if !ok {
			return nil, fmt.Errorf("error creating user: %v", err)
		}
		if me.Number == duplicateKeyErrorCode {
			if strings.Contains(me.Message, "user.alias") {
				return nil, errors.NewBadRequestError(fmt.Errorf("alias %s already registered", user.Alias))
			}
			if strings.Contains(me.Message, "user.email") {
				return nil, errors.NewBadRequestError(fmt.Errorf("email %s already registered", user.Email))
			}
		}
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}

func (m *sqlRepository) GetByID(ctx context.Context, ID string) (*usersmodel.User, error) {
	query, args, err := m.sqlBuilder.
		Select(userColumns...).
		From(userTable).
		Where(sq.Eq{idColumn: ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating user query: %v", err)
	}
	row := m.db.QueryRowContext(ctx, query, args...)
	return m.scanUser(row)
}

func (m *sqlRepository) scanUser(row mysql.Scannable) (*usersmodel.User, error) {
	user := &usersmodel.User{}

	scanArgs := []interface{}{
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Alias,
		&user.Email,
		&user.DateCreated,
	}

	if err := row.Scan(scanArgs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Errorf("user not found: %v", err))
		}
		return nil, fmt.Errorf(fmt.Sprintf("error scanning user: %v", err))
	}

	return user, nil
}
