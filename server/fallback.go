package server

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhump/protoreflect/desc"
	gdyn "github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
)

// FallbackServer ...
type FallbackServer struct {
	server  http.Server
	gServer *grpc.Server
	impls   map[string]interface{}
	stub    *gdyn.Stub
	backend string
}

// NewFallbackServer ...
func NewFallbackServer(port, backend string, gServer *grpc.Server) *FallbackServer {
	return &FallbackServer{
		server: http.Server{
			Addr: ":" + port,
		},
		gServer: gServer,
		impls:   make(map[string]interface{}),
		backend: backend,
	}
}

// Start ...
func (f *FallbackServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/$rpc/{service:[.a-zA-Z0-9]+}/{method:[a-zA-Z]+}", f.handler).Headers("Content-Type", "application/x-protobuf")
	f.server.Handler = r

	go func() {
		log.Println("Fallback server listening on port:", f.server.Addr)
		err := f.server.ListenAndServe()
		if err != nil {
			log.Println("Error in fallback server while listening:", err)
		}
	}()
}

// Shutdown ...
func (f *FallbackServer) Shutdown() {
	err := f.server.Shutdown(context.Background())
	if err != nil {
		log.Println("Error shutting down fallback server:", err)
	}
}

// RegisterService ...
func (f *FallbackServer) RegisterService(name string, impl interface{}) {
	f.impls[name] = impl
}

func (f *FallbackServer) handler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	service := v["service"]
	method := v["method"]

	if info, ok := f.gServer.GetServiceInfo()[service]; !ok || !containsMethod(info, method) {
		w.WriteHeader(404)
		return
	}

	// f.invoke(r.Body, w, f.impls[service], method)
}

func (f *FallbackServer) invoke(in io.Reader, out io.Writer, md *desc.MethodDescriptor) error {
	if f.stub == nil {
		cc, err := grpc.Dial(f.backend, grpc.WithInsecure())
		if err != nil {
			log.Printf("Error dialing backend on %s: %v\n", f.backend, err)
			return err
		}

		s := gdyn.NewStub(cc)
		f.stub = &s
	}

	// f.stub.InvokeRpc(context.Background(), md, )

	return nil
}

func containsMethod(serv grpc.ServiceInfo, method string) bool {
	for _, minf := range serv.Methods {
		if minf.Name == method {
			return true
		}
	}

	return false
}
