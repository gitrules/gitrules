package bureau

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/proto/mail"
	"github.com/gitrules/gitrules/proto/member"
)

func Transfer(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.Address,
	fromUserOpt member.User, // optional, if empty string, a lookup forthe user is performed
	toUser member.User,
	amount float64,
) git.Change[form.Map, mail.RequestEnvelope[Request]] {

	govCloned := gov.Clone(ctx, govAddr)
	userOwner := id.CloneOwner(ctx, userAddr)
	chg := Transfer_StageOnly(ctx, userAddr, userOwner, govCloned, fromUserOpt, toUser, amount)
	proto.Commit(ctx, userOwner.Public.Tree(), chg)
	userOwner.Public.Push(ctx)
	return chg
}

func Transfer_StageOnly(
	ctx context.Context,
	userAddr id.OwnerAddress,
	userOwner id.OwnerCloned,
	govCloned gov.Cloned,
	fromUserOpt member.User,
	toUser member.User,
	amount float64,
) git.Change[form.Map, mail.RequestEnvelope[Request]] {

	userCred := id.GetPublicCredentials(ctx, userOwner.Public.Tree())

	// find the user name of userAddr in the community repo
	if fromUserOpt == "" {
		us := member.LookupUserByID_Local(ctx, govCloned, userCred.ID)
		switch len(us) {
		case 0:
			must.Errorf(ctx, "%s not found in community %v", userAddr.Public, govCloned.Address())
		case 1:
			fromUserOpt = us[0]
		default:
			must.Errorf(ctx, "community %v has more than one user at address %v", govCloned.Address(), userAddr.Public)
		}
	}

	request := Request{
		Transfer: &TransferRequest{
			FromUser: fromUserOpt,
			ToUser:   toUser,
			Amount:   amount,
		},
	}

	sendOnly := mail.Request_StageOnly(ctx, userOwner, govCloned.Tree(), BureauTopic, request)
	return git.NewChange(
		"Transfer account tokens.",
		"bureau_transfer",
		form.Map{
			"from_user": fromUserOpt,
			"to_user":   toUser,
			"amount":    amount,
		},
		sendOnly.Result,
		form.Forms{sendOnly},
	)
}
