package metric

import "github.com/gitrules/gitrules/proto/history"

var (
	metricHistoryNS = history.HistoryNS.Append("metric")
	metricHistory   = History{Root: metricHistoryNS}
)

type Event struct {
	Join    *JoinEvent    `json:"join,omitempty"`
	Motion  *MotionEvent  `json:"motion,omitempty"`
	Account *AccountEvent `json:"account,omitempty"`
	Vote    *VoteEvent    `json:"vote,omitempty"`
}
