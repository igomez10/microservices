package comment

import (
	"context"
	"database/sql"
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
				Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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

	newCommentResult, err := s.DB.CreateCommentForUser(ctx, s.DBConn, params)
	if err != nil {
		log.Error().Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).Err(err).Msg("Error creating comment")
		return openapi.Response(http.StatusNotFound, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}), nil
	}

	newCommentID, err := newCommentResult.LastInsertId()
	if err != nil {
		log.Error().Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).Err(err).Msg("Error getting last insert id")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	// get comment from db
	newComment, errGetComment := s.DB.GetComment(ctx, s.DBConn, newCommentID)
	if errGetComment != nil {
		log.Error().Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).Err(errGetComment).Msg("Error getting comment from db")
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	c := FromDBCmtToAPICmt(newComment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int32) (openapi.ImplResponse, error) {
	comment, err := s.DB.GetComment(ctx, s.DBConn, int64(id))
	if err != nil {
		log.Error().Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}
	// get username
	user, errGetUser := s.DB.GetUserByID(ctx, s.DBConn, comment.UserID)
	if errGetUser != nil {
		log.Error().Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).Err(errGetUser).Msg("Error getting usrname for comment author")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	c := FromDBCmtToAPICmt(comment, user)
	return openapi.Response(http.StatusOK, c), nil
}

func (s *CommentService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	// validate the user exists
	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if errGetUser != nil {
		log.Error().
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Err(errGetUser).
			Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	comments, err := s.DB.GetUserComments(ctx, s.DBConn, username)
	if err != nil {
		log.Error().
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Err(err).
			Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	apiComments := make([]openapi.Comment, 0, len(comments))
	for i := range comments {
		if comments[i].DeletedAt.Valid {
			continue
		}

		c := FromDBCmtToAPICmt(comments[i], user)
		apiComments = append(apiComments, c)
	}

	return openapi.Response(http.StatusOK, comments), nil
}

func (s *CommentService) GetUserFeed(ctx context.Context, username string) (openapi.ImplResponse, error) {
	// validate the user exists
	user, errGetUser := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if errGetUser != nil {
		log.Error().
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Err(errGetUser).
			Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	// get followed users
	followedUsers, err := s.DB.GetFollowedUsers(ctx, s.DBConn, user.ID)
	if err != nil {
		log.Error().
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Err(err).
			Msg("Error getting followed users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	log.Info().Msgf("followed users: %v", followedUsers)

	// get comments for each followed user
	comments := make([]openapi.Comment, 0, len(followedUsers)*20)
	for _, currentFollowedUser := range followedUsers {
		userComments, err := s.DB.GetUserComments(ctx, s.DBConn, currentFollowedUser.Username)
		log.Info().Msgf("userComments: \n%v\n", userComments)
		if err != nil {
			log.Error().
				Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
				Err(err).
				Msg("Error getting user comments")
			return openapi.Response(http.StatusNotFound, nil), nil
		}
		for _, currentComment := range userComments {
			if currentComment.DeletedAt.Valid {
				// skip deleted comments
				continue
			}
			apiComment := FromDBCmtToAPICmt(currentComment, currentFollowedUser)
			comments = append(comments, apiComment)
		}
	}
	for i := range comments {
		log.Info().Msgf("comments: \n%v\n", comments[i])
	}

	return openapi.Response(http.StatusOK, comments), nil
}

func FromDBCmtToAPICmt(comment db.Comment, user db.User) openapi.Comment {
	return openapi.Comment{
		Id:        comment.ID,
		Content:   comment.Content,
		LikeCount: int64(comment.LikeCount),
		CreatedAt: comment.CreatedAt,
		Username:  user.Username,
	}
}
