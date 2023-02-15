#### Grpc Server

Since grpc and http are the ports of the same application, to keep consistency of generated hash,
grpc server will start background hash generation function.

Http server can retrieve a valid hash via grpc connection

```bash
go run . serve-grpc --grpc-addr :5105 --hash-ttl 5m
```

#### Http Server
```bash
go run . serve-http --http-addr :8081 --grpc-addr :5105
```

#### Tests
```bash
go test ./... --race
```

Commands also available in .run folder as runnable configurations