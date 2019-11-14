# Balancer

## Algorithm

```go
type Algorithm interface {
	Pick([]*server.Member) (*server.Member, error)
}
```
