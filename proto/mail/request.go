package mail

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/id"
)

func Request_StageOnly[Req form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	req Req,
) git.Change[form.Map, RequestEnvelope[Req]] {

	mkReqEnv := func(_ context.Context, seqNo SeqNo) RequestEnvelope[Req] {
		return RequestEnvelope[Req]{
			SeqNo:   seqNo,
			Request: req,
		}
	}

	chg := SendSignedMakeMsg_StageOnly(ctx, senderCloned, receiver, topic, mkReqEnv)
	msg := chg.Result.Msg

	return git.NewChange(
		fmt.Sprintf("Requested #%d", chg.Result.SeqNo),
		"request",
		form.Map{"topic": topic, "msg": msg},
		msg,
		form.Forms{chg},
	)
}
