package manager

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rancher/longhorn-manager/manager/pb"
)

type GRPCNotifier struct {
	ip           string
	port         int
	callbackChan chan Event
	done         chan struct{}
}

func NewGRPCNotifier(ip string, port int) *GRPCNotifier {
	return &GRPCNotifier{
		done: make(chan struct{}),
		ip:   ip,
		port: port,
	}
}

func (r *GRPCNotifier) GetAddress() string {
	return r.ip + ":" + strconv.Itoa(r.port)
}

func (r *GRPCNotifier) NodeNotifyRPC(ctx context.Context, req *pb.NodeNotifyRequest) (*pb.NodeNotifyResponse, error) {
	r.callbackChan <- Event{
		Type:       EventType(req.Event),
		VolumeName: req.VolumeName,
	}
	return &pb.NodeNotifyResponse{
		Result: "success",
	}, nil
}

func (r *GRPCNotifier) Start(callbackChan chan Event) (err error) {
	l, err := net.Listen("tcp", r.GetAddress())
	if err != nil {
		return errors.Wrap(err, "fail to start GRPC server")
	}
	r.callbackChan = callbackChan

	s := grpc.NewServer()
	pb.RegisterManagerServer(s, r)
	reflection.Register(s)
	go func() {
		if err := s.Serve(l); err != nil {
			logrus.Errorf("fail to serve GRPC server: %v", err)
		}
	}()
	go func() {
		<-r.done
		s.GracefulStop()
	}()
	return nil
}

func (r *GRPCNotifier) Stop() {
	close(r.done)
}

func (r *GRPCNotifier) GetPort() int {
	return r.port
}

func (r *GRPCNotifier) NodeNotify(address string, event *Event) error {
	//FIXME insecure
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "fail to connect to %v", address)
	}
	defer conn.Close()
	c := pb.NewManagerClient(conn)

	resp, err := c.NodeNotifyRPC(context.Background(), &pb.NodeNotifyRequest{
		Event:      string(event.Type),
		VolumeName: event.VolumeName,
	})
	if err != nil {
		return errors.Wrapf(err, "fail to notify %v", address)
	}
	if resp.Result != "success" {
		return fmt.Errorf("fail to notify %v: %+v", address, resp)
	}
	return nil
}
