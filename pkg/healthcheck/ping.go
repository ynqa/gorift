package healthcheck

import (
	"net"
	"time"
)

func Ping(timeout time.Duration) HealthcheckFn {
	return HealthcheckFn(func(req HealthcheckRequest) (HealthcheckReport, error) {
		resp := HealthcheckReport{
			Available: false,
			LastCheck: time.Now(),
		}
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(req.Address.String(), req.Port.String()), timeout)
		if err != nil {
			return resp, err
		}
		conn.Close()
		resp.Available = true
		return resp, nil
	})
}
