package metric_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/caarlos0/env"

	. "web-layout/utils/aliyun/metric"
)

type testClient struct {
	key string
}

func (t testClient) GetCustomMetricMetricList() cms.PutCustomMetricMetricList {
	num := rand.Intn(100)
	list := cms.PutCustomMetricMetricList{
		MetricName: t.key,
		Dimensions: fmt.Sprintf(`{"dimensionality":"%s"}`, t.key),
		Values:     fmt.Sprintf(`{"value":%d}`, num),
	}
	return list
}

func TestCmsGroup(t *testing.T) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		t.Error(err)
	}

	client1 := testClient{key: "test1"}
	client2 := testClient{key: "test2"}
	client3 := testClient{key: "test3"}
	client4 := testClient{key: "test4"}

	g, err := CmsGroup(cfg, client1, client2, client3, client4)
	if err != nil {
		t.Error(err)
		return
	}

	g.Start()

	time.Sleep(5 * time.Minute)
}
