package django_session

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type PGXTestSuite struct {
	suite.Suite
	store *PGXSession
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(PGXTestSuite))
}

func (s *PGXTestSuite) SetupTest() {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		url = "postgres://tests:tests@localhost:5432/tests"
	}
	dbpool, err := pgxpool.Connect(context.Background(), url)
	s.Require().Nil(err)
	s.store = &PGXSession{Pool: dbpool}
}

func (s *PGXTestSuite) TearDownSuite() {
	s.store.Pool.Close()
}

// func (s *PGXTestSuite) TestFound() {
// 	key := "63krc5u7su98m6qr6zyal7v7c8mzxcx2"
// 	dest := BaseSession{}
// 	err := s.store.Fetch(context.Background(), key, &dest)
// 	s.Require().Nil(err)
// 	s.Require().Equal("1", dest.UserID)
// }
