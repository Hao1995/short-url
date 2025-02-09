package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/internal/router/handler/request"
	"github.com/Hao1995/short-url/mocks/internal_/usecase"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ShortUrlHandlerTestSuite struct {
	suite.Suite
	ginEngine *gin.Engine

	now time.Time

	uc   *usecase.UseCase
	impl *ShortUrlHandler
}

func TestShortUrlHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlHandlerTestSuite))
}

func (s *ShortUrlHandlerTestSuite) SetupSuite() {
	s.now = time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC)

	s.uc = usecase.NewUseCase(s.T())
	s.impl = NewShortUrlHandler(s.uc)

	r := gin.Default()
	r.POST("/api/v1/urls", s.impl.Create)
	r.GET("/:id", s.impl.Get)
	s.ginEngine = r
}

func (s *ShortUrlHandlerTestSuite) SetupTest() {}

func (s *ShortUrlHandlerTestSuite) TearDownSubTest() {}

func (s *ShortUrlHandlerTestSuite) TearDownTest() {}

func (s *ShortUrlHandlerTestSuite) TearDownSuite() {}

func (s *ShortUrlHandlerTestSuite) TestCreate() {
	for _, t := range []struct {
		name    string
		req     *request.ShortUrlCreateRequest
		setup   func()
		expCode int
		expResp string
		expErr  error
	}{
		{
			name: "create record successfully",
			req: &request.ShortUrlCreateRequest{
				Url:      "https://example.com/whatever1",
				ExpireAt: s.now,
			},
			setup: func() {
				s.uc.On("Create", mock.Anything, &domain.CreateReqDto{
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now,
				}).Once().Return(&domain.CreateRespDto{
					TargetID: "testid1",
					ShortUrl: "http://localhost/testid1",
				}, nil)
			},
			expCode: 201,
			expResp: "{\"id\":\"testid1\",\"shortUrl\":\"http://localhost/testid1\"}",
		},
		{
			name:    "failed to bind request data",
			req:     &request.ShortUrlCreateRequest{},
			expCode: 422,
			expResp: fmt.Sprintf("{\"error\":\"%s\"}", "unprocessable entity"),
		},
		{
			name: "failed to create a short url",
			req: &request.ShortUrlCreateRequest{
				Url:      "https://example.com/whatever1",
				ExpireAt: s.now,
			},
			setup: func() {
				s.uc.On("Create", mock.Anything, &domain.CreateReqDto{
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now,
				}).Once().Return(nil, errors.New("whatever"))
			},
			expCode: 500,
			expResp: fmt.Sprintf("{\"error\":\"%s\"}", "internal server error"),
		},
	} {
		s.Suite.Run(t.name, func() {
			if t.setup != nil {
				t.setup()
			}

			w := httptest.NewRecorder()
			data, _ := json.Marshal(t.req)
			req, _ := http.NewRequest("POST", "/api/v1/urls", strings.NewReader(string(data)))
			s.ginEngine.ServeHTTP(w, req)

			s.Equal(t.expCode, w.Code)
			s.Equal(t.expResp, w.Body.String())
		})
	}
}

func (s *ShortUrlHandlerTestSuite) TestGet() {
	for _, t := range []struct {
		name        string
		req         *request.ShortUrlGetRequest
		setup       func()
		expCode     int
		expResp     string
		expLocation string
		expErr      error
	}{
		{
			name: "create record successfully",
			req:  &request.ShortUrlGetRequest{ID: "whatever1"},
			setup: func() {
				s.uc.On("Get", mock.Anything, "whatever1").
					Once().
					Return(&domain.GetRespDto{
						Status:   domain.GetRespStatusNormal,
						Url:      "https://example.com/whatever1",
						ExpireAt: s.now,
					}, nil)
			},
			expCode:     302,
			expResp:     "<a href=\"https://example.com/whatever1\">Found</a>.\n\n",
			expLocation: "https://example.com/whatever1",
		},
		{
			name: "record not found, return 404",
			req:  &request.ShortUrlGetRequest{ID: "whatever1"},
			setup: func() {
				s.uc.On("Get", mock.Anything, "whatever1").
					Once().
					Return(&domain.GetRespDto{
						Status:   domain.GetRespStatusNotFound,
						Url:      "",
						ExpireAt: time.Time{},
					}, nil)
			},
			expCode:     404,
			expResp:     fmt.Sprintf("{\"error\":\"%s\"}", "not found"),
			expLocation: "",
		},
		{
			name: "record is expired, return 404",
			req:  &request.ShortUrlGetRequest{ID: "whatever1"},
			setup: func() {
				s.uc.On("Get", mock.Anything, "whatever1").
					Once().
					Return(&domain.GetRespDto{
						Status:   domain.GetRespStatusExpired,
						Url:      "https://example.com/whatever1",
						ExpireAt: s.now,
					}, nil)
			},
			expCode:     404,
			expResp:     fmt.Sprintf("{\"error\":\"%s\"}", "not found"),
			expLocation: "",
		},
	} {
		s.Suite.Run(t.name, func() {
			if t.setup != nil {
				t.setup()
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/"+t.req.ID, nil)
			s.ginEngine.ServeHTTP(w, req)

			s.Equal(t.expCode, w.Code)
			s.Equal(t.expResp, w.Body.String())
			s.Equal(t.expLocation, w.Header().Get("location"))
		})
	}
}
