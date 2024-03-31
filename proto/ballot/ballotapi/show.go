package ballotapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func Show(
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,

) ballotproto.AdTallyMargin {

	return Show_Local(ctx, gov.Clone(ctx, addr), id)
}

func Show_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) ballotproto.AdTallyMargin {

	ad := ballotio.LoadAd_Local(ctx, cloned.Tree(), id)

	tally, _ := must.Try1[ballotproto.Tally](
		func() ballotproto.Tally {
			return loadTally_Local(ctx, cloned.Tree(), id)
		},
	)

	margin, _ := must.Try1[*ballotproto.Margin](
		func() *ballotproto.Margin {
			return GetMargin_Local(ctx, cloned, id)
		},
	)

	return ballotproto.AdTallyMargin{
		Ad:     ad,
		Tally:  tally,
		Margin: margin,
	}
}
