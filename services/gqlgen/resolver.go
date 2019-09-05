package gqlgen

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nassimelhormi/ecrpe-api/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	DB        *sqlx.DB
	SecreyKey string
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
	_, err := r.DB.Exec(
		"INSERT INTO users (username, email, is_teacher, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, user.IsTeacher, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdatedUser) (*models.User, error) {
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

	if _, err := r.DB.Queryx(query.String()); err != nil {
		log.Fatal(err)
	}

	user := models.User{}
	err := r.DB.Get(&user, "SELECT id, username, email, is_teacher FROM users WHERE username = ?", input.Username)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}
func (r *mutationResolver) PurchaseRefresherCourse(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)
	if err := r.DB.Select(users, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users"); err != nil {
		return users, nil
	}
	return users, nil
}
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users WHERE id = ?", id); err != nil {
		return &models.User{}, nil
	}
	return &user, nil
}
func (r *queryResolver) MyCourses(ctx context.Context, userID int) ([]*models.RefresherCourse, error) {
	refCourses := make([]*models.RefresherCourse, 0)
	if err := r.DB.Select(refCourses, `
		SELECT * FROM refresher_courses
		JOIN users_refresher_courses ON refresher_courses.id = users_refresher_courses.refresher_course_id
		JOIN users ON users_refresher_courses.user_id = ?
	`, userID); err != nil {
		return refCourses, err
	}
	return refCourses, nil
}
func (r *queryResolver) RefresherCourses(ctx context.Context, subjectID *int) ([]*models.RefresherCourse, error) {
	refCourses := make([]*models.RefresherCourse, 0)
	if subjectID == nil {
		return refCourses, nil
	}
	if err := r.DB.Select(refCourses, `
		SELECT * FROM refresher_courses
		JOIN subjects_refresher_courses ON refresher_courses.id = subjects_refresher_courses.refresher_course_id
		JOIN subjects ON subjects_refresher_courses.subject_id = ?
	`, subjectID); err != nil {
		return refCourses, err
	}
	return refCourses, nil
}
func (r *queryResolver) Sessions(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)
	if err := r.DB.Select(sessions, `
		SELECT id, title, description, recorded_on, created_at, updated_at FROM sessions
		WHERE refresher_course_id = ?
	`, refresherCourseID); err != nil {
		return sessions, err
	}
	return sessions, nil
}
func (r *queryResolver) MyProfil(ctx context.Context, userID int) (*models.User, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users WHERE id = ?", userID); err != nil {
		return &models.User{}, err
	}
	return &user, nil
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
