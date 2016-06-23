package tag

import (
	"errors"
	"log"
	"time"

	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Request struct {
	User   User
	Result *UserTag
	Ended  chan error
}

var (
	requestPool chan *Request
	timeout     time.Duration

	ErrorRquestPoolFull = errors.New("tag rpc call error for the pool has not enough connections")
	ErrorRpcCallTimeout = errors.New("rpc call time out")
)

func newWorker(addr string, inputChan <-chan *Request, dialOptions []grpc.DialOption) {
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		clog.Log.Fatalf("did not connect: %v", err)
	}
	client := NewTagsClient(conn)
	for request := range inputChan {
		if conn == nil {
			conn, err := grpc.Dial(addr, dialOptions...)
			if err != nil {
				log.Println(conn)
			}
			client = NewTagsClient(conn)
		}
		//TODO context
		r, err := client.GetUserTags(context.Background(), &request.User)
		if err != nil {
			//TODO:
			request.Ended <- err
		} else {
			request.Result = r
			request.Ended <- nil
		}
	}
}

func Init(addr string, poolNum int, callTimeout time.Duration) {
	requestPool = make(chan *Request, poolNum*2)
	dialOptions := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(callTimeout * 5)}
	timeout = callTimeout
	go func() {
		for i := 0; i < poolNum; i++ {
			go newWorker(addr, requestPool, dialOptions)
		}
	}()
}

func RequestTag(ctx context.Context, user User) (*UserTag, error) {
	request := Request{User: user, Ended: make(chan error, 1)}
	select {
	case requestPool <- &request:
	default:
		return nil, ErrorRquestPoolFull
	}
	select {
	case err := <-request.Ended:
		return request.Result, err
	case <-time.After(timeout):
		return nil, ErrorRpcCallTimeout
	}
	return nil, errors.New("should not run to here")
}
