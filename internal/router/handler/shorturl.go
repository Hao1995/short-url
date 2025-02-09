package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/internal/router/handler/request"
	"github.com/Hao1995/short-url/internal/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("not found")
)

type ShortUrlHandler struct {
	uc usecase.UseCase
}

func NewShortUrlHandler(uc usecase.UseCase) *ShortUrlHandler {
	return &ShortUrlHandler{
		uc: uc,
	}
}

// Create creates short_url record and return short url id
func (hlr *ShortUrlHandler) Create(c *gin.Context) {
	var req request.ShortUrlCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("handler.Create. failed to bind json: %s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ErrUnprocessableEntity.Error()})
		return
	}

	obj, err := hlr.uc.Create(c.Request.Context(), &domain.CreateReqDto{Url: req.Url, ExpiredAt: req.ExpiredAt})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServerError.Error()})
		return
	}

	log.Printf("handler.Create. success create a short utl: %s", obj.ShortUrl)
	c.JSON(http.StatusCreated, gin.H{
		"id":       obj.TargetID,
		"shortUrl": obj.ShortUrl,
	})
}

// Get redirects to the original url
func (hlr *ShortUrlHandler) Get(c *gin.Context) {
	var req request.ShortUrlGetRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.Printf("handler.Get. failed to bind uri: %s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ErrUnprocessableEntity.Error()})
		return
	}

	obj, err := hlr.uc.Get(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServerError.Error()})
		return
	}

	if obj.Status == domain.GetRespStatusNotFound || obj.Status == domain.GetRespStatusExpired {
		log.Printf("handler.Get. get abnormal status: %s, return 404", obj.Status)
		c.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound.Error()})
		return
	}

	log.Printf("handler.Get. success redirect to: %s", obj.Url)
	c.Redirect(http.StatusFound, obj.Url)
}
