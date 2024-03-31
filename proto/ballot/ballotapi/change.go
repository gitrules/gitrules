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
)

func Change(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,
	title string,
	description string,

) git.Change[form.Map, ballotproto.Ad] {

	cloned := gov.CloneOwner(ctx, addr)

	chg := Change_StageOnly(ctx, cloned, id, title, description)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Change_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	title string,
	description string,

) git.Change[form.Map, ballotproto.Ad] {

	ad := ballotio.LoadAd_Local(ctx, cloned.Public.Tree(), id)
	ad.Title = title
	ad.Description = description
	git.ToFileStage(ctx, cloned.Public.Tree(), id.AdNS(), ad)

	return git.NewChange(
		fmt.Sprintf("Change ballot %v info", id),
		"ballot_change",
		form.Map{"name": id, "title": title, "description": description},
		ad,
		nil,
	)
}
