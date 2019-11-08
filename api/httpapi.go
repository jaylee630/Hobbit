package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"golang-admin-basic-master/utils/log"
)

type (
	Server struct {
		*log.Logger
		Host      string
		Port      int
		AdminPort int
		admin     *http.Server
		http      *http.Server
		logLevel  string
		Context   context.Context
	}
)

func (s *Server) IsStopping() bool {
	select {
	case _, ok := <-s.Context.Done():
		return !ok
	default:
		return false
	}
}

func (s *Server) StartHTTP() {
	s.Infof("The http server is listen on %s:%d", s.Host, s.Port)

	go func() {
		if e := s.http.ListenAndServe(); e != nil && !s.IsStopping() {
			s.Panic("failed to listen and serve http server.", zap.Error(e))
		}
		s.Info("The http server is exited.")
	}()
	s.Info("The http server is started.")
}

func (s *Server) StartAdmin() {
	s.Infof("The admin server is listen on %s:%d", s.Host, s.AdminPort)

	go func() {
		if e := s.admin.ListenAndServe(); e != nil && !s.IsStopping() {
			s.Panic("failed to listen and serve admin http server.", zap.Error(e))
		}
		s.Info("The admin server is exited.")
	}()

	s.Info("The admin server is started.")
}

func (s *Server) Start() {

	s.StartHTTP()
	if s.admin != nil {
		s.StartAdmin()
	}
}

func (s *Server) stopAdminServer() {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if err := s.admin.Shutdown(ctx); err != nil {
		s.Warn("failed to stop admin http server.", zap.Error(err))
	}
}

func (s *Server) stopHTTPServer() {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		s.Warn("failed to stop admin http server.", zap.Error(err))
	}
}

func (s *Server) Stop() {
	s.stopHTTPServer()
	if s.AdminPort != 0 {
		s.stopAdminServer()
	}
	s.Info("All API Server has been stopped. ")
}

func (s *Server) initAdminServer() {

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/pprof", pprof.Handler("pprof"))

	s.admin = &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.AdminPort),
	}
}

func (s *Server) Init(ctx context.Context, handler http.Handler) {

	// init context
	s.Context = ctx

	s.http = &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.Port),
	}

	if s.AdminPort != 0 {
		s.initAdminServer()
	}
}

func NewHTTPServer(logger *log.Logger, host string, port, adminPort int, devMode bool, logLevel string) *Server {

	svc := &Server{
		Logger:    logger,
		Host:      host,
		Port:      port,
		AdminPort: adminPort,
		logLevel:  logLevel,
	}
	return svc
}
