package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"codegen/internal/entity"
	"codegen/internal/repository/dto"
	"codegen/internal/service"
	"codegen/pkg"
)

type ProductHandler interface {
	Search(c echo.Context) error
	GetList(c echo.Context) error
	GetSingle(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type productHandler struct {
	srv service.ProductService
}

func NewProductHandler(srv service.ProductService) ProductHandler {
	return &productHandler{srv: srv}
}

// Search			godoc
//
// @Tags			Products
// @Param			name	query		string false "Name"
// @Param			page	query		int false "Page"
// @Param			limit	query		int false "Limit (max=100)"
// @Success		200	{object}	pkg.ProductsResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products/search [get]
func (h *productHandler) Search(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	// filter := dto.ProductFilter{Skip: 0, Take: 10}
	// filter.Id, _ = getIdAsUUID(c)
	// filter.Name = c.Param("name")
	// filter.Skip = getParamAsInt16("skip")
	// filter.Take = getParamAsInt16("take")

	req := pkg.ProductFilterReq{}
	if err := c.Bind(&req); err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Info().Interface("req", req).Msg("")

	// page, err := getParamAsInt16(c, "page", 1)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("")
	// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// }

	// limit, err := getParamAsInt16(c, "limit", 10)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("")
	// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// }
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	filter := &dto.ProductFilter{
		Id:   req.Id,
		Name: req.Name,
		Skip: ((req.Page - 1) * req.Limit),
		Take: req.Limit,
	}

	logger.Info().Interface("filter", filter).Msg("")

	totalCount, products, err := h.srv.SearchProducts(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := &pkg.ProductsResponse{Count: totalCount, Products: products}
	return errors.WithStack(c.JSON(http.StatusOK, res))
}

// GetList		godoc
//
// @Tags			Products
// @Success		200	{object}	pkg.ProductsResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products [get]
func (h *productHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	products, err := h.srv.GetAllProducts(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := &pkg.ProductsResponse{Count: len(products), Products: products}
	return errors.WithStack(c.JSON(http.StatusOK, res))
}

// GetSingle	godoc
//
// @Tags			Products
// @Param			id	path			string true "Product Identifier" format(uuid)
// @Success		200	{object}	pkg.ProductResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products/{id} [get]
func (h *productHandler) GetSingle(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	id, err := getIdAsUUID(c)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	product, err := h.srv.GetProductById(ctx, id)

	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := &pkg.ProductResponse{Product: product}
	return errors.WithStack(c.JSON(http.StatusOK, res))
}

// Create			godoc
//
// @Tags			Products
// @Accept		json
// @Produce		json
// @Param			request body 	pkg.ProductCreateUpdateReq true "query params"
// @Success		201	{object}	pkg.IdResult[uuid.UUID]
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products [post]
func (h *productHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	req := pkg.ProductCreateUpdateReq{}
	if err := shouldBindJSON(c, &req); err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	product := &entity.Product{}
	product.BrandId = req.BrandId
	product.Name = req.Name
	product.Sku = &req.Sku
	product.Summary = &req.Summary
	product.Storyline = &req.Storyline
	product.StockQuantity = req.StockQuantity
	product.Price = req.Price

	// defaults
	product.Id = uuid.New()
	product.CreatedBy = getCurrentUserId(c)
	product.CreatedAt = time.Now().UTC()

	err := h.srv.CreateProduct(ctx, product)
	if err != nil {
		logger.Error().Err(err).Msg("product insert error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := &pkg.IdResult[uuid.UUID]{Id: product.Id}
	return errors.WithStack(c.JSON(http.StatusCreated, res))
}

// Update			godoc
//
// @Tags			Products
// @Accept		json
// @Produce		json
// @Param			id	path			string true "Product Identifier" format(uuid)
// @Param			request body	pkg.ProductCreateUpdateReq true "query params"
// @Success		204
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products/{id} [put]
func (h *productHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	id, err := getIdAsUUID(c)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := pkg.ProductCreateUpdateReq{}
	if err := shouldBindJSON(c, &req); err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	product, err := h.srv.GetProductById(ctx, id)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if product == nil {
		logger.Error().Msg("product not found")
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}

	product.BrandId = req.BrandId
	product.Name = req.Name
	product.Sku = &req.Sku
	product.Summary = &req.Summary
	product.Storyline = &req.Storyline
	product.StockQuantity = req.StockQuantity
	product.Price = req.Price

	// defaults
	updatedBy := getCurrentUserId(c)
	updatedAt := time.Now().UTC()

	product.UpdatedBy = &updatedBy
	product.UpdatedAt = &updatedAt

	if err := h.srv.UpdateProduct(ctx, product); err != nil {
		logger.Error().Err(err).Msg("product update error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return errors.WithStack(c.NoContent(http.StatusNoContent))
}

// Delete			godoc
//
// @Tags			Products
// @Param			id	path			string true "Product Identifier" format(uuid)
// @Success		204
// @Failure		400	{object}	pkg.ErrorResponse
// @Router		/api/products/{id} [delete]
func (h *productHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	id, err := getIdAsUUID(c)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = h.srv.DeleteProduct(ctx, id); err != nil {
		logger.Error().Err(err).Msg("product delete error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return errors.WithStack(c.NoContent(http.StatusNoContent))
}
