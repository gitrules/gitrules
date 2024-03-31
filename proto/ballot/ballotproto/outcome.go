package ballotproto

import (
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/history/metric"
	"github.com/gitrules/gitrules/proto/member"
)

type Summary string

type Outcome struct {
	Summary      string                                      `json:"summary"`
	Scores       map[string]float64                          `json:"scores"`
	ScoresByUser map[member.User]map[string]StrengthAndScore `json:"scores_by_user"`
	Refunded     map[member.User]account.Holding             `json:"refunded"`
}

func (o Outcome) RefundedHistoryReceipts() metric.Receipts {
	r := metric.Receipts{}
	for user, h := range o.Refunded {
		r = append(r,
			metric.Receipt{
				To:     user.MetricAccountID(),
				Type:   metric.ReceiptTypeRefund,
				Amount: h.MetricHolding(),
			},
		)
	}
	return r
}
