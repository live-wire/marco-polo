package client

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	pb "github.com/live-wire/marco-polo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// MarcoPoloClient is the type of MarcoPolo client object
type MarcoPoloClient struct {
	Src            string
	serviceClient  pb.MarcoPoloServiceClient
	grpcClientConn *grpc.ClientConn
	l              sync.Mutex
}

// InitMarcoPoloClient initializes and returns an instance of MarcoPolo client
func InitMarcoPoloClient(address string, src string) (*MarcoPoloClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		src = "default"
	}
	c := pb.NewMarcoPoloServiceClient(conn)
	return &MarcoPoloClient{Src: src, serviceClient: c, grpcClientConn: conn}, nil
}

// Cleanup cleans dangling MarcoPolo client references
func (x *MarcoPoloClient) Cleanup() {
	x.l.Lock()
	defer x.l.Unlock()
	x.serviceClient = nil
	x.grpcClientConn.Close()
}

// Parses IP address from request
// Used the implementation from https://blog.golang.org/context/userip/userip.go
func fromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	log.Println("ip", ip)
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	log.Println("userIp", userIP)
	return userIP, nil
}

func fromHeader(req *http.Request) (string, error) {
	forward := req.Header.Get("X-Forwarded-For")
	log.Println("forward", forward)
	forwards := strings.Split(forward, ",")
	if len(forward) < 1 {
		log.Println("No X-Forwarded-For header found")
		ip, err := fromRequest(req)
		if err != nil {
			return "", err
		}
		return ip.String(), nil
	}
	return forwards[0], nil
}

// Consume sends relevant information to MarcoPolo Service
func (x *MarcoPoloClient) consumeFromRequest(req *http.Request) {
	x.l.Lock()
	defer x.l.Unlock()
	if x.serviceClient == nil {
		log.Println("gRPC connection is closed.")
		return
	}
	ip, err := fromHeader(req)
	if err != nil {
		log.Println("ip extraction error", err)
		return
	}
	x.Consume(ip, nil)
}

// Consume sends an IP address point to MarcoPolo server
func (x *MarcoPoloClient) Consume(ip string, tags map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	msg := &pb.Message{Ip: ip, Src: x.Src, Tags: tags}
	_, err := x.serviceClient.Consume(ctx, msg)
	if err != nil {
		log.Printf("Could not consume: %v", err)
		return
	}
	log.Printf("Consumed: %v", msg)
}

// MarcoPoloDecorator decorates HandleFunc for quick integration
func (x *MarcoPoloClient) MarcoPoloDecorator(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		x.consumeFromRequest(req)
		f(w, req)
	}
}
