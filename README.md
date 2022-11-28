# GO_POC_GRPC
Simple GOlang blockchain POC using grpc 

## Run Servers

#### For node at 8085:
```go run server/main.go -port=:8085```


#### For node at 3000:
```go run server/main.go -port=:3000```


#### For node at 50051:
```go run server/main.go -port=:50051```


## Send Transaction

```go run client/main.go -port=:8085 -id=1 -hash=hehe```

