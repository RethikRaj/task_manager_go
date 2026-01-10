package ctx

import "context"

func GetUserFromContext(ctx context.Context) (ContextUser, bool) {
	user, ok := ctx.Value(UserKey).(ContextUser)

	return user, ok
}
