
# ProtoBus

ProtoBus borrows simple intra-process Message Bus from [Grafana][Grafana], but has a much better performance.

### Publish Example

```go
func ExampleInProcBus_Publish() {
	bus := protobus.ProvideBus()

	protobus.AddEventListener(bus, func(ctx context.Context, query *emptypb.Empty) error {
		fmt.Printf("event received: %T", query)
		return nil
	})

	bus.Publish(context.Background(), &emptypb.Empty{})
	// Output: event received: *emptypb.Empty
}
```

### Error Handling Example

```go
func main() {
	bus := protobus.ProvideBus()
	var ErrEmpty = errors.New("empty")

	protobus.AddEventListener(bus, func(ctx context.Context, query *emptypb.Empty) error {
		if query == nil {
			return ErrEmpty
		}
		fmt.Printf("event received: %T", query)
		return nil
	})
	if err := bus.Publish(context.Background(), (*emptypb.Empty)(nil)); err != nil {
		fmt.Println("error received:", err)
		fmt.Println("error is ErrEmpty?", errors.Is(err, ErrEmpty))
	}
	bus.Publish(context.Background(), &emptypb.Empty{})
	// Output:
	// error received: empty
	// error is ErrEmpty? true
	// event received: *emptypb.Empty
}
```

### Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/Myriad-Dreamin/protobus
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkEventCtxPublish-24         	 2662341	       458.2 ns/op	     176 B/op	       6 allocs/op
BenchmarkProtoEventCtxPublish-24    	34715662	        32.10 ns/op	      48 B/op	       1 allocs/op
PASS
	github.com/Myriad-Dreamin/protobus	coverage: 100% of statements
ok  	github.com/Myriad-Dreamin/protobus	2.830s
```


[Grafana]: https://github.com/grafana/grafana/blob/b7fc837c359ddc4101ba5bccc26e03d9a41d6967/pkg/bus/bus.go
