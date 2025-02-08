package usecase

import (
	"context"
	"fmt"
	"hash/crc32"
	"testing"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/mocks/internal_/usecase"

	"github.com/stretchr/testify/suite"
)

type ShortUrlUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
	now time.Time

	repo *usecase.UseCase
	impl UseCase
}

func TestShortUrlUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlUseCaseTestSuite))
}

func (s *ShortUrlUseCaseTestSuite) SetupSuite() {
	s.now = time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC)
	now = func() time.Time {
		return s.now
	}

	s.repo = usecase.NewUseCase(s.T())
	s.impl = NewShortUrlUseCase(s.repo)
}

func (s *ShortUrlUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ShortUrlUseCaseTestSuite) TearDownSubTest() {}

func (s *ShortUrlUseCaseTestSuite) TearDownTest() {}

func (s *ShortUrlUseCaseTestSuite) TearDownSuite() {}

func (s *ShortUrlUseCaseTestSuite) TestCreate() {
	for _, t := range []struct {
		name   string
		req    *domain.CreateReqDto
		setup  func()
		expID  string
		expErr error
	}{
		{
			name: "create record successfully",
			req: &domain.CreateReqDto{
				Url:       "https://example.com/whatever1",
				ExpiredAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
			},
			setup: func() {
				targetID := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever1")))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:       "https://example.com/whatever1",
					TargetID:  targetID,
					ExpiredAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("testid1", nil)
			},
			expID:  "testid1",
			expErr: nil,
		},
		{
			name: "failed to create record due to duplicated target_id",
			req: &domain.CreateReqDto{
				Url:       "https://example.com/whatever2",
				TargetID:  "testid1",
				ExpiredAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
			},
			setup: func() {
				targetID := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever2")))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:       "https://example.com/whatever2",
					TargetID:  targetID,
					ExpiredAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("", domain.ErrDuplicatedKey)
			},
			expID:  "",
			expErr: domain.ErrDuplicatedKey,
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			id, err := s.impl.Create(ctx, t.req)
			s.ErrorIs(err, t.expErr)
			s.Equal(t.expID, id)
		})
	}
}

func (s *ShortUrlUseCaseTestSuite) TestGet() {
	for _, t := range []struct {
		name   string
		req    string
		setup  func()
		expObj *domain.GetRespDto
		expErr error
	}{
		{
			name: "create record successfully",
			req:  "testid1",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid1").Once().Return(&domain.GetRespDto{
					Url:       "https://example.com/whatever1",
					ExpiredAt: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expObj: &domain.GetRespDto{
				Url:       "https://example.com/whatever1",
				ExpiredAt: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			expErr: nil,
		},
		{
			name: "failed to create record due to duplicated target_id",
			req:  "testid2",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid2").Once().Return(nil, domain.ErrExpired)
			},
			expObj: nil,
			expErr: domain.ErrExpired,
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			obj, err := s.impl.Get(ctx, t.req)
			s.ErrorIs(err, t.expErr)
			s.Equal(t.expObj, obj)
		})
	}
}
