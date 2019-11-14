package algorithm

import (
	"math/rand"

	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/metrics"
	"github.com/gorift/gorift/pkg/server"
)

type p2c struct{}

func NewP2C() Algorithm {
	return &p2c{}
}

func (p *p2c) Pick(members []*server.Member) (*server.Member, error) {
	n := len(members)
	if n < 1 {
		return nil, xerrors.New("there are no members")
	} else if n == 1 {
		return members[0], nil
	} else if n == 2 {
		return pickWithTotalPicked(members[0], members[1])
	} else {
		m1 := members[rand.Intn(n)]
		m2 := members[rand.Intn(n)]
		return pickWithTotalPicked(m1, m2)
	}
}

func pickWithTotalPicked(m1, m2 *server.Member) (*server.Member, error) {
	m1metric, err := m1.GetMetrics(metrics.TotalPickedLabel)
	if err != nil {
		return nil, err
	}
	m2metric, err := m2.GetMetrics(metrics.TotalPickedLabel)
	if err != nil {
		return nil, err
	}
	m1load := m1metric.(int)
	m2load := m2metric.(int)

	res := m1
	if m1load > m2load {
		res = m2
	}
	return res, nil
}
