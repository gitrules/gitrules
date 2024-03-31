package ballotapi

import (
	"context"

	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func GetMargin(
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,

) *ballotproto.Margin {

	cloned := gov.Clone(ctx, addr)
	return GetMargin_Local(ctx, cloned, id)
}

func GetMargin_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) *ballotproto.Margin {

	t := cloned.Tree()
	ad, policy := ballotio.LoadAdPolicy_Local(ctx, t, id)
	tally := loadTally_Local(ctx, t, id)
	return policy.Margin(ctx, cloned, &ad, &tally)
}
