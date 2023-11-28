package webserver

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/raycastillo3/clickCountApp/pb"
	"google.golang.org/grpc"

	"golang.org/x/sync/errgroup"
	"storj.io/common/sync2"
)

func Run(ctx context.Context, httpAddr, rpcAddr, webdir string, cacheInterval time.Duration) error {
	apiServer, err := NewAPIServer(httpAddr, rpcAddr, webdir, cacheInterval)
	if err != nil {
		return err
	}
	defer apiServer.Close()

	var group errgroup.Group
	group.Go(func() error {
		return apiServer.ServeHTTP(ctx)
	})
	group.Go(func() error {
		return apiServer.cache.loop.Run(ctx, apiServer.refreshCache)
	})

	return group.Wait()
}

type Cache struct {
	mu     sync.Mutex
	loop   *sync2.Cycle
	values JSONClicks
}

type APIServer struct {
	grpcConn   *grpc.ClientConn
	clickCount pb.ClickCountAppClient
	webdir     string
	httpAddr   string
	cache      Cache
}

func NewAPIServer(httpAddr, rpcAddr, webdir string, cacheInterval time.Duration) (*APIServer, error) {
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	// make a grpc proto-specific client
	client := pb.NewClickCountAppClient(conn)

	return &APIServer{
		grpcConn:   conn,
		clickCount: client,
		webdir:     webdir,
		httpAddr:   httpAddr,
		cache: Cache{
			loop: sync2.NewCycle(cacheInterval),
		},
	}, nil
}

func (a *APIServer) Close() error {
	return a.grpcConn.Close()
}

func (a *APIServer) refreshCache(ctx context.Context) error {
	fmt.Println("refreshing cache")

	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()
	//SetClicks
	_, err := a.clickCount.SetClicks(ctx, &pb.SetClicksRequest{
		ClickCounts: &pb.ClickCounts{
			Item:      a.cache.values.Item,
			AddToCart: a.cache.values.AddToCart,
			Buy:       a.cache.values.Buy,
		},
	})
	if err != nil {
		return err
	}
	//GetClicks
	getRes, err := a.clickCount.GetClicks(ctx, &pb.GetClicksRequest{})
	if err != nil {
		return err
	}
	a.cache.values.Item = getRes.ClickCounts.Item
	a.cache.values.AddToCart = getRes.ClickCounts.AddToCart
	a.cache.values.Buy = getRes.ClickCounts.Buy

	return nil

}

type JSONClicks struct {
	Item      int64 `json:"itemClicks"`
	AddToCart int64 `json:"addToCartClicks"`
	Buy       int64 `json:"buyClicks"`
}

func (a *APIServer) ServeHTTP(ctx context.Context) error {

	http.HandleFunc("/api/clicks", func(w http.ResponseWriter, r *http.Request) {
		a.getClicksHandler(ctx, w, r)
	})
	http.HandleFunc("/api/clicks/item", func(w http.ResponseWriter, r *http.Request) {
		a.clickCountHandler(ctx, w, r, "item")
	})
	http.HandleFunc("/api/clicks/addToCart", func(w http.ResponseWriter, r *http.Request) {
		a.clickCountHandler(ctx, w, r, "addToCart")
	})
	http.HandleFunc("/api/clicks/buy", func(w http.ResponseWriter, r *http.Request) {
		a.clickCountHandler(ctx, w, r, "buy")
	})

	// static files
	fs := http.FileServer(http.Dir(filepath.Join(a.webdir, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// home page
	http.HandleFunc("/", a.indexHandler)

	return http.ListenAndServe(a.httpAddr, nil)
}
