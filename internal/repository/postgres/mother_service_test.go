package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMotherServiceRepository_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		responseDelayDuration := 100
		randomDelayMin := 50
		randomDelayMax := 150
		serviceAddress := "http://service.example.com"

		motherService := &entity.MotherService{
			Model:                    gorm.Model{CreatedAt: now, UpdatedAt: now},
			Name:                     "mother1",
			ExceptionRate:            0.1,
			ResponseDelayRate:        0.2,
			ResponseDelayDuration:    &responseDelayDuration,
			RandomResponseDelayMin:   &randomDelayMin,
			RandomResponseDelayMax:   &randomDelayMax,
			ProvisioningStatus:       entity.ProvisioningStatusPending,
			ServiceDeploymentAddress: &serviceAddress,
			DatabaseName:             "test_db",
			DatabaseTableName:        "test_table",
			KafkaFactorialTopic:      "factorial",
			KafkaLiveFeedTopic:       "live_feed",
			StoppedAt:                &now,
			RestartedAt:              &now,
			StartedAt:                &now,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name","exception_rate","response_delay_rate","response_delay_duration","random_response_delay_min","random_response_delay_max","provisioning_status","service_deployment_address","database_name","database_table_name","kafka_livefeed_topic","kafka_factorial_topic","stopped_at","restarted_at","started_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING "id"`)).
			WithArgs(
				now, now, nil,
				motherService.Name,
				motherService.ExceptionRate,
				motherService.ResponseDelayRate,
				motherService.ResponseDelayDuration,
				motherService.RandomResponseDelayMin,
				motherService.RandomResponseDelayMax,
				motherService.ProvisioningStatus,
				motherService.ServiceDeploymentAddress,
				motherService.DatabaseName,
				motherService.DatabaseTableName,
				motherService.KafkaLiveFeedTopic,
				motherService.KafkaFactorialTopic,
				motherService.StoppedAt,
				motherService.RestartedAt,
				motherService.StartedAt,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.Create(context.Background(), motherService)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), motherService.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewMotherServiceRepository(db)
		now := time.Now()

		motherService := &entity.MotherService{
			Model:               gorm.Model{CreatedAt: now, UpdatedAt: now},
			Name:                "mother1",
			ExceptionRate:       0.0,
			ResponseDelayRate:   0.0,
			ProvisioningStatus:  entity.ProvisioningStatusPending,
			DatabaseName:        "test_db",
			DatabaseTableName:   "test_table",
			KafkaLiveFeedTopic:  "live_feed",
			KafkaFactorialTopic: "factorial",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "mother_services" ("created_at","updated_at","deleted_at","name","exception_rate","response_delay_rate","response_delay_duration","random_response_delay_min","random_response_delay_max","provisioning_status","service_deployment_address","database_name","database_table_name","kafka_livefeed_topic","kafka_factorial_topic","stopped_at","restarted_at","started_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING "id"`)).
			WithArgs(
				now, now, nil,
				motherService.Name,
				motherService.ExceptionRate,
				motherService.ResponseDelayRate,
				motherService.ResponseDelayDuration,
				motherService.RandomResponseDelayMin,
				motherService.RandomResponseDelayMax,
				motherService.ProvisioningStatus,
				motherService.ServiceDeploymentAddress,
				motherService.DatabaseName,
				motherService.DatabaseTableName,
				motherService.KafkaLiveFeedTopic,
				motherService.KafkaFactorialTopic,
				motherService.StoppedAt,
				motherService.RestartedAt,
				motherService.StartedAt,
			).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		err = repo.Create(context.Background(), motherService)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create mother service record")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMotherServiceRepository_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		responseDelayDuration := 100
		randomDelayMin := 50
		randomDelayMax := 150
		serviceAddress := "http://service.example.com"

		expectedMotherService := &entity.MotherService{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:                     "mother1",
			ExceptionRate:            0.1,
			ResponseDelayRate:        0.2,
			ResponseDelayDuration:    &responseDelayDuration,
			RandomResponseDelayMin:   &randomDelayMin,
			RandomResponseDelayMax:   &randomDelayMax,
			ProvisioningStatus:       entity.ProvisioningStatusProvisioned,
			ServiceDeploymentAddress: &serviceAddress,
			DatabaseName:             "test_db",
			DatabaseTableName:        "test_table",
			KafkaLiveFeedTopic:       "live_feed",
			KafkaFactorialTopic:      "factorial",
			StoppedAt:                &now,
			RestartedAt:              &now,
			StartedAt:                &now,
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"provisioning_status", "service_deployment_address",
				"database_name", "database_table_name", "kafka_livefeed_topic", "kafka_factorial_topic",
				"stopped_at", "restarted_at", "started_at",
			}).
				AddRow(
					expectedMotherService.ID,
					expectedMotherService.CreatedAt,
					expectedMotherService.UpdatedAt,
					nil,
					expectedMotherService.Name,
					expectedMotherService.ExceptionRate,
					expectedMotherService.ResponseDelayRate,
					expectedMotherService.ResponseDelayDuration,
					expectedMotherService.RandomResponseDelayMin,
					expectedMotherService.RandomResponseDelayMax,
					expectedMotherService.ProvisioningStatus,
					expectedMotherService.ServiceDeploymentAddress,
					expectedMotherService.DatabaseName,
					expectedMotherService.DatabaseTableName,
					expectedMotherService.KafkaLiveFeedTopic,
					expectedMotherService.KafkaFactorialTopic,
					expectedMotherService.StoppedAt,
					expectedMotherService.RestartedAt,
					expectedMotherService.StartedAt,
				))

		result, err := repo.GetByID(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMotherService.ID, result.ID)
		assert.Equal(t, expectedMotherService.Name, result.Name)
		assert.Equal(t, expectedMotherService.ExceptionRate, result.ExceptionRate)
		assert.Equal(t, expectedMotherService.ResponseDelayRate, result.ResponseDelayRate)
		assert.Equal(t, expectedMotherService.ProvisioningStatus, result.ProvisioningStatus)
		assert.Equal(t, expectedMotherService.DatabaseName, result.DatabaseName)
		assert.Equal(t, expectedMotherService.DatabaseTableName, result.DatabaseTableName)
		assert.Equal(t, expectedMotherService.KafkaLiveFeedTopic, result.KafkaLiveFeedTopic)
		assert.Equal(t, expectedMotherService.KafkaFactorialTopic, result.KafkaFactorialTopic)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(999), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByID(context.Background(), uint64(999))
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."id" = $1 AND "mother_services"."deleted_at" IS NULL ORDER BY "mother_services"."id" LIMIT $2`)).
			WithArgs(uint64(1), 1).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetByID(context.Background(), uint64(1))
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMotherServiceRepository_GetAll(t *testing.T) {
	t.Run("success case with empty result", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"provisioning_status", "service_deployment_address",
				"database_name", "database_table_name", "kafka_livefeed_topic", "kafka_factorial_topic",
				"stopped_at", "restarted_at", "started_at",
			}))

		result, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC`)).
			WillReturnError(errors.New("database connection failed"))

		result, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get mother service records")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success case with order by id desc", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)
		repo := NewMotherServiceRepository(db)

		now := time.Now()
		serviceAddress1 := "http://service1.example.com"
		serviceAddress2 := "http://service2.example.com"

		expectedMotherServices := []*entity.MotherService{
			{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:                     "mother2",
				ExceptionRate:            0.0,
				ResponseDelayRate:        0.0,
				ProvisioningStatus:       entity.ProvisioningStatusPending,
				ServiceDeploymentAddress: &serviceAddress2,
				DatabaseName:             "test_db2",
				DatabaseTableName:        "test_table2",
				KafkaLiveFeedTopic:       "live_feed",
				KafkaFactorialTopic:      "factorial",
			},
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:                     "mother1",
				ExceptionRate:            0.1,
				ResponseDelayRate:        0.2,
				ProvisioningStatus:       entity.ProvisioningStatusProvisioned,
				ServiceDeploymentAddress: &serviceAddress1,
				DatabaseName:             "test_db1",
				DatabaseTableName:        "test_table1",
				KafkaLiveFeedTopic:       "live_feed",
				KafkaFactorialTopic:      "factorial",
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "mother_services" WHERE "mother_services"."deleted_at" IS NULL ORDER BY id DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
				"exception_rate", "response_delay_rate", "response_delay_duration",
				"random_response_delay_min", "random_response_delay_max",
				"provisioning_status", "service_deployment_address",
				"database_name", "database_table_name", "kafka_livefeed_topic", "kafka_factorial_topic",
				"stopped_at", "restarted_at", "started_at",
			}).
				AddRow(
					expectedMotherServices[0].ID,
					expectedMotherServices[0].CreatedAt,
					expectedMotherServices[0].UpdatedAt,
					nil,
					expectedMotherServices[0].Name,
					expectedMotherServices[0].ExceptionRate,
					expectedMotherServices[0].ResponseDelayRate,
					expectedMotherServices[0].ResponseDelayDuration,
					expectedMotherServices[0].RandomResponseDelayMin,
					expectedMotherServices[0].RandomResponseDelayMax,
					expectedMotherServices[0].ProvisioningStatus,
					expectedMotherServices[0].ServiceDeploymentAddress,
					expectedMotherServices[0].DatabaseName,
					expectedMotherServices[0].DatabaseTableName,
					expectedMotherServices[0].KafkaLiveFeedTopic,
					expectedMotherServices[0].KafkaLiveFeedTopic,
					expectedMotherServices[0].StoppedAt,
					expectedMotherServices[0].RestartedAt,
					expectedMotherServices[0].StartedAt,
				).
				AddRow(
					expectedMotherServices[1].ID,
					expectedMotherServices[1].CreatedAt,
					expectedMotherServices[1].UpdatedAt,
					nil,
					expectedMotherServices[1].Name,
					expectedMotherServices[1].ExceptionRate,
					expectedMotherServices[1].ResponseDelayRate,
					expectedMotherServices[1].ResponseDelayDuration,
					expectedMotherServices[1].RandomResponseDelayMin,
					expectedMotherServices[1].RandomResponseDelayMax,
					expectedMotherServices[1].ProvisioningStatus,
					expectedMotherServices[1].ServiceDeploymentAddress,
					expectedMotherServices[1].DatabaseName,
					expectedMotherServices[1].DatabaseTableName,
					expectedMotherServices[1].KafkaLiveFeedTopic,
					expectedMotherServices[1].KafkaLiveFeedTopic,
					expectedMotherServices[1].StoppedAt,
					expectedMotherServices[1].RestartedAt,
					expectedMotherServices[1].StartedAt,
				))

		result, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedMotherServices[0].ID, result[0].ID)
		assert.Equal(t, expectedMotherServices[0].Name, result[0].Name)
		assert.Equal(t, expectedMotherServices[0].DatabaseName, result[0].DatabaseName)
		assert.Equal(t, expectedMotherServices[1].ID, result[1].ID)
		assert.Equal(t, expectedMotherServices[1].Name, result[1].Name)
		assert.Equal(t, expectedMotherServices[1].DatabaseName, result[1].DatabaseName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
