package main

// [{Name:PostService Package:post GoName:PostService Methods:[{Name:GetPost GoName:GetPost In:0x1400042ed80 Out:0x1400042eea0 Options:url:"/api/post"}] GPRCPort:50051}]
import (
	"log"
	"fmt"
	"net/http"
	"time"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/pb/post"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	r := mux.NewRouter()
	PostServiceHandler, err := GetPostServiceHandler()
	if err != nil {
		log.Fatalf("failed to get PostService handler: %v", err)
	}
	r.HandleFunc("/api/post", PostServiceHandler.GetPostHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("serving on :8000")
	log.Fatal(srv.ListenAndServe())
}

type postService struct {
	client post.PostServiceClient
}

func GetPostServiceHandler() (*postService, error) {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial 50051: %w", err)
	}
	return &postService{
		client: post.NewPostServiceClient(conn),
	}, nil
}
func makeGetPostRequest() *post.Empty {
	return &post.Empty{}
}
func (sh *postService) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	in := makeGetPostRequest()
	resp, err := sh.client.GetPost(ctx, in)
	if err != nil {
		writeError(
			w, fmt.Errorf("failed to call GetPost: %w", err),
			http.StatusInternalServerError,
		)
		return
	}
	b, err := protojson.Marshal(resp)
	if err != nil {
		writeError(
			w, fmt.Errorf("failed to convert response to json: %w", err),
			http.StatusInternalServerError,
		)
		return
	}
	w.Write(b)
}
func writeError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
