package waimea

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

var (
	WaimeaAccountID = account.AccountIDFromLine(
		account.Term("waimea"),
	)
	PennyAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("waimea"),
			account.Term("penny"),
		),
	)
)

func ConcernAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("waimea-concern"),
		),
	)
}

func ProposalAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("waimea-proposal"),
		),
	)
}

func ProposalBountyAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("waimea-proposal-bounty"),
		),
	)
}

func ProposalRewardAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("waimea-proposal-reward"),
		),
	)
}

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

	must.Try(
		func() {
			account.Create_StageOnly(
				ctx,
				cloned,
				PennyAccountID,
				WaimeaAccountID,
				fmt.Sprintf("penny account for the Waimea Protocol"),
			)
		},
	)

}
