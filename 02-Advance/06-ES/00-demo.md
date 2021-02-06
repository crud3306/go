

```golang
import (
	"context"
	"guardian/util"
	"sync"

	"gopkg.in/olivere/elastic.v6"

	"reflect"
)

var (
	esClient *elastic.Client
	esLock   sync.Mutex
)

func GetESClient() *elastic.Client {
	esLock.Lock()
	defer esLock.Unlock()

	if esClient == nil {
		esClient = InitESClient()
	}

	return esClient
}

func InitESClient() *elastic.Client {
	esConfig := library.GetESConfig()

	var err error
	esClient, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(esConfig.Host...),
		elastic.SetBasicAuth(esConfig.User, esConfig.Password),
	)
	if err != nil {
		panic(err)
	}
	return esClient
}


```