package gqlgen

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nassimelhormi/ecrpe-api/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root@/ecrpe")
	//defer db.Close()
}

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*models.User, error) {
	user := &models.User{
		Username:    input.Username,
		PhoneNumber: *input.PhoneNumber,
		Email:       input.Email,
		CurrentRank: *input.CurrentRank,
		IsTeacher:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := db.Exec(
		"INSERT INTO users (username, phone_number, email, current_rank, is_teacher, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Username,
		user.PhoneNumber,
		user.Email,
		user.CurrentRank,
		user.IsTeacher,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdatedUser) (*models.User, error) {
	user := &models.User{}
	query := "UPDATE users SET"

	if input.PhoneNumber != nil && *input.PhoneNumber != "" {
		query += fmt.Sprintf(" phone_number = '%s', ", *input.PhoneNumber)
	}
	if input.Email != nil && *input.Email != "" {
		query += fmt.Sprintf(" email = '%s', ", *input.Email)
	}
	if input.CurrentRank != nil && *input.CurrentRank != 0 {
		query += fmt.Sprintf(" current_rank = %d, ", *input.CurrentRank)
	}
	query += fmt.Sprintf(" username = '%[1]s' WHERE username = '%[1]s'", input.Username)

	_, err := db.Exec(
		query,
	)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow(
		"SELECT id, username, phone_number, email, is_teacher FROM users WHERE username = ?",
		input.Username,
	)
	errScan := row.Scan(&user.ID, &user.Username, &user.PhoneNumber, &user.Email, &user.IsTeacher)
	if errScan != nil {
		log.Fatal(errScan)
	}

	return user, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	panic("not implemented")
}

type userResolver struct{ *Resolver }

func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	panic("not implemented")
}
func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	panic("not implemented")
}
