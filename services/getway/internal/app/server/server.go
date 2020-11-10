package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/sirupsen/logrus"
	"gitlab.com/taskProvider/services/getway/internal/app/httpmiddle"
	"gitlab.com/taskProvider/services/getway/internal/app/services"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type server struct {
	mux          *chi.Mux
	logger       *logrus.Logger
	endpoint     *string
	middle       httpmiddle.Service
	service      services.Handler
}

// NewServer ...
func NewServer() error {
	s := &server{
		logger:       logrus.New(),
		mux:          chi.NewRouter(),
	}

	if err := s.configureMiddleware(); err != nil {
		return err
	}

	return nil
}

// MiddleWare
func (s *server) configureMiddleware() error {

	// middlewars
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Use(middleware.Timeout(60 * time.Second))

	grpcGwMux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(s.httpResponseModifier),
	)

	// Инициализируем сервис Login
	initUserConn, err := s.service.Init(grpcGwMux); 
	defer initUserConn.Close()
	if err != nil {
		return err
	}
	
	// Роуты V1 user
	s.mux.Route("/api/v1/user", func(route chi.Router) {
		route.Handle("/*", grpcGwMux)
	})

	return http.ListenAndServe(":8081", s.mux)
}

// Status code modificator
func (s *server) httpResponseModifier(ctx context.Context, w http.ResponseWriter, p protoreflect.ProtoMessage) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)

		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	}

	return nil
}