package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/evg555/antibrutforce/api/pb"
	"github.com/evg555/antibrutforce/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

func (a *APITestSuite) SetupTest() {
	a.truncateTables()
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

func (a *APITestSuite) truncateTables() {
	query := `TRUNCATE TABLE whitelist; TRUNCATE TABLE blacklist;`
	_, err := a.db.Exec(query)
	a.Require().NoError(err)
}

func (a *APITestSuite) TearDownSuite() {
	a.truncateTables()

	err := a.db.Close()
	a.Require().NoError(err)

	err = a.conn.Close()
	a.Require().NoError(err)
}

func (a *APITestSuite) TestAuth() {
	netAddress := "172.1.1.0/24"
	ipAddress := "172.1.1.1"
	login := "login"
	password := "password"

	authReq := &pb.AuthRequest{}
	res, err := a.client.Auth(context.Background(), authReq)
	a.Error(err)

	st, ok := status.FromError(err)
	a.True(ok)
	a.Equal(codes.InvalidArgument, st.Code())

	// In blacklist
	req := &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.AddIpBlacklist(context.Background(), req)
	a.NoError(err)

	authReq = &pb.AuthRequest{
		Login:    login,
		Password: password,
		Ip:       ipAddress,
	}
	res, err = a.client.Auth(context.Background(), authReq)
	a.NoError(err)
	a.NotNil(res)
	a.False(res.Ok)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpBlacklist(context.Background(), req)
	a.NoError(err)

	// In whitelist (priority)
	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.AddIpBlacklist(context.Background(), req)
	a.NoError(err)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.AddIpWhitelist(context.Background(), req)
	a.NoError(err)

	authReq = &pb.AuthRequest{
		Login:    login,
		Password: password,
		Ip:       ipAddress,
	}
	res, err = a.client.Auth(context.Background(), authReq)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpBlacklist(context.Background(), req)
	a.NoError(err)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpWhitelist(context.Background(), req)
	a.NoError(err)

	// Brutforce on login
	attempts := 10

	for attempts > 0 {
		authReq = &pb.AuthRequest{
			Login:    login,
			Password: gofakeit.Password(true, true, true, true, true, 10),
			Ip:       gofakeit.IPv4Address(),
		}

		res, err = a.client.Auth(context.Background(), authReq)
		a.NoError(err)
		a.NotNil(res)
		a.Require().True(res.Ok, "Attempts: %d", 10-attempts)

		attempts--
	}

	res, err = a.client.Auth(context.Background(), authReq)
	a.NoError(err)
	a.NotNil(res)
	a.False(res.Ok)
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

	// Invalid subnet
	req := &pb.IpRequest{}
	res, err := a.client.AddIpWhitelist(context.Background(), req)
	a.Error(err)

	st, ok := status.FromError(err)
	a.True(ok)
	a.Equal(codes.InvalidArgument, st.Code())

	// Positive case
	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.AddIpWhitelist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	var rows1 []storage.Subnet
	err = a.db.Select(&rows1, `SELECT * FROM whitelist WHERE subnet=$1 LIMIT 1`, netAddress)
	a.NoError(err)
	a.Equal(rows1[0].Address, netAddress)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpWhitelist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	var rows2 []storage.Subnet
	err = a.db.Select(&rows2, `SELECT * FROM whitelist WHERE subnet=$1 LIMIT 1`, netAddress)
	a.NoError(err)
	a.Empty(rows2)
}

func (a *APITestSuite) TestIpBlackList() {
	netAddress := "172.1.1.0/24"

	// Invalid subnet
	req := &pb.IpRequest{}
	res, err := a.client.AddIpBlacklist(context.Background(), req)
	a.Error(err)

	st, ok := status.FromError(err)
	a.True(ok)
	a.Equal(codes.InvalidArgument, st.Code())

	// Positive case
	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.AddIpBlacklist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	var rows1 []storage.Subnet
	err = a.db.Select(&rows1, `SELECT * FROM blacklist WHERE subnet=$1 LIMIT 1`, netAddress)
	a.NoError(err)
	a.Equal(rows1[0].Address, netAddress)

	req = &pb.IpRequest{NetAddress: netAddress}
	res, err = a.client.DeleteIpBlacklist(context.Background(), req)
	a.NoError(err)
	a.NotNil(res)
	a.True(res.Ok)

	var rows2 []storage.Subnet
	err = a.db.Select(&rows2, `SELECT * FROM blacklist WHERE subnet=$1 LIMIT 1`, netAddress)
	a.NoError(err)
	a.Empty(rows2)
}
