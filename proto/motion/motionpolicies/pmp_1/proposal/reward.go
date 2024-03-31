package proposal

import (
	"sort"

	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/history/metric"
	"github.com/gitrules/gitrules/proto/member"
)

type Reward struct {
	To     member.User     `json:"to"`
	Amount account.Holding `json:"amount"`
}

func (x Reward) MetricReceipt() metric.Receipt {
	return metric.Receipt{
		To:     metric.AccountID(member.UserAccountID(x.To)),
		Type:   metric.ReceiptTypeReward,
		Amount: x.Amount.MetricHolding(),
	}
}

type Rewards []Reward

func (x Rewards) Len() int {
	return len(x)
}

func (x Rewards) Less(i, j int) bool {
	return x[i].To < x[j].To
}

func (x Rewards) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Rewards) Sort() {
	sort.Sort(x)
}

func (x Rewards) MetricReceipts() metric.Receipts {
	r := make(metric.Receipts, len(x))
	for i := range x {
		r[i] = x[i].MetricReceipt()
	}
	return r
}
