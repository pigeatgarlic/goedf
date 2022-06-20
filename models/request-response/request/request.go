package request

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/protoc"
	proto "github.com/golang/protobuf/proto"
)

type UserRequest struct {
	ID      int
	Target  string
	Headers map[string]string

	Data map[string]string
}

func ReverseTranslate(req *protoc.UserRequest) *UserRequest {
	return &UserRequest{
		ID: int(req.Id),
		Target: req.Target,
		Headers: req.Headers,
		Data: req.Data,
	}
}

func (req *UserRequest) protocTranslate() protoc.UserRequest {
	return protoc.UserRequest{
		Id: int64(req.ID),
		Target: req.Target,
		Headers: req.Headers,
		Data: req.Data,
	}
}

func (req *UserRequest) ToProtobytes() ([]byte ,error ){
	proto_req := req.protocTranslate();
	return proto.Marshal(&proto_req);
}

func FromProtobytes(data *[]byte) (*UserRequest,error ){
	var req protoc.UserRequest;
	err :=proto.Unmarshal(*data,&req);
	if err != nil {
		return nil,err;
	}
	proto_req := ReverseTranslate(&req);
	return proto_req,nil;
}