# Metrics

```go
type Metric interface {
	Add(interface{}) error
	Get() interface{}
}
```
