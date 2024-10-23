package lib

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
	"github.com/gitrules/gitrules/proto/notice"
	"github.com/google/go-github/v66/github"
)

func DisplayNotices_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
) {

	t := cloned.Tree()
	motions := motionapi.ListMotions_Local(ctx, t)
	for _, motion := range motions {
		issueNum, err := MotionIDToIssueNumber(motion.ID)
		if err != nil {
			base.Errorf("encountered motion %v whose id cannot be converted into a github issue number", motion.ID)
			continue
		}
		queue := motionapi.LoadMotionNotices_Local(ctx, cloned, motion.ID)
		flushNotices(ctx, repo, ghc, cloned, queue, issueNum)
		motionapi.SaveMotionNotices_StageOnly(ctx, cloned, motion.ID, queue)
	}
}

func flushNotices(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
	queue *notice.NoticeQueue,
	issueNum int,
) {

	var w bytes.Buffer

	notShown := 0
	for _, nstate := range queue.NoticeStates {

		// check if notice already displayed, based on governance records
		if nstate.IsShown() {
			continue
		}

		// TODO: check if notice already displayed, according to github

		fmt.Fprintf(&w, "#### Notice `%v`\n%s\n\n", nstate.ID, nstate.Notice.Body)
		nstate.MarkShown()
		notShown++
	}

	if notShown > 0 {
		replyToIssue(ctx, repo, ghc, issueNum, "GitRules notices", w.String())
	}
}
