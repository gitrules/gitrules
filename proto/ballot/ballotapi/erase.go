package ballotapi

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/id"
)

func Erase(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotID ballotproto.BallotID,

) git.Change[form.Map, bool] {

	cloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Erase_StageOnly(ctx, cloned, ballotID)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Erase_StageOnly(
	ctx context.Context,
	cloned id.OwnerCloned,
	ballotID ballotproto.BallotID,

) git.Change[form.Map, bool] {

	t := cloned.Public.Tree()

	// verify ad is present
	ballotio.LoadAd_Local(ctx, t, ballotID)

	// erase
	ballotproto.BallotKV.Remove(ctx, ballotproto.BallotNS, cloned.PublicClone().Tree(), ballotID)

	return git.NewChange(
		fmt.Sprintf("Erased ballot %v", ballotID),
		"ballot_erase",
		form.Map{"name": ballotID},
		true,
		nil,
	)
}
