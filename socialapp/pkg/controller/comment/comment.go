package comment

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/db"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/rs/zerolog/log"
)

// s *CommentService openapi.CommentApiServicer
type CommentService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *CommentService) CreateComment(ctx context.Context, comment openapi.Comment) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateComment")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	// validate user exists
	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, comment.Username)
	if errGetUser != nil {
		switch errGetUser {
		case sql.ErrNoRows:
			return openapi.Response(http.StatusNotFound, openapi.Error{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}), nil

		default:
			log.Error().
				Err(errGetUser).
				Msg("Error getting user")

			return openapi.Response(http.StatusInternalServerError, openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}), nil
		}
	}

	params := db.CreateCommentForUserParams{
		Username: comment.Username,
		Content:  comment.Content,
	}

	createdComment, err := s.DB.CreateCommentForUser(ctx, s.DBConn, params)
	if err != nil {
		log.Error().Err(err).Msg("Error creating comment")
		return openapi.Response(http.StatusNotFound, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	c := converter.FromDBCmtToAPICmt(createdComment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.GetComment")
	defer span.End()
	comment, err := s.DB.GetComment(ctx, s.DBConn, int64(id))
	if err != nil {
		log.Error().Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}
	// get username
	user, errGetUser := s.DB.GetUserByID(ctx, s.DBConn, comment.UserID)
	if errGetUser != nil {
		log.Error().Err(errGetUser).Msg("Error getting usrname for comment author")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	c := converter.FromDBCmtToAPICmt(comment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetUserFeed(ctx context.Context) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CommentService.GetUserFeed")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	// validate the user exists
	// get username from context
	username, exists := contexthelper.GetUsernameInContext(ctx)
	if !exists {
		log.Error().
			Msg("Error getting user from context")

		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if errGetUser != nil {
		log.Error().
			Err(errGetUser).
			Msg("Error getting user")
		return openapi.Response(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	// get followed users
	followedUsers, err := s.DB.GetFollowedUsers(ctx, s.DBConn, user.ID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error getting followed users")
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
		log.Info().Msgf("userComments: \n%v\n", userComments)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error getting user comments")
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
		log.Info().Msgf("comments: \n%v\n", comments[i])
	}

	return openapi.Response(http.StatusOK, comments), nil
}
