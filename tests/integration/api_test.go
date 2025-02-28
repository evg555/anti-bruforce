package integration

import (
	"context"
	"fmt"
	"github.com/evg555/antibrutforce/api/pb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type APITestSuite struct {
	suite.Suite
	db     *sqlx.DB
	conn   *grpc.ClientConn
	client pb.AppServiceClient
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (a *APITestSuite) SetupSuite() {
	a.setupDB()
	a.setupGRPSClient()
}

func (a *APITestSuite) setupDB() {
	var err error

	dsn := fmt.Sprintf("postgresql://dbuser:dbpass@localhost:5432/dbname?sslmode=disable")

	db, err := sqlx.Open("pgx", dsn)
	a.Require().NoError(err)

	a.db = db
}

func (a *APITestSuite) setupGRPSClient() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	a.Require().NoError(err)

	a.conn = conn
	a.client = pb.NewAppServiceClient(conn)
}

func (a *APITestSuite) TearDownSuite() {
	err := a.db.Close()
	a.Require().NoError(err)

	err = a.conn.Close()
	a.Require().NoError(err)
}

func (a *APITestSuite) TestAuth() {
	req := &pb.AuthRequest{
		Login:    "login",
		Password: "pass",
		Ip:       "172.1.1.1",
	}
	res, err := a.client.Auth(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	// Проверяем, что данные сохранились в БД
	//var user User
	//err = db.First(&user, res.Id).Error
	//assert.NoError(t, err)
	//assert.Equal(t, "John Doe", user.Name)
	//assert.Equal(t, 30, user.Age)
}

func (a *APITestSuite) TestBucketReset() {
	req := &pb.BucketResetRequest{
		Password: "pass",
		Ip:       "172.1.1.1",
	}
	res, err := a.client.BucketReset(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	// Проверяем, что данные сохранились в БД
	//var user User
	//err = db.First(&user, res.Id).Error
	//assert.NoError(t, err)
	//assert.Equal(t, "John Doe", user.Name)
	//assert.Equal(t, 30, user.Age)
}

func (a *APITestSuite) TestIpWhitelist() {
	netAddress := "172.1.1.0/24"

	req := &pb.IpRequest{NetAddress: netAddress}
	res, err := a.client.AddIpWhitelist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	// Проверяем, что данные сохранились в БД
	//var user User
	//err = db.First(&user, res.Id).Error
	//assert.NoError(t, err)
	//assert.Equal(t, "John Doe", user.Name)
	//assert.Equal(t, 30, user.Age)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpWhitelist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)
}

func (a *APITestSuite) TestIpBlackList() {
	netAddress := "172.1.1.0/24"

	req := &pb.IpRequest{NetAddress: netAddress}
	res, err := a.client.AddIpBlacklist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	// Проверяем, что данные сохранились в БД
	//var user User
	//err = db.First(&user, res.Id).Error
	//assert.NoError(t, err)
	//assert.Equal(t, "John Doe", user.Name)
	//assert.Equal(t, 30, user.Age)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpBlacklist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)
}
