package routes

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/di"
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var Config = &configs.Env{
	MongoURI: "mongodb://localhost:27017",
	DBName:   fmt.Sprintf("test-%d", time.Now().UnixMilli()),
}
var RepositoryContext = di.NewRepositoryContext(Config)
var ServiceContext = di.NewServiceContext(RepositoryContext)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	dbClient := configs.ConnectDB(*Config)
	_ = dbClient.Database(Config.DBName).Drop(context.Background())

	os.Exit(exitCode)
}
