package watcher

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"

	authenticator "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/module/auth"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/tenant-pool/tenant"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/protoc"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/request"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TenantWatcher struct {
	protoc.UnimplementedStreamServiceServer
	auth *authenticator.Authenticator

	grpcServer *grpc.Server

	tenantChannel  chan (*tenant.Tenant)
	requestChannel chan (*request.UserRequest)
	closeChannel   chan (uint64)
}

func InitTenantWatcher(log logger.Logger,
	conf *config.SecurityConfig,
	watcherconf *config.GatewayConfig,
) (*TenantWatcher, error) {
	var ret TenantWatcher
	ret.auth = authenticator.InitAuthenticator(conf)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", watcherconf.WatcherPort))
	if err != nil {
		log.Fatal(err.Error())
	}

	ret.grpcServer = grpc.NewServer()
	protoc.RegisterStreamServiceServer(ret.grpcServer, &ret)
	go ret.grpcServer.Serve(lis)
	return &ret, nil
}

func (watcher *TenantWatcher) StreamRequest(client protoc.StreamService_StreamRequestServer) error {
	var ok bool
	var err error
	var usr *user.User
	var headers metadata.MD
	headers, ok = metadata.FromIncomingContext(client.Context())
	if !ok {
		return errors.New("Unauthorized")
	}

	token := headers["Authorization"]
	usr, err = watcher.auth.ValidateToken(token[0], "User")
	if err != nil {
		return nil
	}

	SessionID := rand.Uint64()
	new := tenant.NewTenant(SessionID, usr)
	watcher.tenantChannel <- new

	defer func() {
		watcher.closeChannel <- SessionID
	}()

	go func() {
		for {
			response := new.ListenonResponse();
			transresp := response.ProtocTranslate();
			client.Send(&transresp);
		}
	}()

	for {
		in, err := client.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}

		request := request.ReverseTranslate(in)
		request.Headers["SessionID"] = fmt.Sprintf("%d", SessionID)
		request.Headers["UserID"] = fmt.Sprintf("%d", usr.ID)
		request.Headers["Username"] = usr.UserName
		request.Headers["RequestID"] = fmt.Sprintf("%d", request.ID)
		watcher.requestChannel <- request
	}
}

func (watcher *TenantWatcher) WaitClose() uint64 {
	return <-watcher.closeChannel
}

func (watcher *TenantWatcher) WaitTenant() *tenant.Tenant {
	return <-watcher.tenantChannel
}

func (watcher *TenantWatcher) WaitRequest() *request.UserRequest {
	return <-watcher.requestChannel
}
