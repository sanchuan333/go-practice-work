package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var fakeHandler = FakeH{}

type FakeH struct{}

func (h FakeH) ServeHTTP(http.ResponseWriter, *http.Request) {
	fmt.Println("handler a request")
}

type App struct {
	eg             *errgroup.Group
	stop           chan struct{}
	isStopped      bool
	httpServices   []http.Server
	signalListener []func(ctx context.Context) error
}

func (a *App) New(ctx context.Context, services []http.Server, sListener []func(ctx context.Context) error) context.Context {
	a.eg, ctx = errgroup.WithContext(ctx)
	a.httpServices = services
	a.isStopped = false
	a.stop = make(chan struct{})
	a.signalListener = sListener
	return ctx
}

func (a *App) Run(ctx context.Context) {
	a.startHttpServer(ctx)
	a.startSignalListen(ctx)

	a.eg.Go(func() error {
		select {
		case <-ctx.Done():
			if !a.isStopped {
				close(a.stop)
				a.isStopped = true
			}

		}
		return nil
	})
	if err := a.eg.Wait(); err != nil {
		fmt.Println("server down", err)
		if !a.isStopped {
			close(a.stop)
			a.isStopped = true
		}
	}
}

func (a *App) startHttpServer(ctx context.Context) {
	for _, s := range a.httpServices {
		server := s
		a.eg.Go(func() error {
			go func() {
				<-a.stop
				fmt.Println("server stop", server.Addr)
				_ = server.Shutdown(ctx)
			}()
			fmt.Println("start server", server.Addr)
			return server.ListenAndServe()
		})
	}
}

func (a *App) startSignalListen(ctx context.Context) {
	for _, s := range a.signalListener {
		f := s
		a.eg.Go(func() error {
			return f(ctx)
		})
	}
}

const serverOneAddr = "127.0.0.1:3001"
const serverTwoAddr = "127.0.0.1:3002"

func listenSignal(ctx context.Context) error {
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println("start listen signal")
	select {
	case s := <-c:
		fmt.Println("get signal:", s)
		return errors.New("signal-" + s.String())
	case <-ctx.Done():
		fmt.Println("close listen signal")
		return nil
	}
}

func fakeListenS(ctx context.Context) error {
	time.Sleep(10 * time.Second)
	fmt.Println("get fake signal stop")
	return errors.New("fake-signal-stop")
	//return nil
}

func main() {
	app := App{}
	services := []http.Server{{Addr: serverOneAddr, Handler: fakeHandler}, {Addr: serverTwoAddr, Handler: fakeHandler}}
	listeners := []func(ctx context.Context) error{listenSignal, fakeListenS}
	ctx := app.New(context.Background(), services, listeners)
	app.Run(ctx)
}
