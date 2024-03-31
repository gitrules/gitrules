package sv

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
)

func (qv SV) Cancel(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,
) git.Change[form.Map, ballotproto.Outcome] {

	// refund users
	refunded := map[member.User]account.Holding{}
	for user, spent := range tally.Charges {
		refund := account.H(account.PluralAsset, spent)
		account.Transfer_StageOnly(
			ctx,
			govOwner.PublicClone(),
			ballotproto.BallotEscrowAccountID(ad.ID),
			member.UserAccountID(user),
			refund,
			fmt.Sprintf("refund from cancelling ballot %v", ad.ID),
		)
		refunded[user] = refund
	}

	return git.NewChange(
		fmt.Sprintf("cancelled ballot %v and refunded voters", ad.ID),
		"ballot_qv_cancel",
		form.Map{"id": ad.ID},
		ballotproto.Outcome{
			Summary:      "cancelled",
			Scores:       tally.Scores,
			ScoresByUser: tally.ScoresByUser,
			Refunded:     refunded,
		},
		nil,
	)
}
