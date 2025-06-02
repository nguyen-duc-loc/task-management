package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/nguyen-duc-loc/task-management/backend/internal/database/mock"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomUser(t *testing.T) (user store.User, password string) {
	password = util.RandomPrintableString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = store.User{
		ID:             rand.Int64(),
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
	}

	return
}

type eqCreateUserParamsMatcher struct {
	arg      store.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg store.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg,
		password,
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user store.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Data    store.User `json:"data"`
		Success bool       `json:"success"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	gotUser := response.Data
	require.Equal(t, user.Username, gotUser.Username)
	require.Empty(t, gotUser.HashedPassword)
}

func TestCreateUserHandler(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.CreateUserParams{
					Username: user.Username,
				}
				storage.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.User{}, store.ErrUniqueViolation)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username": user.Username,
				"password": "short",
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mockdb.NewMockStorage(ctrl)
			tc.buildStubs(storage)

			server, err := NewServer(storage)
			require.NoError(t, err)
			server.RegisterRoutes()
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserHandler(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.CreateUserParams{
					Username: user.Username,
				}
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(arg.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.User{}, store.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		}, {
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username": user.Username,
				"password": "short",
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mockdb.NewMockStorage(ctrl)
			tc.buildStubs(storage)

			server, err := NewServer(storage)
			require.NoError(t, err)
			server.RegisterRoutes()
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
