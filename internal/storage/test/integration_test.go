//go:build integration

package test

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/storage"
)

type IntegrationTestSuite struct {
	suite.Suite
	closeDB    func()
	repository *storage.Repository
}

func (suite *IntegrationTestSuite) SetupTest() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	suite.repository, suite.closeDB, err = storage.NewRepository(context.Background(), logger.NewDisabled())
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *IntegrationTestSuite) TearDownTest() {
	suite.closeDB()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
