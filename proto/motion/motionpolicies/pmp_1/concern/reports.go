package concern

import "github.com/gitrules/gitrules/proto/ballot/ballotproto"

type CloseReport struct {
}

type CancelReport struct {
	PriorityPollOutcome ballotproto.Outcome `json:"priority_poll_outcome"`
}
