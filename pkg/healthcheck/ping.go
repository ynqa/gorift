package healthcheck

import (
	"net"
	"time"
)

func Ping(timeout time.Duration) Func {
	return Func(func(req Request) (Report, error) {
		report := Report{
			Available: false,
			LastCheck: time.Now(),
		}
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(req.Address.String(), req.Port.String()), timeout)
		if err != nil {
			return report, err
		}
		conn.Close()
		report.Available = true
		return report, nil
	})
}
