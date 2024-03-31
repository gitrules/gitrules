package member

import (
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/proto/kv"
)

var (
	membersNS = proto.RootNS.Append("members")

	usersNS = membersNS.Append("users")
	usersKV = kv.KV[User, UserProfile]{}

	groupsNS = membersNS.Append("groups")
	groupsKV = kv.KV[Group, form.None]{}

	userGroupsNS  = membersNS.Append("user_groups")
	userGroupsKKV = kv.KKV[User, Group, bool]{}

	groupUsersNS  = membersNS.Append("group_users")
	groupUsersKKV = kv.KKV[Group, User, bool]{}
)

type UserProfile struct {
	ID            id.ID            `json:"id"`
	PublicAddress id.PublicAddress `json:"public_address"`
}

func UserAccountID(user User) account.AccountID {
	return account.AccountIDFromLine(account.Pair("user", string(user)))
}
