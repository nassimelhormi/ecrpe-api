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
			return nil, gqlerror.Errorf("%w", user.Error)
		}
	case gqlgen.RoleUser:
		if user := interceptors.ForUserContext(ctx); !user.IsAuth {
			return nil, gqlerror.Errorf("%w", user.Error)
		}
	}
	return next(ctx)
}
