package handler

import (
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewDatabaseMetadataHandler(databaseMetaDataService DatabaseMetadata) *DatabaseMetadataHandler {
	return &DatabaseMetadataHandler{
		databaseMetaDataService: databaseMetaDataService,
	}
}

type DatabaseMetadataHandler struct {
	databaseMetaDataService DatabaseMetadata
}

// GetAll godoc
//
//	@Summary		Get all databases
//	@Description	Get list of all non-template PostgreSQL databases
//	@Tags			database-metadata
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.DatabaseMetadataDatabasesResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/databases [get]
func (handler *DatabaseMetadataHandler) GetAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("database-metadata-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get-databases-handler")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		result, err := handler.databaseMetaDataService.GetAll(traceCtx)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "usecase_error"))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.DatabaseMetadataDatabasesResponse{
			Data: result,
		})
	}
}

// GetTablesByDBNamePost godoc
//
//	@Summary		Get tables by database name
//	@Description	Get list of tables inside a specific database
//	@Tags			database-metadata
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.GetTablesRequest	true	"Database name request"
//	@Success		200		{object}	response.TablesByType
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/databases/tables [post]
func (handler *DatabaseMetadataHandler) GetTablesByDBNamePost() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("database-metadata-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get-tables-by-database-handler")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.GetTablesRequest)
		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		dbName := strings.TrimSpace(req.DatabaseName)
		if dbName == "" {
			span.SetAttributes(attribute.String("error.type", "invalid_db_name"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("database.name", dbName))

		result, err := handler.databaseMetaDataService.GetTablesByDBName(traceCtx, dbName)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "usecase_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		if result == nil {
			return ctx.Status(http.StatusOK).JSON(&response.TablesByType{
				MotherTables: []string{},
				TestTables:   []string{},
			})
		}

		return ctx.Status(http.StatusOK).JSON(&response.TablesByType{
			MotherTables: result.MotherTables,
			TestTables:   result.TestTables,
		})
	}
}
