package filter

import (
	"github.com/gorift/gorift/pkg/server"
)

type FilterFn func([]*server.Member) []*server.Member

func Availables() FilterFn {
	return FilterFn(func(members []*server.Member) []*server.Member {
		res := make([]*server.Member, 0)
		for _, v := range members {
			if v.HealthStatus.Available {
				res = append(res, v)
			}
		}
		return res
	})
}
