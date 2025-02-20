# restful

- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go get -u github.com/gin-gonic/gin
- go get -u github.com/gin-contrib/cors
- go get -u gorm.io/driver/postgres
- go get -u gorm.io/gorm
- go get -u github.com/google/uuid
- go get github.com/valkey-io/valkey-go
- go get github.com/polygon-io/client-go/rest
- go get github.com/polygon-io/client-go/rest/models

```shell
protoc -I ./protos --go_out ./pkg/pbm --go_opt paths=source_relative --go-grpc_out=./pkg/pbm --go-grpc_opt=paths=source_relative ./protos/*.proto
```