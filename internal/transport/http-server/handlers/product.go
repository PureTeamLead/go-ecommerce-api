package handlers

import (
	"errors"
	"eshop/internal/infrastructure/errs"
	"eshop/internal/transport/http-server/dto"
	"eshop/internal/transport/http-server/dto/requests"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) AddProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Add Product"))

	var req requests.AddProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	logging.Debug("adding to database product...")
	id, err := h.prrs.AddProduct(req.Name, req.Category, req.Price)
	if err != nil {
		logging.Error("failed add to db new product", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Server error"))
	}

	logging.Info("Operation success", zap.Any("id", id))
	return e.JSON(http.StatusCreated, dto.NewOkReponse("product added with id", id))
}

func (h *Handler) DeleteProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Delete Product"))

	var req requests.DeleteGetProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	id, err := h.prrs.DeleteProduct(req.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNoProductFound) {
			logging.Info("no id in products db", zap.Any("id", req.ID))
			return e.JSON(http.StatusNotFound, dto.NewErrorResponse(errs.ErrNoProductFound, "No product found"))
		}

		logging.Error("refused to communicate with db", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Server error"))
	}

	logging.Info("Operation success", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("product deleted with id", id))
}

func (h *Handler) UpdateProductInfo(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Update Product"))

	var req requests.UpdateProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	updatedProduct, err := h.prrs.UpdateProduct(req.ID, req.Name, req.Category, req.Price)
	if err != nil {
		if errors.Is(err, errs.ErrNoProductFound) {
			logging.Info("no id in products db", zap.Any("id", req.ID))
			return e.JSON(http.StatusNotFound, dto.NewErrorResponse(errs.ErrNoProductFound, "No product found"))
		}

		logging.Error("DB fail", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Failed to update product info"))
	}

	logging.Info("Operation success", zap.Any("updated product", updatedProduct))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully updated product", updatedProduct))
}

func (h *Handler) GetProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Get Product"))

	var req requests.DeleteGetProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	product, err := h.prrs.GetProduct(req.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNoProductFound) {
			logging.Info("no id in products db", zap.Any("id", req.ID))
			return e.JSON(http.StatusNotFound, dto.NewErrorResponse(errs.ErrNoProductFound, "No product found"))
		}

		logging.Error("DB fail", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Failed to fetch product"))
	}

	logging.Info("Operation success", zap.Any("fetched product", product))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully fetched product", product))
}

func (h *Handler) GetAllProducts(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Get All Products"))

	products, err := h.prrs.GetAllProducts()
	if err != nil {
		logging.Error("DB fail", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Failed to fetch all products"))
	}

	logging.Info("Operation success", zap.Any("fetched all products", products))
	return e.JSON(http.StatusOK, dto.NewOkReponse("Successfully fetched all products", products))
}
