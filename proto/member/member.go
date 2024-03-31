// Package member implements governance member management services
package member

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/history/metric"
	"github.com/gitrules/gitrules/proto/history/trace"
)

const (
	Everybody = Group("everybody")
)

type User string

func (u User) IsNone() bool {
	return u == ""
}

func (u User) MetricUser() metric.User {
	return metric.User(u)
}

func (u User) MetricAccountID() metric.AccountID {
	return metric.AccountID(UserAccountID(u))
}

type Group string

func AddMember(ctx context.Context, addr gov.Address, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddMember_StageOnly(ctx, cloned, user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddMember_StageOnly(ctx context.Context, cloned gov.Cloned, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Set(ctx, userGroupsNS, cloned.Tree(), user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, cloned.Tree(), group, user, true)

	// log
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "add_user_to_group",
		Args:   trace.M{"user": user, "group": group},
		Result: nil,
	})

	return git.NewChangeNoResult(fmt.Sprintf("Added user %v to group %v", user, group), "member_add_member")
}

func IsMember(ctx context.Context, addr gov.Address, user User, group Group) bool {
	return IsMember_Local(ctx, gov.Clone(ctx, addr), user, group)
}

func IsMember_Local(ctx context.Context, cloned gov.Cloned, user User, group Group) bool {
	var userHasGroup, groupHasUser bool
	must.Try(
		func() { userHasGroup = userGroupsKKV.Get(ctx, userGroupsNS, cloned.Tree(), user, group) },
	)
	must.Try(
		func() { groupHasUser = groupUsersKKV.Get(ctx, groupUsersNS, cloned.Tree(), group, user) },
	)
	return userHasGroup && groupHasUser
}

func RemoveMember(ctx context.Context, addr gov.Address, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveMember_StageOnly(ctx, cloned, user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveMember_StageOnly(ctx context.Context, cloned gov.Cloned, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Remove(ctx, userGroupsNS, cloned.Tree(), user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, cloned.Tree(), group, user)

	// log
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "remove_user_from_group",
		Args:   trace.M{"user": user, "group": group},
		Result: nil,
	})

	return git.NewChangeNoResult(fmt.Sprintf("Removed user %v from group %v", user, group), "member_remove_member")
}

func ListUserGroups(ctx context.Context, addr gov.Address, user User) []Group {
	return ListUserGroups_Local(ctx, gov.Clone(ctx, addr), user)
}

func ListUserGroups_Local(ctx context.Context, cloned gov.Cloned, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, cloned.Tree(), user)
}

func ListGroupUsers(ctx context.Context, addr gov.Address, group Group) []User {
	return ListGroupUsers_Local(ctx, gov.Clone(ctx, addr), group)
}

func ListGroupUsers_Local(ctx context.Context, cloned gov.Cloned, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, cloned.Tree(), group)
}
