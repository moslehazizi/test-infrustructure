package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/usecase"
)

func NewTestCategoryHandler(cfg *config.Config, testCategoryService usecase.TestCategoryService) *TestCategoryHandler {
	return &TestCategoryHandler{
		cfg,
		testCategoryService,
	}
}

type TestCategoryHandler struct {
	cfg                 *config.Config
	testCategoryService usecase.TestCategoryService
}
