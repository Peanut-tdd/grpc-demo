
## grpc demo

#### 生成pb文件到指定目录
```aiignore
protoc -I./proto --go_out=./genproto --go_opt  paths=source_relative --go-grpc_out=./genproto  --go-grpc_opt=paths=source_relative ./proto/user/*.proto
```