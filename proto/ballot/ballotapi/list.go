package ballotapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
)

func List(
	ctx context.Context,
	addr gov.Address,

) ballotproto.Advertisements {

	return List_Local(ctx, gov.Clone(ctx, addr))
}

func List_Local(
	ctx context.Context,
	cloned gov.Cloned,

) ballotproto.Advertisements {

	ids := ballotproto.BallotKV.ListKeys(ctx, ballotproto.BallotNS, cloned.Tree())

	ads := ballotproto.Advertisements{}
	for _, id := range ids {
		ad := git.FromFile[ballotproto.Ad](ctx, cloned.Tree(), id.AdNS())
		if ballotio.TryLookupPolicy(ctx, ad.Policy) != nil {
			ads = append(ads, ad)
		}
	}

	ads.Sort()
	return ads
}

func ListFilter(
	ctx context.Context,
	addr gov.Address,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,

) ballotproto.Advertisements {

	return ListFilter_Local(ctx, gov.Clone(ctx, addr), onlyOpen, onlyClosed, onlyFrozen, withParticipant)
}

func ListFilter_Local(
	ctx context.Context,
	cloned gov.Cloned,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,

) ballotproto.Advertisements {

	ads := List_Local(ctx, cloned)
	if onlyOpen {
		ads = ballotproto.FilterOpenClosedAds(false, ads)
	}
	if onlyClosed {
		ads = ballotproto.FilterOpenClosedAds(true, ads)
	}
	if onlyFrozen {
		ads = ballotproto.FilterFrozenAds(true, ads)
	}
	if withParticipant != "" {
		userGroups := member.ListUserGroups_Local(ctx, cloned, withParticipant)
		ads = ballotproto.FilterWithParticipants(userGroups, ads)
	}
	ads.Sort()
	return ads
}
