package mongo_test

import (
	"testing"

	. "web-layout/utils/mongo"
)

func TestConnect(t *testing.T) {
	conf := &Cfg{
		Size:          10,
		Addr:          []string{"localhost:27017"},
		Username:      "test",
		Password:      "test",
		DBName:        "test",
		AuthMechanism: "SCRAM-SHA-1",
	}

	db, err := conf.Connect()
	if err != nil || db == nil {
		t.Error("err = ", err)
		return
	}
}
