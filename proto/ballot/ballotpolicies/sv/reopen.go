package sv

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func (qv SV) Reopen(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,
) git.Change[form.Map, form.None] {

	return git.NewChange(
		fmt.Sprintf("reopened ballot %v", ad.ID),
		"ballot_qv_reopen",
		form.Map{"id": ad.ID},
		form.None{},
		nil,
	)
}
