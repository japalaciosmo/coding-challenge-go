package product

import (
	"coding-challenge-go/pkg/api/seller"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	LIST_PAGE_SIZE = 10
)

func NewController(repository *repository, sellerRepository *seller.Repository, sellerEmailProvider *seller.EmailProvider) *controller {
	return &controller{
		repository:          repository,
		sellerRepository:    sellerRepository,
		sellerEmailProvider: sellerEmailProvider,
	}
}

type controller struct {
	repository          *repository
	sellerRepository    *seller.Repository
	sellerEmailProvider *seller.EmailProvider
}

func (pc *controller) List(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.repository.list((request.Page-1)*LIST_PAGE_SIZE, LIST_PAGE_SIZE)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	productsJson, err := json.Marshal(products)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

func (pc *controller) ListV2(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.repository.list((request.Page-1)*LIST_PAGE_SIZE, LIST_PAGE_SIZE)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	for i, _ := range products {
		var sllr productSeller
		sllr.UUID = products[i].SellerUUID
		sllr.Link.Self.Href = fmt.Sprintf("http://localhost:8080/sellers/%s", products[i].SellerUUID)
		products[i].SellerUUID = ""
		products[i].Seller = &sllr
	}
	productsJson, err := json.Marshal(products)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

func (pc *controller) Get(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	productJson, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productJson)
}

func (pc *controller) GetV2(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}
	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}
	var sllr productSeller
	sllr.UUID = product.SellerUUID
	sllr.Link.Self.Href = fmt.Sprintf("http://localhost:8080/sellers/%s", product.SellerUUID)
	product.SellerUUID = ""
	product.Seller = &sllr

	productJson, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productJson)
}

func (pc *controller) Post(c *gin.Context) {
	request := &struct {
		Name   string `form:"name"`
		Brand  string `form:"brand"`
		Stock  int    `form:"stock"`
		Seller string `form:"seller"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := pc.sellerRepository.FindByUUID(request.Seller)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller by UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
		return
	}

	if seller == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller is not found"})
		return
	}

	product := &product{
		UUID:       uuid.New().String(),
		Name:       request.Name,
		Brand:      request.Brand,
		Stock:      request.Stock,
		SellerUUID: seller.UUID,
	}

	err = pc.repository.insert(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	jsonData, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *controller) Put(c *gin.Context) {
	queryRequest := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(queryRequest.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	request := &struct {
		Name  string `form:"name"`
		Brand string `form:"brand"`
		Stock int    `form:"stock"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	oldStock := product.Stock

	product.Name = request.Name
	product.Brand = request.Brand
	product.Stock = request.Stock

	err = pc.repository.update(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	if oldStock != product.Stock {
		seller, err := pc.sellerRepository.FindByUUID(product.SellerUUID)

		if err != nil {
			log.Error().Err(err).Msg("Fail to query seller by UUID")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
			return
		}

		pc.sellerEmailProvider.StockChanged(oldStock, product.Stock, *seller, product.Name)
	}

	jsonData, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *controller) Delete(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	err = pc.repository.delete(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (pc *controller) ListTopTen(c *gin.Context) {
	request := &struct {
		Page int `form:"top,default=10"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellers, err := pc.sellerRepository.ListTop(request.Page)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}
