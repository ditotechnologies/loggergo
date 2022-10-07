package loggergo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"sync"
	"time"
)

type internalEndpointResolver struct {
	endpointResolverMutex sync.Mutex
	endpointFetchTime     time.Time
	endpoint              *timestreamwrite.Endpoint
}

var endpointResolver = &internalEndpointResolver{}

func (er *internalEndpointResolver) fetchEndpointDescription(ctx context.Context, service *timestreamwrite.TimestreamWrite) (*timestreamwrite.Endpoint, error) {
	er.endpointResolverMutex.Lock()
	defer er.endpointResolverMutex.Unlock()

	if er.endpoint != nil {
		dMinutes := int64(time.Now().Sub(er.endpointFetchTime).Minutes())
		if dMinutes < *er.endpoint.CachePeriodInMinutes {
			// can use the cached version.
			return er.endpoint, nil
		} else {
			er.endpoint = nil
		}
	}

	er.endpointFetchTime = time.Now()

	// else, need to actually describe the endpoints
	input := timestreamwrite.DescribeEndpointsInput{}
	output, err := service.DescribeEndpoints(&input)
	if err != nil {
		return nil, err
	}

	if len(output.Endpoints) == 0 {
		return nil, fmt.Errorf("could not find endpoints")
	}

	er.endpoint = output.Endpoints[0]

	return er.endpoint, nil
}
