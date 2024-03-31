package sv

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func (qv SV) Close(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,
) git.Change[form.Map, ballotproto.Outcome] {

	return git.NewChange(
		fmt.Sprintf("closed ballot %v", ad.ID),
		"ballot_qv_close",
		form.Map{"id": ad.ID},
		ballotproto.Outcome{
			Summary:      "closed",
			Scores:       tally.Scores,
			ScoresByUser: tally.ScoresByUser,
			Refunded:     nil,
		},
		nil,
	)
}
