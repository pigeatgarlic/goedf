package tenantwatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	authenticator "github.com/pigeatgarlic/goedf/chassis/gateway/module/auth"
	registrator "github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/microservice"
	"github.com/pigeatgarlic/goedf/models/request-response/protoc"
	"github.com/pigeatgarlic/goedf/models/request-response/response"
	"github.com/pigeatgarlic/goedf/models/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type HelperWatcher struct {
	protoc.UnimplementedHelperServiceServer
	registrator registrator.ServiceRegistrator

	auth *authenticator.Authenticator

	grpcServer *grpc.Server
}

func InitHelperEndpoint(log logger.Logger,
	conf *config.SecurityConfig,
	helper *config.GatewayConfig,
	registrator registrator.ServiceRegistrator) (*HelperWatcher, error) {
	var ret HelperWatcher
	ret.auth = authenticator.InitAuthenticator(conf)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", helper.HelperPort))
	if err != nil {
		log.Fatal(err.Error())
	}

	ret.grpcServer = grpc.NewServer()
	protoc.RegisterHelperServiceServer(ret.grpcServer, &ret)
	go ret.grpcServer.Serve(lis)
	return &ret, nil
}

func (watcher *HelperWatcher) HelperRequest(ctx context.Context, req *protoc.UserRequest) (*protoc.UserResponse, error) {
	var err error
	var resp *response.UserResponse
	var usr *user.User

	peer, _ := peer.FromContext(ctx)
	req.Headers["UserIPAddress"] = peer.Addr.String()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("Missing metadata")
	}
	token := md["Authorization"]
	if token[0] == "" {
		return &protoc.UserResponse{
			Id:    req.Id,
			Error: "Empty token",
			Data:  map[string]string{},
		}, nil
	}

	usr, err = watcher.auth.ValidateToken(token[0], "Admin")
	if err != nil {
		return &protoc.UserResponse{
			Id:    req.Id,
			Error: err.Error(),
			Data:  map[string]string{},
		}, nil
	}
	resp.Data, err = watcher.authorizedEndpoint(req.Target, usr, req.Headers, req.Data)

	ret := resp.ProtocTranslate()
	return &ret, err
}

func (watcher *HelperWatcher) authorizedEndpoint(target string,
	user *user.User,
	headers map[string]string,
	data map[string]string) (ret map[string]string, err error) {
	switch target {
	case "registerFeature":
		var tag map[string]string
		err = json.Unmarshal([]byte(data["Tags"]), &tag)
		if err != nil {
			return map[string]string{}, err
		}

		var endpoints []int
		err = json.Unmarshal([]byte(data["Endpoints"]), &endpoints)
		if err != nil {
			return map[string]string{}, err
		}

		feature := microservice.Feature{
			ID:          rand.Int(),
			Name:        data["Name"],
			Tags:        tag,
			Authority:   "gRPC Gateway",
			EndpointIDs: endpoints,
		}

		feature.Tags["RegisterAt"] = time.Now().Format(time.RFC3339)
		err = watcher.registrator.RegisterFeature(&feature)
		if err != nil {
			return map[string]string{}, err
		}

		return map[string]string{
			"FeatureID": fmt.Sprintf("%d", feature.ID),
		}, nil
	}

	return nil, errors.New("unknown endpoint")
}
