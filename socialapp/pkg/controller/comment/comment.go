package comment

import (
	"context"
	"net/http"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	"github.com/rs/zerolog/log"
)

// s *CommentService openapi.CommentApiServicer
type CommentService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *CommentService) CreateComment(ctx context.Context, comment openapi.Comment) (openapi.ImplResponse, error) {
	params := db.CreateCommentForUserParams{
		Username: comment.Username,
		Content:  comment.Content,
	}

	newComment, err := s.DB.CreateCommentForUser(ctx, s.DBConn, params)
	if err != nil {
		log.Error().Err(err).Msg("Error creating comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, newComment), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int32) (openapi.ImplResponse, error) {
	commnet, err := s.DB.GetComment(ctx, s.DBConn, id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

func (s *CommentService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	commnet, err := s.DB.GetUserComments(ctx, s.DBConn, username)
	if err != nil {
		log.Error().Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}
