package sv

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/proto/member"
)

func (qv SV) VerifyElections(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.Address,
	voterCloned id.OwnerCloned,
	govCloned gov.Cloned,
	ad *ballotproto.Ad,
	prior *ballotproto.Tally,
	elections ballotproto.Elections,
) {

	voterCred := id.GetPublicCredentials(ctx, voterCloned.Public.Tree())
	user := member.LookupUserByID_Local(ctx, govCloned, voterCred.ID)
	if len(user) == 0 {
		must.Errorf(ctx, "cannot find user with id %v in the community", voterCred.ID)
	}

	// tally writes to the gov repo, but the repo is throw-away and won't be committed
	qv.tally(ctx, govCloned, ad, prior, map[member.User]ballotproto.Elections{user[0]: elections}, true)
}
