package ballotproto

import (
	"github.com/gitrules/gitrules/lib/util"
	"github.com/gitrules/gitrules/proto/member"
)

func AdsToBallotNames(ads []Ad) []string {
	names := make([]string, len(ads))
	for i := range ads {
		names[i] = ads[i].ID.GitPath()
	}
	return names
}

func FilterFrozenAds(frozen bool, ads []Ad) []Ad {
	r := []Ad{}
	for _, ad := range ads {
		if ad.Frozen == frozen {
			r = append(r, ad)
		}
	}
	return r
}

func FilterOpenClosedAds(closed bool, ads []Ad) []Ad {
	r := []Ad{}
	for _, ad := range ads {
		if ad.Closed == closed {
			r = append(r, ad)
		}
	}
	return r
}

func FilterWithParticipants(groups []member.Group, ads []Ad) []Ad {
	r := []Ad{}
	for _, ad := range ads {
		if util.IsIn(ad.Participants, groups...) {
			r = append(r, ad)
		}
	}
	return r
}
