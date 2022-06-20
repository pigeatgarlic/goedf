package response

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/protoc"
	proto "github.com/golang/protobuf/proto"
)

type UserResponse struct {
	ID int
	SessionID uint64
	
	Error string
	Data map[string]string
}


func reverseTranslate(req *protoc.UserResponse) *UserResponse {
	return &UserResponse{
		ID: int(req.Id),
		Error: req.Error,
		Data: req.Data,
	}
}

func (req *UserResponse) ProtocTranslate() protoc.UserResponse {
	return protoc.UserResponse{
		Id: int64(req.ID),
		Error: req.Error,
		Data: req.Data,
	}
}

func (req *UserResponse) ToProtobytes() ([]byte ,error ){
	proto_req := req.ProtocTranslate();
	return proto.Marshal(&proto_req);
}

func FromProtobytes(data *[]byte) (*UserResponse,error ){
	var req protoc.UserResponse;
	err :=proto.Unmarshal(*data,&req);
	if err != nil {
		return nil,err;
	}
	proto_req := reverseTranslate(&req);
	return proto_req,nil;
}