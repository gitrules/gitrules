package ballotapi

import (
	"context"

	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func Capitalization_Local(
	ctx context.Context,
	cloned gov.Cloned,
	ballotName ballotproto.BallotID,
) float64 {

	return Show_Local(ctx, cloned, ballotName).Tally.Capitalization()
}
