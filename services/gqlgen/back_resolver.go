package gqlgen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/nassimelhormi/ecrpe-api/models"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/redis"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/utils"
	"github.com/vektah/gqlparser/gqlerror"
	"golang.org/x/crypto/bcrypt"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	DB              *sqlx.DB
	SecreyKey       string
	IPAddressCache  *redis.Cache
	VideoEncodingCh chan models.Video
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (bool, error) {
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
	}
	if _, err := r.DB.Exec(
		"INSERT INTO users (username, email, is_teacher, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, user.IsTeacher, time.Now(), time.Now(),
	); err != nil {
		return false, gqlerror.Errorf("cannot create user")
	}
	return true, nil
}
func (r *mutationResolver) RefreshToken(ctx context.Context) (*models.Token, error) {
	refreshToken := interceptors.ForRefreshToken(ctx)
	if refreshToken == "" {
		return &models.Token{}, gqlerror.Errorf("No refresh token provided")
	}
	// Check refreshToken validity
	userAuth := models.UserAuth{}
	if err := r.DB.Get(userAuth, `
		SELECT u.username, ua.id, ua.user_id
		FROM user u, user_auths ua 
		WHERE ua.refresh_token = ?
	`, refreshToken); err != nil {
		return &models.Token{}, gqlerror.Errorf("")
	}
	// Get IP Address from user
	userIP := interceptors.ForIPAddress(ctx)
	if _, err := r.DB.Queryx(`
		UPDATE user_auths SET ip_address=?, is_revoked=?, revoked_at=?, on_refresh=? WHERE id = ?
	`, userIP, 1, time.Now(), 1, userAuth.ID); err != nil {
		return &models.Token{}, gqlerror.Errorf("token error")
	}
	// Create new jwt then new refresh token
	pl := utils.CustomPayload{
		Payload: jwt.Payload{
			Audience:       jwt.Audience{"https://ecrpe.fr"},
			ExpirationTime: jwt.NumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:       jwt.NumericDate(time.Now()),
		},
		Username: userAuth.Username,
		UserID:   userAuth.ID,
	}
	jwt, err := jwt.Sign(pl, jwt.NewHS256([]byte(r.SecreyKey)))
	if err != nil {
		return &models.Token{}, gqlerror.Errorf("error occured during new jwt creation")
	}
	tokens := models.Token{
		JWT:          string(jwt),
		RefreshToken: utils.HexKeyGenerator(16),
	}
	// push refresh token
	if _, err = r.DB.Exec(
		"INSERT INTO user_auths (ip_address, refresh_token, delivered_at, on_refresh, user_id) VALUES (?, ?, ?, ?, ?, ?)",
		userIP, tokens.RefreshToken, time.Now(), 1, userAuth.UserID,
	); err != nil {
		return nil, gqlerror.Errorf("error occured during new refresh token creation")
	}
	return &tokens, nil
}
func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdatedUser) (bool, error) {
	// *add empty string front end so that input.email != nil and input.Username != nil get removed*
	if (input.Email != nil && *input.Email != "") && (input.Username != nil && *input.Username != "") {
		return false, gqlerror.Errorf("Nothing to update")
	}
	query := strings.Builder{}
	query.WriteString("UPDATE users SET ")

	if input.Email != nil && *input.Email != "" {
		query.WriteString(fmt.Sprintf("email = '%s'", *input.Email))
		if input.Username != nil && *input.Username != "" {
			query.WriteString(", ")
		}
	}
	if input.Username != nil && *input.Username != "" {
		query.WriteString(fmt.Sprintf("username = '%s'", *input.Username))
	}
	query.WriteString(fmt.Sprintf(" WHERE username = '%s'", *input.Username))
	if _, err := r.DB.Queryx(query.String()); err != nil {
		return false, gqlerror.Errorf("Error occured during your update profil")
	}
	return true, nil
}
func (r *mutationResolver) PurchaseRefresherCourse(ctx context.Context, refresherCourseID int) ([]*models.Session, error) {

	sessions := make([]*models.Session, 0)
	// 	** if not logged with no subject_id **
	// 	SELECT * FROM refresher_courses
	//
	//	** if not logged with subject_id **
		// 	SELECT *
		// 	FROM refresher_courses AS rc
		// 	JOIN subjects_refresher_courses AS src
		//		ON rc.id = src.refresher_course_id
		//	JOIN subjects AS s
		//		ON s.id = src.subject_id
		// 	WHERE s.id = ?
	//
	// 	** if logged with subject_id **
	// 	SELECT id, year, is_finished, created_at, updated_at,
	//  IF(id
	//			NOT IN(
	//				SELECT refresher_course_id FROM users_refresher_courses
	//				WHERE user_id = ?
	//			), price
	//  )
	// 	FROM refresher_courses AS rc
	// 	JOIN subjects_refresher_courses AS src
	//		ON rc.id = src.refresher_course_id
	//	JOIN subjects AS s
	//		ON s.id = src.subject_id
	// 	WHERE s.id = ?
	//
	// 	** logged without subject_id **
	// 	SELECT id, year, is_finished, created_at, updated_at,
	//  IF(id
	//			NOT IN(
	//				SELECT refresher_course_id FROM users_refresher_courses
	//				WHERE user_id = ?
	//			), price
	//  )
	// 	FROM refresher_courses

	// paypal system
	return sessions, nil
}
func (r *mutationResolver) CreateRefresherCourse(ctx context.Context, input NewSessionCourse) (bool, error) {
	// **add refresherID**
	refresherCourse := struct {
		Year string `db:"year"`
		Name string `db:"name"`
	}{}
	if err := r.DB.Get(&refresherCourse, `
		SELECT rc.id, rc.year, s.name 
		FROM refresher_courses as rc
		JOIN subjects_refresher_courses as src on (rc.id = ?)
		JOIN subjects as s on (src.subject_id = s.id)
		WHERE s.active = 1
	`); err != nil {
		return false, gqlerror.Errorf("no found")
	}
	// Create new session
	sessionObj, err := r.DB.Exec(`
		INSERT INTO sessions (title, type, description, part, recorded_on, created_at, refresher_course_id)
		VALUES (?,?,?,?,?,?,?)
	`, input.Title, input.Type, input.Description, input.Part, input.RecordedOn, time.Now(), input.RefresherCourseID)
	if err != nil {
		return false, gqlerror.Errorf("Error occured during process")
	}
	sessionID, err := sessionObj.LastInsertId()
	if err != nil {
		return false, gqlerror.Errorf("Error occured during process")
	}
	// Generate session PATH
	pathToUpload, err := utils.MakeSessionPath(refresherCourse, sessionID)
	if err != nil {
		return false, gqlerror.Errorf("Path generation not working")
	}
	fmt.Println(pathToUpload)
	// Create video then handle process treatment

	// Create class papers then paths
	if len(input.DocFiles) > 0 {

	}
	return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Login(ctx context.Context, input UserLogin) (*models.Token, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, encrypted_pwd FROM users WHERE email = ?", input.Email); err != nil {
		return &models.Token{}, gqlerror.Errorf("This email doesn't exist")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPWD), []byte(input.Password)); err != nil {
		return &models.Token{}, gqlerror.Errorf("Password incorrect")
	}
	pl := utils.CustomPayload{
		Payload: jwt.Payload{
			Audience:       jwt.Audience{"https://ecrpe.fr"},
			ExpirationTime: jwt.NumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:       jwt.NumericDate(time.Now()),
		},
		Username: user.Username,
		UserID:   user.ID,
	}

	jwt, err := jwt.Sign(pl, jwt.NewHS256([]byte(r.SecreyKey)))
	if err != nil {
		return &models.Token{}, gqlerror.Errorf("Cannot proceed further")
	}
	tokens := models.Token{
		JWT:          string(jwt),
		RefreshToken: utils.HexKeyGenerator(16),
	}

	userIP := interceptors.ForIPAddress(ctx)
	if _, err = r.DB.Exec(
		"INSERT INTO user_auths (ip_address, refresh_token, delivered_at, on_login, user_id) VALUES (?, ?, ?, ?, ?)",
		userIP, tokens.RefreshToken, time.Now(), 1, user.ID,
	); err != nil {
		return &models.Token{}, gqlerror.Errorf("Cannot proceed further")
	}
	return &tokens, nil
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
		SELECT * FROM refresher_courses as rc
		JOIN subjects_refresher_courses as src ON rc.id = src.refresher_course_id
		JOIN subjects as s ON src.subject_id = s.id WHERE s.id = ?
	`, subjectID); err != nil {
		return refCourses, gqlerror.Errorf("cannot retrieve refresher courses from your subject choice")
	}
	return refCourses, nil
}
func (r *queryResolver) RefresherCourse(ctx context.Context, refresherCourseID int) (*models.RefresherCourse, error) {
	refresherCourse := models.RefresherCourse{}
	if err := r.DB.Get(&refresherCourse, `
		SELECT * FROM refresher_courses
	`, refresherCourseID); err != nil {
		return &models.RefresherCourse{}, gqlerror.Errorf("refresherCourse")
	}
	return &refresherCourse, nil
}
func (r *queryResolver) VideoUserCheck(ctx context.Context) (bool, error) {
	user := interceptors.ForUserContext(ctx)
	if !user.IsAuth {
		return false, gqlerror.Errorf("%w", user.Error)
	}
	userIPAddress := interceptors.ForIPAddress(ctx)
	lastIPCached, ok := r.IPAddressCache.GetIP(string(user.UserID))
	if !ok {
		return false, gqlerror.Errorf("Cannot retrieve IPAddress")
	}
	err := utils.IPsChecker(userIPAddress, lastIPCached)
	if err != nil {
		r.IPAddressCache.DeleteIP(string(user.UserID))
		return false, gqlerror.Errorf("%w", err)
	}
	r.IPAddressCache.AddIP(string(user.UserID), string(user.Username))
	return true, nil
}
func (r *queryResolver) MyProfil(ctx context.Context, userID int) (*models.User, error) {
	user := models.User{}
	if err := r.DB.Get(&user, "SELECT id, username, email, is_teacher, created_at, updated_at FROM users WHERE id = ?", userID); err != nil {
		return &models.User{}, gqlerror.Errorf("Cannot access your profil")
	}
	return &user, nil
}
func (r *queryResolver) MyRefrescherCourses(ctx context.Context, userID int) ([]*models.RefresherCourse, error) {
	refCourses := make([]*models.RefresherCourse, 0)
	// add checking paypal payment
	if err := r.DB.Select(refCourses, `
		SELECT * FROM refresher_courses
		JOIN users_refresher_courses ON refresher_courses.id = users_refresher_courses.refresher_course_id
		JOIN users ON users_refresher_courses.user_id = ?
	`, userID); err != nil {
		return refCourses, gqlerror.Errorf("cannot retrieve your refresher courses purchased")
	}
	return refCourses, nil
}
func (r *queryResolver) SessionCourse(ctx context.Context, sessionID int) (*models.Session, error) {
	panic("not implemented")
}
