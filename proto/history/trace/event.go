package trace

import (
	"github.com/gitrules/gitrules/proto/history"
)

var (
	traceHistoryNS = history.HistoryNS.Append("trace")
	traceHistory   = History{Root: traceHistoryNS}
)

type Event struct {
	Op     string `json:"op"`
	Note   string `json:"note"`
	Args   M      `json:"args,omitempty"`
	Result M      `json:"result,omitempty"`
}

type M = map[string]any
