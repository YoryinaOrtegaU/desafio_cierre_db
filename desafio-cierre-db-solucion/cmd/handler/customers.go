package handler

import (
	"github.com/bootcamp-go/desafio-cierre-db.git/internal/customers"
	"github.com/bootcamp-go/desafio-cierre-db.git/internal/domain"
	"github.com/gin-gonic/gin"
)

type Customers struct {
	s customers.Service
}

func NewHandlerCustomers(s customers.Service) *Customers {
	return &Customers{s}
}

func (c *Customers) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customers, err := c.s.ReadAll()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, customers)
	}
}

func (c *Customers) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customer := domain.Customers{}
		err := ctx.ShouldBindJSON(&customer)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		err = c.s.Create(&customer)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"data": customer})
	}
}

func (c *Customers) PostAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customers := []domain.Customers{}

		err := ctx.ShouldBindJSON(&customers)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		for _, customer := range customers {
			err = c.s.Create(&customer)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(201, gin.H{"data": customers})
	}
}
