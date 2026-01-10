package comment

import (
	"context"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5"
)

// s *CommentService openapi.CommentApiServicer
var _ openapi.CommentAPIServicer = (*CommentService)(nil)

type CommentService struct {
	DB                 dbpgx.Querier
	DBConn             dbpgx.DBTX
	SnowflakeGenerator snowflake.IDGenerator
}

func (s *CommentService) CreateComment(ctx context.Context, comment openapi.Comment) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateComment")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)
	// validate user exists
	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, comment.Username)
	if errGetUser != nil {
		switch errGetUser {
		case pgx.ErrNoRows:
			return openapi.Response(http.StatusNotFound, openapi.Error{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}), nil

		default:
			logger.Error("Error getting user", "error", errGetUser)

			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}), nil
		}
	}

	// Generate snowflake ID for the comment
	commentID, err := s.SnowflakeGenerator.NextID()
	if err != nil {
		logger.Error("Error generating snowflake ID for comment", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	params := db.CreateCommentForUserWithIDParams{
		ID:       commentID,
		Username: comment.Username,
		Content:  comment.Content,
	}

	createdComment, err := s.DB.CreateCommentForUserWithID(ctx, s.DBConn, params)
	if err != nil {
		logger.Error("Error creating comment", "error", err)
		return openapi.Response(http.StatusNotFound, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	c := converter.FromDBCmtToAPICmt(createdComment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int64) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.GetComment")
	defer span.End()
	comment, err := s.DB.GetComment(ctx, s.DBConn, id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return openapi.Response(http.StatusNotFound, openapi.Error{
				Code:    http.StatusNotFound,
				Message: "Comment not found",
			}), nil
		default:
			logger := contexthelper.GetLoggerInContext(ctx)
			logger.Error("Error getting comment", "error", err)
			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error getting comment",
			}), nil
		}
	}
	// get username
	user, errGetUser := s.DB.GetUserByID(ctx, s.DBConn, comment.UserID)
	if errGetUser != nil {
		logger := contexthelper.GetLoggerInContext(ctx)
		logger.Error("Error getting username for comment author", "error", errGetUser)
		return openapi.Response(http.StatusNotFound, openapi.Error{
			Code:    http.StatusNotFound,
			Message: "Error getting username for comment author",
		}), nil
	}

	c := converter.FromDBCmtToAPICmt(comment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetUserFeed(ctx context.Context) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.GetUserFeed")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)
	// validate the user exists
	// get username from context
	username, exists := contexthelper.GetUsernameInContext(ctx)
	if !exists {
		logger.Error("Error getting user from context")

		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if errGetUser != nil {
		switch errGetUser {
		case pgx.ErrNoRows:
			return openapi.Response(http.StatusNotFound, openapi.Error{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}), nil
		default:
			logger.Error("Error getting user", "error", errGetUser)
			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}), nil
		}
	}

	// get followed users
	followedUsers, err := s.DB.GetFollowedUsers(ctx, s.DBConn, user.ID)
	if err != nil {
		logger.Error("Error getting followed users", "error", err)
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	// get comments for each followed user
	comments := make([]openapi.Comment, 0, len(followedUsers)*20)
	for _, currentFollowedUser := range followedUsers {
		userComments, err := s.DB.GetUserComments(ctx, s.DBConn, db.GetUserCommentsParams{
			Username: currentFollowedUser.Username,
			Limit:    20,
			Offset:   0,
		})
		logger.Info("userComments", "comments", userComments)
		if err != nil {
			logger.Error("Error getting user comments", "error", err)
			return openapi.Response(http.StatusNotFound, nil), nil
		}
		for _, currentComment := range userComments {
			if currentComment.DeletedAt.Valid {
				// skip deleted comments
				continue
			}
			apiComment := converter.FromDBCmtToAPICmt(currentComment, currentFollowedUser)
			comments = append(comments, apiComment)
		}
	}
	for i := range comments {
		logger.Info("comments", "comment", comments[i])
	}

	return openapi.Response(http.StatusOK, comments), nil
}
