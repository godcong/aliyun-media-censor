package service

import (
	"context"
	"fmt"
	"github.com/godcong/node-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"testing"
	"time"
)

type testManagerServer struct {
}

// Back ...
func (*testManagerServer) Back(ctx context.Context, request *proto.ManagerCallbackRequest) (*proto.ManagerReply, error) {
	log.Printf("%+v", request)
	return &proto.ManagerReply{}, nil
}

// TestGrpcBack_Callback ...
func TestGRPCBack_Callback(t *testing.T) {
	config = LoadConfig("../config.toml")
	StartQueue(context.Background(), 2)
	go func() {
		time.Sleep(5 * time.Second)
		key := "origin/5c35cc6b5ec8a925a4143001/e84976d3567f339635eb0d49cccae72c/0050.mp4"
		stream := NewStreamerWithConfig(Config(), key)
		//stream.Dir, stream.FileName = filepath.Split(key)
		stream.ObjectKey = key
		stream.SetEncrypt(false)
		stream.StreamerCallback = NewBack()
		//stream.SetURI("")
		//stream.FileDest = config.Media.Upload
		//stream.SetSrc(config.Media.Transfer)
		queue.Set(stream.ID, StatusQueuing, 0)
		Push(stream)
	}()

	server := grpc.NewServer()

	_ = syscall.Unlink("/tmp/manager.sock")
	lis, err := net.Listen("unix", "/tmp/manager.sock")
	//port = s.Path

	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	proto.RegisterManagerServiceServer(server, &testManagerServer{})
	// Register reflection service on gRPC server.
	reflection.Register(server)
	log.Printf("Listening and serving TCP on %s\n", "/tmp/manager.sock")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// TestGRPCServer_RemoteDownload ...
func TestGRPCServer_RemoteDownload(t *testing.T) {
	resp, err := http.PostForm("http://localhost:7781/v0/rd", url.Values{"key": []string{"origin/5c35cc6b5ec8a925a4143001/e84976d3567f339635eb0d49cccae72c/0050.mp4"}})
	if err != nil {
		t.Log(err)
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	t.Log(string(bytes), err)
}
