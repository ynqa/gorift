package algorithm

import (
	"math/rand"

	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/server"
)

type random struct {
}

func NewRandom() Algorithm {
	return &random{}
}

func (r *random) Pick(members []*server.Member) (*server.Member, error) {
	n := len(members)
	if n < 1 {
		return nil, xerrors.New("there are no members")
	} else if n == 1 {
		return members[0], nil
	}
	return members[rand.Intn(n)], nil
}
