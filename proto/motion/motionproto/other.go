package motionproto

import "github.com/gitrules/gitrules/proto/history/metric"

type Decision string

func (x Decision) String() string {
	return string(x)
}

func (x Decision) IsEmpty() bool {
	return x == ""
}

func (x Decision) IsAccept() bool {
	return x == Accept
}

func (x Decision) IsReject() bool {
	return x == Reject
}

func (x Decision) MetricDecision() metric.MotionDecision {
	return metric.MotionDecision(x)
}

var (
	Accept Decision = "accept"
	Reject Decision = "reject"
)
