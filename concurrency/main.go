package concurrency

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

var fakeHandler FakeH = FakeH{}

type FakeH struct{}

func (h FakeH) ServeHTTP(http.ResponseWriter, *http.Request) {
	fmt.Println("handler a request")
}

func startServer(addr string, handler http.Handler) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.ListenAndServe()
}

func serverGroup(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return nil
	})

}
