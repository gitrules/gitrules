package member

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/id"
)

func FindClonedUser_Local(
	ctx context.Context,
	cloned gov.Cloned,
	userCloned id.OwnerCloned,

) User {

	voterCred := id.GetPublicCredentials(ctx, userCloned.Public.Tree())
	users := LookupUserByID_Local(ctx, cloned, voterCred.ID)
	must.Assertf(ctx, len(users) > 0, "user not found in community")
	return users[0]
}
