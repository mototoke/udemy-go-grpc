## 実行コマンド

- 27.gRPCのコンパイル
```
protoc -I. --go_out=. --go-grpc_out=. proto/*.proto

go mod tidy
```

- 31.サーバーストリーミングRPCのスキーマ定義
```
protoc -I. --go_out=. --go-grpc_out=. proto/*.proto


```

