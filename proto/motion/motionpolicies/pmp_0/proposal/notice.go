package proposal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/pmp_0"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/proto/notice"
)

func cancelNotice(
	ctx context.Context,
	motion motionproto.Motion,
	againstPopular bool,
	outcome ballotproto.Outcome,
) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This unmerged PR, managed as GitRules proposal `%v`, has been cancelled 🌂\n\n", motion.ID)

	if againstPopular {
		fmt.Fprintf(&w, "⚠️ Note that the PR was cancelled against the popular vote.\n\n")
	}

	fmt.Fprintf(&w, "The PR approval tally was `%0.6f`.\n\n", outcome.Scores[pmp_0.ProposalBallotChoice])

	// refunded
	fmt.Fprintf(&w, "Refunds issued:\n")
	for _, refund := range ballotproto.FlattenRefunds(outcome.Refunded) {
		fmt.Fprintf(&w, "- User @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
	}
	fmt.Fprintln(&w, "")

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by user:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- User @%v contributed `%0.6f` votes\n", user, ss[pmp_0.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(ctx, w.String())
}

func closeNotice(
	ctx context.Context,
	prop motionproto.Motion,
	againstPopular bool,
	outcome ballotproto.Outcome,
	resolved motionproto.Motions,
	bounty account.Holding,
	bountyDonated bool,
	rewards Rewards,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This PR, managed as GitRules proposal `%v`, has been closed 🎉\n\n", prop.ID)

	if againstPopular {
		fmt.Fprintf(&w, "⚠️ Note that the PR was merged against the popular vote.\n\n")
	}

	fmt.Fprintf(&w, "The PR approval tally was `%0.6f`.\n\n", outcome.Scores[pmp_0.ProposalBallotChoice])

	// bounty
	if bountyDonated {
		fmt.Fprintf(&w, "Bounty of `%0.6f` credits was donated to the community's matching fund.\n\n", bounty.Quantity)
	} else {
		fmt.Fprintf(&w, "Bounty of `%0.6f` credits was awarded to @%v.\n\n", bounty.Quantity, prop.Author)
	}

	// resolved issues
	fmt.Fprintf(&w, "Resolved issues:\n")
	for _, con := range resolved {
		fmt.Fprintf(&w, "- [Issue #%v](%v)\n", con.ID, con.TrackerURL)
	}
	fmt.Fprintln(&w, "")

	// rewarded reviewers
	fmt.Fprintf(&w, "Rewarded PR reviewers:\n")
	for _, reward := range rewards {
		fmt.Fprintf(&w, "- Reviewer @%v was awarded `%0.6f` credits\n", reward.To, reward.Amount.Quantity)
	}
	fmt.Fprintln(&w, "")

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by user:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- User @%v contributed `%0.6f` votes\n", user, ss[pmp_0.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(ctx, w.String())
}
