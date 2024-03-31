package sync

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/bureau"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
)

func Sync(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	maxPar int,
) git.Change[form.Map, form.Map] {

	// collect votes and tally all open ballots
	tallyChg := ballotapi.TallyAll(ctx, govAddr, maxPar)

	// process bureau requests by users
	bureauChg := bureau.Process(ctx, govAddr, member.Everybody)

	return git.NewChange(
		"Governance-community sync",
		"sync_sync",
		form.Map{},
		form.Map{
			"tally_result":  tallyChg.Result,
			"bureau_result": bureauChg.Result,
		},
		form.Forms{tallyChg, bureauChg},
	)
}
