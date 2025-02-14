package handlers

import (
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/errs"
	"eshop/internal/transport/http-server/dto"
	"eshop/internal/transport/http-server/dto/requests"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type ProductHandler interface {
	AddProduct(e echo.Context) error
	DeleteProduct(e echo.Context) error
	UpdateProductInfo(e echo.Context) error
	GetProduct(e echo.Context) error
}

//newProduct := entities.NewProduct(product.Price, product.Name, product.Category)

func (h *HandlerStruct) AddProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Add Product"))

	var req requests.AddUpdateProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	newProduct := entities.NewProduct(req.Price, req.Name, req.Category)
	logging.Info("new product created", zap.Any("id", newProduct.ID))

	logging.Debug("adding to database new created product...")
	id, err := h.prrs.AddProduct(newProduct)
	if err != nil {
		logging.Error("failed add to db new product", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Server error"))
	}

	logging.Info("Product successfully added", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("product added"))
}

func (h *HandlerStruct) DeleteProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Delete Product"))

	var req requests.DeleteGetProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	id, err := h.prrs.DeleteProduct(req.ID)
	if err != nil {
		logging.Error("refused to communicate with db", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, dto.NewErrorResponse(errs.ErrDB, "Server error"))
	}

	logging.Info("product successfully deleted", zap.Any("id", id))
	return e.JSON(http.StatusOK, dto.NewOkReponse("product deleted"))
}

func (h *HandlerStruct) UpdateProductInfo(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Update Product"))

	var req requests.AddUpdateProduct
	if err := e.Bind(&req); err != nil {
		logging.Error("failed to bind request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, dto.NewErrorResponse(errs.ErrBadRequest, "Got an unexpected fields in request"))
	}

	updatedProduct := entities.UpdateProduct(req.ID, req.)

	h.prrs.UpdateProduct()
}

func (h *HandlerStruct) GetProduct(e echo.Context) error {
	logging := h.logger.With(zap.String("Use Case", "Get Product"))

}
