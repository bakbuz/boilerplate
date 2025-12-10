package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codegen/pkg"
	"codegen/utils/random"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Client interface {
	GetList() (pkg.ProductsResponse, error)
	GetSingle(uuid.UUID) (pkg.ProductResponse, error)
	Create(pkg.ProductCreateUpdateReq) (pkg.IdResult[uuid.UUID], error)
	Update(uuid.UUID, pkg.ProductCreateUpdateReq) error
	Delete(uuid.UUID) error
}

type client struct {
	http        http.Client
	baseAddress string
}

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("unexpected status code: %d. details: %s", e.Code, e.Message)
}

func NewClient(baseAddress string) Client {
	return &client{baseAddress: baseAddress}
}

func sampleReq() pkg.ProductCreateUpdateReq {
	return pkg.ProductCreateUpdateReq{
		Name: "e2e_test_" + random.Str(4),
	}
}

func (c *client) GetList() (pkg.ProductsResponse, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseAddress+"/api/products", nil)
	if err != nil {
		return pkg.ProductsResponse{}, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return pkg.ProductsResponse{}, err
	}

	if response.StatusCode != http.StatusOK {
		var result pkg.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return pkg.ProductsResponse{}, err
		}

		return pkg.ProductsResponse{}, &APIError{
			Code:    response.StatusCode,
			Message: result.Message,
		}
	}

	var result pkg.ProductsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return pkg.ProductsResponse{}, err
	}

	return result, nil
}

func (c *client) GetSingle(id uuid.UUID) (pkg.ProductResponse, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseAddress+"/api/products/"+fmt.Sprint(id), nil)
	if err != nil {
		return pkg.ProductResponse{}, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return pkg.ProductResponse{}, err
	}

	if response.StatusCode != http.StatusOK {
		var result pkg.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return pkg.ProductResponse{}, err
		}

		return pkg.ProductResponse{}, &APIError{
			Code:    response.StatusCode,
			Message: result.Message,
		}
	}

	var result pkg.ProductResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return pkg.ProductResponse{}, err
	}

	return result, nil
}

func (c *client) Create(req pkg.ProductCreateUpdateReq) (pkg.IdResult[uuid.UUID], error) {

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return pkg.IdResult[uuid.UUID]{}, err
	}

	request, err := http.NewRequest(http.MethodPost, c.baseAddress+"/api/products", bytes.NewReader(jsonBytes))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return pkg.IdResult[uuid.UUID]{}, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return pkg.IdResult[uuid.UUID]{}, err
	}

	if response.StatusCode != http.StatusCreated {
		var result pkg.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return pkg.IdResult[uuid.UUID]{}, err
		}

		return pkg.IdResult[uuid.UUID]{}, &APIError{
			Code:    response.StatusCode,
			Message: result.Message,
		}
	}

	var result pkg.IdResult[uuid.UUID]
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return pkg.IdResult[uuid.UUID]{}, err
	}

	return result, nil
}

func (c *client) Update(id uuid.UUID, req pkg.ProductCreateUpdateReq) error {

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, c.baseAddress+"/api/products/"+fmt.Sprint(id), bytes.NewReader(jsonBytes))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		var result pkg.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return err
		}

		return &APIError{
			Code:    response.StatusCode,
			Message: result.Message,
		}
	}

	return nil
}

func (c *client) Delete(id uuid.UUID) error {
	request, err := http.NewRequest(http.MethodDelete, c.baseAddress+"/api/products/"+fmt.Sprint(id), nil)
	if err != nil {
		return err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		var result pkg.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return err
		}

		return &APIError{
			Code:    response.StatusCode,
			Message: result.Message,
		}
	}

	return nil
}
