package controller

import (
	model "go-api/models"
	usecase "go-api/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(usecase usecase.ProductUseCase) productController {
	return productController{
		productUseCase: usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	products, err := p.productUseCase.GetProducts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, products)
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var product model.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedProduct, err := p.productUseCase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *productController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUseCase.GetProductById(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: "Produto não encontrado no banco de dados.",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	existingProduct, err := p.productUseCase.GetProductById(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if existingProduct == nil {
		response := model.Response{
			Message: "Produto não encontrado no banco de dados.",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	var updateData model.Product

	err = ctx.BindJSON(&updateData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if updateData.Name != "" {
		existingProduct.Name = updateData.Name
	}

	if updateData.Price != 0 {
		existingProduct.Price = updateData.Price
	}

	updatedProduct, err := p.productUseCase.UpdateProduct(*existingProduct)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	deleted, err := p.productUseCase.DeleteProduct(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if !deleted {
		response := model.Response{
			Message: "Produto não encontrado no banco de dados.",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := model.Response{
		Message: "Produto deletado com sucesso.",
	}
	ctx.JSON(http.StatusOK, response)
}
