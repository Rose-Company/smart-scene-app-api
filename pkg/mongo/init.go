package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodb struct {
	prefix       string
	host         string
	user         string
	pass         string
	databaseName string
	client       *mongo.Client
	db           *mongo.Database
}

//func NewMongoDB(prefix, hostEnv, userEnv, passEnv, databaseName string) *Mongodb {
//	host, _ := os.LookupEnv(hostEnv)
//	user, _ := os.LookupEnv(userEnv)
//	pass, _ := os.LookupEnv(passEnv)
//
//	return &Mongodb{
//		prefix:       prefix,
//		host:         host,
//		user:         user,
//		pass:         pass,
//		databaseName: databaseName,
//	}
//}

func (t *Mongodb) Get() interface{} {
	return t.db
}

func (t *Mongodb) Run() (err error) {
	err = t.validate()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(`Mongodb://%s:%s@%s/`, t.user, t.pass, t.host)

	opt := options.ClientOptions{}
	opt.ApplyURI(url)
	opt.SetTLSConfig(&tls.Config{})

	if false {
		srvMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, e *event.CommandStartedEvent) {
				fmt.Println(fmt.Sprintf("[QUERY] %v", e.Command))
			},
			Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
				fmt.Println(fmt.Sprintf("[%vms] %v", e.DurationNanos/1000000, e.CommandFinishedEvent))
			},
			Failed: func(_ context.Context, e *event.CommandFailedEvent) {
				fmt.Println(e.Failure)
			},
		}
		opt.SetMonitor(srvMonitor)
	}

	err = opt.Validate()
	if err != nil {
		return err
	}

	t.client, err = mongo.Connect(context.TODO(), &opt)
	if err != nil {
		return err
	}

	t.db = t.client.Database(t.databaseName)
	return nil
}

func (t *Mongodb) validate() error {
	if t.prefix == "" {
		return fmt.Errorf("prefix is empty")
	}

	if t.host == "" {
		return fmt.Errorf("host is empty")
	}

	if t.user == "" {
		return fmt.Errorf("user is empty")
	}

	if t.pass == "" {
		return fmt.Errorf("pass is empty")
	}

	if t.databaseName == "" {
		return fmt.Errorf("databaseName is empty")
	}

	return nil
}

func (t *Mongodb) GetPrefix() string {
	return t.prefix
}

func (t *Mongodb) Stop() <-chan bool {
	stop := make(chan bool)
	go func() {
		stop <- true
	}()
	return stop
}
