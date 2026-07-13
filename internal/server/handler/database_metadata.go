package handler

import (
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"control-panel-service/pkg/responsewriter"
	"encoding/json"
	"net/http"
	"strings"

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

// // GetAll godoc
// //
// //	@Summary		Get all databases
// //	@Description	Get list of all non-template PostgreSQL databases
// //	@Tags			database-metadata
// //	@Accept			json
// //	@Produce		json
// //	@Success		200	{object}	response.DatabaseMetadataDatabasesResponse
// //	@Failure		500	{object}	response.ErrorResponse
// //	@Router			/api/v1/databases [get]
func (handler *DatabaseMetadataHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("database-metadata-handler")
	traceCtx, span := tracer.Start(r.Context(), "get-databases-handler")
	defer span.End()

	requestID := logger.GetRequestID(r.Context())
	span.SetAttributes(attribute.String("request_id", requestID))

	result, err := handler.databaseMetaDataService.GetAll(traceCtx)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "usecase_error"))
		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)
		return
	}

	responsewriter.WriteJSON(w, http.StatusOK, response.DatabaseMetadataDatabasesResponse{
		Data: result,
	})
}

// // GetTablesByDBNamePost godoc
// //
// //	@Summary		Get tables by database name
// //	@Description	Get list of tables inside a specific database
// //	@Tags			database-metadata
// //	@Accept			json
// //	@Produce		json
// //	@Param			request	body		request.GetTablesRequest	true	"Database name request"
// //	@Success		200		{object}	response.TablesByType
// //	@Failure		400		{object}	response.ErrorResponse
// //	@Failure		500		{object}	response.ErrorResponse
// //	@Router			/api/v1/databases/tables [post]
func (handler *DatabaseMetadataHandler) GetTablesByDBNamePost(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("database-metadata-handler")
	traceCtx, span := tracer.Start(r.Context(), "get-tables-by-database-handler")
	defer span.End()

	requestID := logger.GetRequestID(r.Context())
	span.SetAttributes(attribute.String("request_id", requestID))

	req := new(request.GetTablesRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		span.SetAttributes(attribute.String("error.type", "bad_request"))
		httpError := pkg.ToHTTPError(pkg.ErrBadRequest)
		httpError.WriteError(w)
		return
	}

	dbName := strings.TrimSpace(req.DatabaseName)
	if dbName == "" {
		span.SetAttributes(attribute.String("error.type", "invalid_db_name"))
		httpError := pkg.ToHTTPError(pkg.ErrBadRequest)
		httpError.WriteError(w)
		return
	}

	span.SetAttributes(attribute.String("database.name", dbName))

	result, err := handler.databaseMetaDataService.GetTablesByDBName(traceCtx, dbName)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "usecase_error"), attribute.String("error.message", err.Error()))
		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)
		return
	}

	if result == nil {
		responsewriter.WriteJSON(w, http.StatusOK, response.TablesByType{
			MotherTables: []string{},
			TestTables:   []string{},
		})
	}

	responsewriter.WriteJSON(w, http.StatusOK, response.TablesByType{
		MotherTables: result.MotherTables,
		TestTables:   result.TestTables,
	})

}
