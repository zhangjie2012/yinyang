package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(host string, port int) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	s := &Server{
		HttpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: router,
		},
	}

	s.RegisterRouter(router)

	return s
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go s.Run(ctx, wg)
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// graceful close
	defer func() {
		waiting := 5 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), waiting)
		defer cancel()

		if err := s.HttpServer.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failure, error=%s", err)
		}

		log.Infof("server exiting")
	}()

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("server run on: %s", s.HttpServer.Addr)

	select {
	case <-ctx.Done():
		return
	}
}
