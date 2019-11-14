package algorithm

import (
	"github.com/gorift/gorift/pkg/server"
)

type Algorithm interface {
	Pick([]*server.Member) (*server.Member, error)
}
