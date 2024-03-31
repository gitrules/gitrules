package notice

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/ns"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
)

func SaveNoticeQueue(
	ctx context.Context,
	addr gov.Address,
	filepath ns.NS,
	queue *NoticeQueue,
) git.Change[*NoticeQueue, form.None] {

	cloned := gov.Clone(ctx, addr)
	chg := SaveNoticeQueue_StageOnly(ctx, cloned, filepath, queue)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func SaveNoticeQueue_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	filepath ns.NS,
	queue *NoticeQueue,
) git.Change[*NoticeQueue, form.None] {

	git.ToFileStage[*NoticeQueue](ctx, cloned.Tree(), filepath, queue)
	return git.NewChange[*NoticeQueue, form.None](
		"Save notice queue",
		"notice_save_queue",
		queue,
		form.None{},
		nil,
	)
}

func LoadNoticeQueue(
	ctx context.Context,
	addr gov.Address,
	filepath ns.NS,
) *NoticeQueue {

	return LoadNoticeQueue_Local(ctx, gov.Clone(ctx, addr), filepath)
}

func LoadNoticeQueue_Local(
	ctx context.Context,
	cloned gov.Cloned,
	filepath ns.NS,
) *NoticeQueue {

	queue, err := git.TryFromFile[*NoticeQueue](ctx, cloned.Tree(), filepath)
	if git.IsNotExist(err) {
		return NewNoticeQueue()
	}
	must.NoError(ctx, err)
	return queue
}
