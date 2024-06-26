package bureau

import (
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/proto/member"
)

const BureauTopic = "bureau"

type Request struct {
	Transfer *TransferRequest `json:"transfer"`
}

type Requests []Request

type TransferRequest struct {
	FromUser member.User `json:"from_user"`
	ToUser   member.User `json:"to_user"`
	Amount   float64     `json:"amount"`
}

type FetchedRequest struct {
	User     member.User      `json:"requesting_user"`
	Address  id.PublicAddress `json:"requesting_address"`
	Requests Requests         `json:"requests"`
}

type FetchedRequests []FetchedRequest
