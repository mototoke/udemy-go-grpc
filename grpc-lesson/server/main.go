package main

import (
	"bytes"
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	grpc_middlerware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ファイルサービスサーバー

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := "D:/Work/udemy_go_grpc/workspace/grpc-lesson/strage/" + filename

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return status.Error(codes.NotFound, "file was not found")
	}

	file, err := os.Open(path)

	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)

		// ファイル読み込みされなくなったまたは読み込み完了した時
		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res) // streamを投げる
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}

	return nil // streamはEOFを返してresponseが終了
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := "D:/Work/udemy_go_grpc/workspace/grpc-lesson/strage"

	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			// フォルダでなければファイル名を追加
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}

	return res, nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data)
	}
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("recieved data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("recieved %vbytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("request data: %+v", req)

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}
		log.Printf("response data: %+v", resp)

		return resp, nil
	}
}

func authorize(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	if token != "test-token" {
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}
	return ctx, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middlerware.ChainUnaryServer(
			myLogging(),
			grpc_auth.UnaryServerInterceptor(authorize))))
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
