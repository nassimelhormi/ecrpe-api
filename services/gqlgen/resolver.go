package gqlgen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vektah/gqlparser/gqlerror"

	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/utils"

	"github.com/gbrlsnchs/jwt"

	"github.com/jmoiron/sqlx"
	"github.com/nassimelhormi/ecrpe-api/models"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
	"golang.org/x/crypto/bcrypt"
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
		return nil, gqlerror.Errorf("cannot create user")
	}
	return user, nil
}
func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdatedUser) (*models.User, error) {
	if user := interceptors.ForContext(ctx); !user.IsAuth {
		return &models.User{}, gqlerror.Errorf("%w", user.Error)
	}
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
		return &models.User{}, gqlerror.Errorf("cannot update user")
	}
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher FROM users WHERE username = ?", input.Username); err != nil {
		return &models.User{}, gqlerror.Errorf("cannot retrieve user")
	}
	return &user, nil
}
func (r *mutationResolver) PurchaseRefresherCourse(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateSession(ctx context.Context, input NewSession) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) RefreshToken(ctx context.Context) (*models.Token, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)
	if err := r.DB.Select(users, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users"); err != nil {
		return users, gqlerror.Errorf("cannot retrieve users")
	}
	return users, nil
}
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users WHERE id = ?", id); err != nil {
		return &models.User{}, gqlerror.Errorf("cannot retrieve user")
	}
	return &user, nil
}
func (r *queryResolver) AuthUser(ctx context.Context, input UserLogin) (*models.Token, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, encrypted_pwd FROM users WHERE email = ?", input.Email); err != nil {
		return &models.Token{}, gqlerror.Errorf("cannot find user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPWD), []byte(input.Password)); err != nil {
		return &models.Token{}, gqlerror.Errorf("pwd incorrect")
	}
	tokens := models.Token{}
	pl := jwt.Payload{
		Subject:        user.Username,
		Audience:       jwt.Audience{"https://ecrpe.fr"},
		ExpirationTime: jwt.NumericDate(time.Now().Add(30 * time.Minute)),
		IssuedAt:       jwt.NumericDate(time.Now()),
	}
	jwt, err := jwt.Sign(pl, jwt.NewHS256([]byte(r.SecreyKey)))
	if err != nil {
		return &models.Token{}, gqlerror.Errorf("error jwt")
	}
	tokens.JWT = string(jwt)
	tokens.RefreshToken = utils.HexKeyGenerator(16)
	if _, err := r.DB.Queryx(`
		UPDATE user_auths SET ip_adress=?, refresh_token=?, delivered_at=?, on_login=?, user_id=?
	`, "ipadress", tokens.RefreshToken, time.Now(), 1, user.ID); err != nil {
		return &models.Token{}, gqlerror.Errorf("token error")
	}
	return &tokens, nil
}
func (r *queryResolver) MyCourses(ctx context.Context, userID int) ([]*models.RefresherCourse, error) {
	refCourses := make([]*models.RefresherCourse, 0)
	if user := interceptors.ForContext(ctx); !user.IsAuth {
		return refCourses, gqlerror.Errorf("%w", user.Error)
	}
	if err := r.DB.Select(refCourses, `
		SELECT * FROM refresher_courses
		JOIN users_refresher_courses ON refresher_courses.id = users_refresher_courses.refresher_course_id
		JOIN users ON users_refresher_courses.user_id = ?
	`, userID); err != nil {
		return refCourses, gqlerror.Errorf("cannot retrieve your refresher courses purchased")
	}
	return refCourses, nil
}
func (r *queryResolver) RefresherCourses(ctx context.Context, subjectID *int) ([]*models.RefresherCourse, error) {
	refCourses := make([]*models.RefresherCourse, 0)
	if subjectID == nil {
		if err := r.DB.Select(refCourses, "SELECT * FROM refresher_courses"); err != nil {
			return refCourses, gqlerror.Errorf("cannot retrieve refresher courses, try again")
		}
		return refCourses, nil
	}
	if err := r.DB.Select(refCourses, `
		SELECT * FROM refresher_courses
		JOIN subjects_refresher_courses ON refresher_courses.id = subjects_refresher_courses.refresher_course_id
		JOIN subjects ON subjects_refresher_courses.subject_id = ?
	`, subjectID); err != nil {
		return refCourses, gqlerror.Errorf("cannot retrieve refresher courses from your subject choice")
	}
	return refCourses, nil
}
func (r *queryResolver) Sessions(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)
	if user := interceptors.ForContext(ctx); !user.IsAuth {
		return sessions, gqlerror.Errorf("%w", user.Error)
	}
	if err := r.DB.Select(sessions, `
		SELECT id, title, description, recorded_on, created_at, updated_at FROM sessions
		WHERE refresher_course_id = ?
	`, refresherCourseID); err != nil {
		return sessions, gqlerror.Errorf("cannot retrieve sessions from this refresher course")
	}
	return sessions, nil
}
func (r *queryResolver) MyProfil(ctx context.Context, userID int) (*models.User, error) {
	if user := interceptors.ForContext(ctx); !user.IsAuth {
		return &models.User{}, gqlerror.Errorf("%w", user.Error)
	}
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users WHERE id = ?", userID); err != nil {
		return &models.User{}, gqlerror.Errorf("cannot access your profil, try again")
	}
	return &user, nil
}
func (r *queryResolver) OneUserAuth(ctx context.Context) (string, error) {
	if user := interceptors.ForContext(ctx); !user.IsAuth {
		return "", gqlerror.Errorf("%w", user.Error)
	}
	// redis check ip from subject jwt (username)
	return "", nil
}
