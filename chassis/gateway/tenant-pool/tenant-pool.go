package tenantpool

import (
	"fmt"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/tenant-pool/broadcast"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/tenant-pool/tenant"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/response"
)

type TenantPool struct {
	// map [SessionID] tenant
	pool map[uint64]*tenant.Tenant

	broadcaster *broadcast.Broadcast

	logger logger.Logger
}

func InitTenantPool(ps_config *config.PubsubConfig,
	log logger.Logger) (*TenantPool, error) {
	var ret TenantPool;
	var err error

	ret.pool = make(map[uint64]*tenant.Tenant);
	ret.logger = log;
	ret.broadcaster, err = broadcast.InitBroadcaster(ps_config, log)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			resp, err := ret.broadcaster.Subscribe()
			if err != nil {
				ret.logger.Error(err.Error())
			}
			destination := ret.pool[resp.SessionID]
			if destination != nil {
				destination.SendResponse(resp)
			}
		}
	}()

	return &ret, nil
}

func (pool *TenantPool) NewTenant(tnt *tenant.Tenant) {
	pool.pool[tnt.SessionID] = tnt
}

func (pool *TenantPool) SendResponse(resp *response.UserResponse) {
	var destination *tenant.Tenant
	if destination = pool.pool[resp.SessionID]; destination == nil {
		go pool.broadcaster.Publish(resp)
		return 
	}
	go destination.SendResponse(resp)
}

func (pool *TenantPool) KillTenant(ID uint64) {
	pool.logger.Infor(fmt.Sprintf("Tenant %d exited from pool",ID))
	tnt := pool.pool[ID]
	tnt.Terminate()
	pool.pool[ID] = nil
}
