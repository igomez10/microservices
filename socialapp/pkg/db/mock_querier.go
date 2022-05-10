package db

import "context"

const CREATECOMMENT = "CreateComment"
const CREATEUSER = "CreateUser"
const DELETEUSER = "DeleteUser"
const GETCOMMENT = "GetComment"
const GETUSERBYEMAIL = "GetUserByEmail"
const GETUSERBYID = "GetUserByID"
const LISTCOMMENT = "ListComment"
const LISTUSERS = "ListUsers"

type MockQuerier struct {
	CommentToReturn  Comment
	UserToReturn     User
	CommentsToReturn []Comment
	UsersToReturn    []User
	CallRegister     map[string][]interface{}
}

func NewMockQuerier() *MockQuerier {
	m := &MockQuerier{}
	m.CallRegister = map[string][]interface{}{}
	return m
}

func (m *MockQuerier) WithCommentToReturn(c Comment) {
	m.CommentToReturn = c
}
func (m *MockQuerier) WithCommentsToReturn(cs []Comment) {
	m.CommentsToReturn = cs
}
func (m *MockQuerier) WithUserToReturn(u User) {
	m.UserToReturn = u
}
func (m *MockQuerier) WithUsersToReturn(us []User) {
	m.UsersToReturn = us
}

func (m *MockQuerier) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	callName := CREATECOMMENT
	m.CallRegister[callName] = append(m.CallRegister[callName], arg)
	return m.CommentToReturn, nil
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	callName := CREATEUSER
	m.CallRegister[callName] = append(m.CallRegister[callName], arg)
	return m.UserToReturn, nil
}

// We do not allow deletes
func (m *MockQuerier) DeleteUser(ctx context.Context, id int32) error {
	callName := DELETEUSER
	m.CallRegister[callName] = append(m.CallRegister[callName], id)
	return nil
}

// Comment
func (m *MockQuerier) GetComment(ctx context.Context, id int32) (Comment, error) {
	callName := GETCOMMENT
	m.CallRegister[callName] = append(m.CallRegister[callName], id)
	return m.CommentToReturn, nil
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (User, error) {
	callName := GETUSERBYEMAIL
	m.CallRegister[callName] = append(m.CallRegister[callName], email)
	return m.UserToReturn, nil
}

//User
func (m *MockQuerier) GetUserByID(ctx context.Context, id int32) (User, error) {
	callName := GETUSERBYID
	m.CallRegister[callName] = []interface{}{"hello"}
	return m.UserToReturn, nil
}

func (m *MockQuerier) ListComment(ctx context.Context) ([]Comment, error) {
	callName := LISTCOMMENT
	m.CallRegister[callName] = append(m.CallRegister[callName], true)
	return m.CommentsToReturn, nil
}

func (m *MockQuerier) ListUsers(ctx context.Context) ([]User, error) {
	callName := LISTUSERS
	m.CallRegister[callName] = append(m.CallRegister[callName], true)
	return m.UsersToReturn, nil
}
