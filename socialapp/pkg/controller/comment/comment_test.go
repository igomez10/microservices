package comment

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type stubQuerier struct {
	db.Querier
	getUserByUsernameFunc func(ctx context.Context, dbtx db.DBTX, username string) (db.User, error)
	getCommentFunc        func(ctx context.Context, dbtx db.DBTX, id int64) (db.GetCommentRow, error)
	getUserByIDFunc       func(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error)
	createLikeFunc        func(ctx context.Context, dbtx db.DBTX, arg db.CreateLikeParams) error
	deleteLikeFunc        func(ctx context.Context, dbtx db.DBTX, arg db.DeleteLikeParams) error
}

func (s *stubQuerier) GetUserByUsername(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
	return s.getUserByUsernameFunc(ctx, dbtx, username)
}

func (s *stubQuerier) GetComment(ctx context.Context, dbtx db.DBTX, id int64) (db.GetCommentRow, error) {
	return s.getCommentFunc(ctx, dbtx, id)
}

func (s *stubQuerier) GetUserByID(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
	return s.getUserByIDFunc(ctx, dbtx, id)
}

func (s *stubQuerier) CreateLike(ctx context.Context, dbtx db.DBTX, arg db.CreateLikeParams) error {
	return s.createLikeFunc(ctx, dbtx, arg)
}

func (s *stubQuerier) DeleteLike(ctx context.Context, dbtx db.DBTX, arg db.DeleteLikeParams) error {
	return s.deleteLikeFunc(ctx, dbtx, arg)
}

func newCommentTestContext(username string) context.Context {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = contexthelper.SetUsernameInContext(req, username)
	return contexthelper.SetLoggerInContext(req.Context(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
}

func TestLikeComment(t *testing.T) {
	var createLikeArg db.CreateLikeParams
	getCommentCalls := 0

	service := &CommentService{
		DB: &stubQuerier{
			getUserByUsernameFunc: func(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
				return db.User{ID: 7, Username: username}, nil
			},
			getCommentFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.GetCommentRow, error) {
				getCommentCalls++
				likeCount := int64(0)
				if getCommentCalls > 1 {
					likeCount = 1
				}
				return db.GetCommentRow{
					ID:        id,
					Content:   "hello world",
					LikeCount: likeCount,
					UserID:    11,
					CreatedAt: pgtype.Timestamp{Time: time.Unix(1700000000, 0).UTC(), Valid: true},
				}, nil
			},
			getUserByIDFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
				return db.User{ID: id, Username: "author"}, nil
			},
			createLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.CreateLikeParams) error {
				createLikeArg = arg
				return nil
			},
			deleteLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.DeleteLikeParams) error {
				return nil
			},
		},
	}

	res, err := service.LikeComment(newCommentTestContext("some2"), openapi.LikeRequest{CommentId: "42"})
	if err != nil {
		t.Fatalf("LikeComment returned error: %v", err)
	}
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	body, ok := res.Body.(openapi.Comment)
	if !ok {
		t.Fatalf("expected comment body, got %T", res.Body)
	}
	if body.LikeCount != "1" {
		t.Fatalf("expected like_count \"1\", got %q", body.LikeCount)
	}
	if createLikeArg.UserID != 7 || createLikeArg.CommentID != 42 {
		t.Fatalf("unexpected create like params: %+v", createLikeArg)
	}
}

func TestUnlikeComment(t *testing.T) {
	var deleteLikeArg db.DeleteLikeParams
	getCommentCalls := 0

	service := &CommentService{
		DB: &stubQuerier{
			getUserByUsernameFunc: func(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
				return db.User{ID: 7, Username: username}, nil
			},
			getCommentFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.GetCommentRow, error) {
				getCommentCalls++
				likeCount := int64(1)
				if getCommentCalls > 1 {
					likeCount = 0
				}
				return db.GetCommentRow{
					ID:        id,
					Content:   "hello world",
					LikeCount: likeCount,
					UserID:    11,
					CreatedAt: pgtype.Timestamp{Time: time.Unix(1700000000, 0).UTC(), Valid: true},
				}, nil
			},
			getUserByIDFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
				return db.User{ID: id, Username: "author"}, nil
			},
			createLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.CreateLikeParams) error {
				return nil
			},
			deleteLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.DeleteLikeParams) error {
				deleteLikeArg = arg
				return nil
			},
		},
	}

	res, err := service.UnlikeComment(newCommentTestContext("some2"), openapi.LikeRequest{CommentId: "42"})
	if err != nil {
		t.Fatalf("UnlikeComment returned error: %v", err)
	}
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	body, ok := res.Body.(openapi.Comment)
	if !ok {
		t.Fatalf("expected comment body, got %T", res.Body)
	}
	if body.LikeCount != "0" {
		t.Fatalf("expected like_count \"0\", got %q", body.LikeCount)
	}
	if deleteLikeArg.UserID != 7 || deleteLikeArg.CommentID != 42 {
		t.Fatalf("unexpected delete like params: %+v", deleteLikeArg)
	}
}

func TestLikeCommentRequiresAuthenticatedUser(t *testing.T) {
	service := &CommentService{}

	res, err := service.LikeComment(context.Background(), openapi.LikeRequest{CommentId: "42"})
	if err != nil {
		t.Fatalf("LikeComment returned error: %v", err)
	}
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}

func TestLikeCommentReturnsNotFoundForMissingComment(t *testing.T) {
	service := &CommentService{
		DB: &stubQuerier{
			getUserByUsernameFunc: func(ctx context.Context, dbtx db.DBTX, username string) (db.User, error) {
				return db.User{ID: 7, Username: username}, nil
			},
			getCommentFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.GetCommentRow, error) {
				return db.GetCommentRow{}, pgx.ErrNoRows
			},
			getUserByIDFunc: func(ctx context.Context, dbtx db.DBTX, id int64) (db.User, error) {
				return db.User{}, nil
			},
			createLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.CreateLikeParams) error {
				return nil
			},
			deleteLikeFunc: func(ctx context.Context, dbtx db.DBTX, arg db.DeleteLikeParams) error {
				return nil
			},
		},
	}

	res, err := service.LikeComment(newCommentTestContext("some2"), openapi.LikeRequest{CommentId: "404"})
	if err != nil {
		t.Fatalf("LikeComment returned error: %v", err)
	}
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.Code)
	}
}
