package comment

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/eventRecorder"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/events"
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// parseAPIID parses an int64 ID that arrived on the API surface as a string
// (see openapi.yaml: int64 fields are serialized as JSON strings to avoid
// JS precision loss). Returns a 400 ApiError payload on parse failure.
func parseAPIID(s string) (int64, *openapi.Error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, &openapi.Error{
			Code:    http.StatusBadRequest,
			Message: "invalid id: must be a numeric string",
		}
	}
	return id, nil
}

// s *CommentService openapi.CommentApiServicer
var _ openapi.CommentAPIServicer = (*CommentService)(nil)

type CommentService struct {
	DB                 dbpgx.Querier
	DBConn             *pgxpool.Pool
	EventRecorder      eventRecorder.EventRecorder
	SnowflakeGenerator snowflake.IDGenerator
}

func (s *CommentService) loadAPIComment(ctx context.Context, id int64) (openapi.Comment, *openapi.Error, error) {
	comment, err := s.DB.GetComment(ctx, s.DBConn, id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return openapi.Comment{}, &openapi.Error{
				Code:    http.StatusNotFound,
				Message: "Comment not found",
			}, nil
		default:
			logger := contexthelper.GetLoggerInContext(ctx)
			logger.Error("Error getting comment", "error", err)
			return openapi.Comment{}, &openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error getting comment",
			}, nil
		}
	}

	user, err := s.DB.GetUserByID(ctx, s.DBConn, comment.UserID)
	if err != nil {
		logger := contexthelper.GetLoggerInContext(ctx)
		logger.Error("Error getting username for comment author", "error", err)
		return openapi.Comment{}, &openapi.Error{
			Code:    http.StatusNotFound,
			Message: "Error getting username for comment author",
		}, nil
	}

	return openapi.Comment{
		Id:        strconv.FormatInt(comment.ID, 10),
		Content:   comment.Content,
		LikeCount: strconv.FormatInt(comment.LikeCount, 10),
		CreatedAt: comment.CreatedAt.Time,
		Username:  user.Username,
	}, nil, nil
}

func (s *CommentService) SearchComments(ctx context.Context, username string, startTime time.Time, endTime time.Time) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.SearchComments")
	defer span.End()

	logger := contexthelper.GetLoggerInContext(ctx)

	if startTime.IsZero() {
		startTime = time.Unix(0, 0).UTC()
	}
	if endTime.IsZero() {
		endTime = time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)
	}
	if startTime.After(endTime) {
		return openapi.Response(http.StatusBadRequest, openapi.Error{
			Code:    http.StatusBadRequest,
			Message: "start_time must be before or equal to end_time",
		}), nil
	}

	dbComments, err := s.DB.SearchComments(ctx, s.DBConn, db.SearchCommentsParams{
		Column1: username,
		CreatedAt: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
	})
	if err != nil {
		logger.Error("Error searching comments", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error searching comments",
		}), nil
	}

	apiComments := make([]openapi.Comment, len(dbComments))
	for i := range dbComments {
		apiComments[i] = openapi.Comment{
			Id:        strconv.FormatInt(dbComments[i].ID, 10),
			Content:   dbComments[i].Content,
			LikeCount: strconv.FormatInt(dbComments[i].LikeCount, 10),
			CreatedAt: dbComments[i].CreatedAt.Time,
			Username:  dbComments[i].Username,
		}
	}

	return openapi.Response(http.StatusOK, apiComments), nil
}

func (s *CommentService) CreateComment(ctx context.Context, comment openapi.CreateCommentRequest) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateComment")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	authenticatedUsername, exists := contexthelper.GetUsernameInContext(ctx)
	if !exists {
		return openapi.Response(http.StatusUnauthorized, openapi.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}), nil
	}
	if comment.Username != authenticatedUsername {
		return openapi.Response(http.StatusForbidden, openapi.Error{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}), nil
	}

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

	tx, err := s.DBConn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}
	defer tx.Rollback(ctx)

	createdComment, err := s.DB.CreateCommentForUserWithID(ctx, tx, params)
	if err != nil {
		logger.Error("Error creating comment", "error", err)
		return openapi.Response(http.StatusNotFound, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	if err := s.EventRecorder.RecordEvent(ctx, tx, eventRecorder.Event{
		AggregateID:   createdComment.ID,
		AggregateType: events.AggregateTypeComment,
		EventType:     events.EventTypeCommentCreated,
		Payload: events.CommentCreated{
			CommentID:      createdComment.ID,
			AuthorID:       user.ID,
			AuthorUsername: user.Username,
			Content:        createdComment.Content,
		},
	}); err != nil {
		logger.Error("Error recording event", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	c := converter.FromDBCmtToAPICmt(createdComment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetComment(ctx context.Context, id string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.GetComment")
	defer span.End()

	commentID, parseErr := parseAPIID(id)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}

	comment, apiErr, err := s.loadAPIComment(ctx, commentID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiErr != nil {
		return openapi.Response(int(apiErr.Code), *apiErr), nil
	}

	return openapi.Response(http.StatusOK, comment), nil
}

func (s *CommentService) LikeComment(ctx context.Context, req openapi.LikeRequest) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.LikeComment")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	username, exists := contexthelper.GetUsernameInContext(ctx)
	if !exists {
		return openapi.Response(http.StatusUnauthorized, openapi.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}), nil
	}

	user, err := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return openapi.Response(http.StatusUnauthorized, openapi.Error{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			}), nil
		default:
			logger.Error("Error getting user", "error", err)
			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}), nil
		}
	}

	commentID, parseErr := parseAPIID(req.CommentId)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}

	if _, apiErr, err := s.loadAPIComment(ctx, commentID); err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	} else if apiErr != nil {
		return openapi.Response(int(apiErr.Code), *apiErr), nil
	}

	if err := s.DB.CreateLike(ctx, s.DBConn, db.CreateLikeParams{
		UserID:    user.ID,
		CommentID: commentID,
	}); err != nil {
		logger.Error("Error creating like", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error creating like",
		}), nil
	}

	comment, apiErr, err := s.loadAPIComment(ctx, commentID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiErr != nil {
		return openapi.Response(int(apiErr.Code), *apiErr), nil
	}

	return openapi.Response(http.StatusOK, comment), nil
}

func (s *CommentService) UnlikeComment(ctx context.Context, req openapi.LikeRequest) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.UnlikeComment")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	username, exists := contexthelper.GetUsernameInContext(ctx)
	if !exists {
		return openapi.Response(http.StatusUnauthorized, openapi.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}), nil
	}

	user, err := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return openapi.Response(http.StatusUnauthorized, openapi.Error{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			}), nil
		default:
			logger.Error("Error getting user", "error", err)
			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}), nil
		}
	}

	commentID, parseErr := parseAPIID(req.CommentId)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}

	if _, apiErr, err := s.loadAPIComment(ctx, commentID); err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	} else if apiErr != nil {
		return openapi.Response(int(apiErr.Code), *apiErr), nil
	}

	if err := s.DB.DeleteLike(ctx, s.DBConn, db.DeleteLikeParams{
		UserID:    user.ID,
		CommentID: commentID,
	}); err != nil {
		logger.Error("Error deleting like", "error", err)
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error deleting like",
		}), nil
	}

	comment, apiErr, err := s.loadAPIComment(ctx, commentID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiErr != nil {
		return openapi.Response(int(apiErr.Code), *apiErr), nil
	}

	return openapi.Response(http.StatusOK, comment), nil
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
			apiComment := openapi.Comment{
				Id:        strconv.FormatInt(currentComment.ID, 10),
				Content:   currentComment.Content,
				LikeCount: strconv.FormatInt(currentComment.LikeCount, 10),
				CreatedAt: currentComment.CreatedAt.Time,
				Username:  currentFollowedUser.Username,
			}
			comments = append(comments, apiComment)
		}
	}
	for i := range comments {
		logger.Info("comments", "comment", comments[i])
	}

	return openapi.Response(http.StatusOK, comments), nil
}
