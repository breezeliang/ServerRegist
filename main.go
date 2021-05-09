package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is version 3"))
}

// 关闭http
func sayBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bye bye"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/bye", sayBye)

	stop := make(chan error)
	eg, ctx := errgroup.WithContext(context.Background())
	for i:=0; i <=1; i++ {
		eg.Go(func() error {
			s := http.Server{
				Addr: fmt.Sprintf(":808%d", i),
				Handler: mux,
			}
			select {
			case <-stop:
				_ = s.Shutdown(ctx)
			}
			err := s.ListenAndServe()
			if err != nil {
				stop <- err
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		os.Exit(0)
	}



}


