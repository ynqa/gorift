# Healthcheck

```go
type HealthcheckFn func(HealthcheckRequest) (HealthcheckReport, error)

type HealthcheckRequest struct {
	Address server.Address
	Port    server.Port
}

type HealthcheckReport struct {
	Available bool
	LastCheck time.Time
}
```
