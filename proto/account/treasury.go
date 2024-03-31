package account

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/proto/gov"
)

var (
	NobodyAccountID = AccountIDFromLine(
		Term("nobody"),
	)
	TreasuryAccountID = AccountIDFromLine(
		Term("treasury"),
	)
	IssueAccountID = AccountIDFromLine(
		Cat(
			Term("treasury"),
			Term("issue"),
		),
	)
	BurnAccountID = AccountIDFromLine(
		Cat(
			Term("treasury"),
			Term("burn"),
		),
	)
)

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

	// create burn pool account
	Create_StageOnly(
		ctx,
		cloned,
		BurnAccountID,
		TreasuryAccountID,
		fmt.Sprintf("create burn account for treasury"),
	)

	// create issue pool account
	Create_StageOnly(
		ctx,
		cloned,
		IssueAccountID,
		TreasuryAccountID,
		fmt.Sprintf("create issue account for treasury"),
	)

}
