package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// grpc.WithInSecure()は非暗号化で非推奨
	// 通常はSSLのオプションを利用する
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect:%v", err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	// callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	// name.txtをダウンロード
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		// サーバーがレスポンスを返すたびにresで受け取る
		res, err := stream.Recv()
		// サーバーのストリームが終了
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}
