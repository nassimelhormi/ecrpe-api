package gqlgen

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nassimelhormi/ecrpe-api/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root@/ecrpe")
	//	defer db.Close()
}

type Resolver struct {
	DB *sql.DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) RefresherCourse() RefresherCourseResolver {
	return &refresherCourseResolver{r}
}
func (r *Resolver) Session() SessionResolver {
	return &sessionResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*models.User, error) {
	user := &models.User{
		Username:  input.Username,
		Email:     input.Email,
		IsTeacher: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := db.Exec(
		"INSERT INTO users (username, email, is_teacher, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username,
		user.Email,
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
	query := strings.Builder{}

	query.WriteString("UPDATE users SET ")

	if input.Email != nil && *input.Email != "" {
		query.WriteString(fmt.Sprintf("email = '%s'", *input.Email))
		query.WriteString(", ")
	}
	if input.Username != nil && *input.Username != "" {
		query.WriteString(fmt.Sprintf("username = '%s'", *input.Username))
	}
	query.WriteString(fmt.Sprintf(" WHERE username = '%s'", *input.Username))

	if _, err := db.Exec(query.String()); err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow(
		"SELECT id, username, email, is_teacher FROM users WHERE username = ?",
		input.Username,
	)
	if errScan := row.Scan(&user.ID, &user.Username, &user.Email, &user.IsTeacher); errScan != nil {
		log.Fatal(errScan)
	}

	return user, nil
}
func (r *mutationResolver) PurchaseRefresherCourse(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	panic("not implemented")
}
func (r *queryResolver) MyCourses(ctx context.Context, userID int) ([]*models.RefresherCourse, error) {
	panic("not implemented")
}
func (r *queryResolver) RefresherCourses(ctx context.Context, subjectID *int) ([]*models.RefresherCourse, error) {
	panic("not implemented")
}
func (r *queryResolver) Sessions(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	panic("not implemented")
}
func (r *queryResolver) MyProfil(ctx context.Context, userID int) (*models.User, error) {
	panic("not implemented")
}

type refresherCourseResolver struct{ *Resolver }

func (r *refresherCourseResolver) CreatedAt(ctx context.Context, obj *models.RefresherCourse) (string, error) {
	panic("not implemented")
}
func (r *refresherCourseResolver) UpdatedAt(ctx context.Context, obj *models.RefresherCourse) (string, error) {
	panic("not implemented")
}

type sessionResolver struct{ *Resolver }

func (r *sessionResolver) RecordedOn(ctx context.Context, obj *models.Session) (*string, error) {
	panic("not implemented")
}
func (r *sessionResolver) CreatedAt(ctx context.Context, obj *models.Session) (string, error) {
	panic("not implemented")
}
func (r *sessionResolver) UpdatedAt(ctx context.Context, obj *models.Session) (string, error) {
	panic("not implemented")
}

type userResolver struct{ *Resolver }

func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	panic("not implemented")
}
func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	panic("not implemented")
}
