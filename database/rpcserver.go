package database

import (
	"context"
	"net"
	"sync"

	"github.com/raycastillo3/clickCountApp/pb"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, rpcAddr string) error {
	// create the in-memory database
	clickCountApp := &ClickCountAppDatabase{}

	// create a grpc server (without TLS)
	s := grpc.NewServer()

	// register the proto-specific methods on the grpc server
	pb.RegisterClickCountAppServer(s, clickCountApp)

	// listen on a tcp socket
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}

	// run the server
	return s.Serve(lis)
}

type ClickCountAppDatabase struct {
	pb.ClickCountAppServer

	mu                                     sync.Mutex
	itemClicks, addToCartClicks, buyClicks int64
}

// Call `SetClicks` on the RPC server to update the values.
func (s *ClickCountAppDatabase) SetClicks(ctx context.Context, r *pb.SetClicksRequest) (*pb.SetClicksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.itemClicks = r.ClickCounts.Item
	s.addToCartClicks = r.ClickCounts.AddToCart
	s.buyClicks = r.ClickCounts.Buy
	return &pb.SetClicksResponse{}, nil
}

// Call `GetClicks` on the RPC server to get updated values from the server.
func (s *ClickCountAppDatabase) GetClicks(ctx context.Context, r *pb.GetClicksRequest) (*pb.GetClicksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.GetClicksResponse{
		ClickCounts: &pb.ClickCounts{
			Item:      s.itemClicks,
			AddToCart: s.addToCartClicks,
			Buy:       s.buyClicks,
		},
	}, nil
}
