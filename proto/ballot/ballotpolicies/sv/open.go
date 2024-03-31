package sv

import (
	"context"

	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
)

func (qv SV) Open(
	ctx context.Context,
	owner gov.OwnerCloned,
	ad *ballotproto.Ad,

) *ballotproto.Tally {

	return &ballotproto.Tally{
		Ad:            *ad,
		Scores:        map[string]float64{},
		ScoresByUser:  map[member.User]map[string]ballotproto.StrengthAndScore{},
		AcceptedVotes: map[member.User]ballotproto.AcceptedElections{},
		RejectedVotes: map[member.User]ballotproto.RejectedElections{},
		Charges:       map[member.User]float64{},
	}
}
