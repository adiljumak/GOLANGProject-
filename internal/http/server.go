package http

import (
	"Lecture6/internal/store"
	"context"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	Address     string
	store       store.Store
	cache       *lru.TwoQueueCache
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {

	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	categoryResource := CategoryResource{s.store, s.cache}
	r.Mount("/categories", categoryResource.Routes())
	return r
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)
	log.Println("Working...")
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
		return
	}

	log.Println("[HTTP] Precessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTerminator() {
	<-s.idleConnsCh
}
