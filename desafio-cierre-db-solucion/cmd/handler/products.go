package handler

import (
	"net/http"

	"github.com/bootcamp-go/desafio-cierre-db.git/internal/domain"
	"github.com/bootcamp-go/desafio-cierre-db.git/internal/products"
	"github.com/gin-gonic/gin"
)

type Products struct {
	s products.Service
}

func NewHandlerProducts(s products.Service) *Products {
	return &Products{s}
}

func (p *Products) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.s.ReadAll()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, products)
	}
}

func (p *Products) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products := domain.Product{}
		err := ctx.ShouldBindJSON(&products)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		err = p.s.Create(&products)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"data": products})
	}
}

func (p *Products) PostAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products := []domain.Product{}

		err := ctx.ShouldBindJSON(&products)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		for _, product := range products {
			err = p.s.Create(&product)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(201, gin.H{"data": products})
	}
}

func (p *Products) TopProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		response, err := p.s.TopProducts()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}
