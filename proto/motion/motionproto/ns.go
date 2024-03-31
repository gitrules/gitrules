package motionproto

import (
	"github.com/gitrules/gitrules/lib/ns"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/kv"
	"github.com/gitrules/gitrules/proto/motion"
)

var (
	MotionNS = proto.RootNS.Append("motion")
	MotionKV = kv.KV[MotionID, Motion]{}
)

func MotionNoticesNS(id MotionID) ns.NS {
	return MotionKV.KeyNS(MotionNS, id).Append("notices.json")
}

func MotionAccountID(motionID MotionID) account.AccountID {
	return account.AccountIDFromLine(account.Pair("motion", motionID.String()))
}

var (
	// PoliciesNS is a namespace for holding individual policy class namespaces.
	PoliciesNS = proto.PolicyNS.Append("motion")

	PolicyStateFilebase = "state.json"
)

func PolicyNS(policyName motion.PolicyName) ns.NS {
	return PoliciesNS.Append(policyName.String())
}
