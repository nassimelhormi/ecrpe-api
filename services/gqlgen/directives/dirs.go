package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
	"github.com/vektah/gqlparser/gqlerror"
)

// HasRole func
func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role gqlgen.Role) (interface{}, error) {
	switch role {
	case gqlgen.RoleTeacher:
		user := interceptors.ForUserContext(ctx)
		if !user.IsAuth {
			return "", gqlerror.Errorf("%w", user.Error)
		}

		return nil, nil
	case gqlgen.RoleUser:
		user := interceptors.ForUserContext(ctx)
		if !user.IsAuth {
			return "", gqlerror.Errorf("%w", user.Error)
		}
	}
	return next(ctx)
}

// RefresherCourseOwner func
func RefresherCourseOwner(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return next(ctx)
}
