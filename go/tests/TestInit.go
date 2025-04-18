package tests

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

var globals common.IResources

func init() {
	config := &types.SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  "tests"}
	secure, err := common.LoadSecurityProvider("security.so")
	if err != nil {
		panic(err)
	}
	globals = resources.NewResources(registry.NewRegistry(), secure, nil, Log, nil, nil, config, nil)
}
