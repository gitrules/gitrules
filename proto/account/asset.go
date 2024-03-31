package account

import "github.com/gitrules/gitrules/proto/history/metric"

var (
	PluralAsset = Asset("plural")
)

type Asset string

func (a Asset) String() string {
	return string(a)
}

func (a Asset) MetricAsset() metric.Asset {
	return metric.Asset(a)
}
