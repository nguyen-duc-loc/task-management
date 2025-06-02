package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	mockdb "github.com/nguyen-duc-loc/task-management/backend/internal/database/mock"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/internal/token"
	"github.com/nguyen-duc-loc/task-management/backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomTask(t *testing.T, creatorID int64) store.Task {
	id, err := gonanoid.New()
	require.NoError(t, err)

	return store.Task{
		ID:        id,
		CreatorID: creatorID,
		Title:     util.RandomPrintableString(50),
		Description: pgtype.Text{
			String: util.RandomPrintableString(300),
			Valid:  true,
		},
		Deadline: time.Now().Add(time.Hour),
	}
}

type eqCreateTaskParamsMatcher struct {
	arg store.CreateTaskParams
}

func (e eqCreateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.CreateTaskParams)
	if !ok {
		return false
	}

	if e.arg.Title != arg.Title {
		return false
	}

	if e.arg.Description != arg.Description {
		return false
	}

	if e.arg.CreatorID != arg.CreatorID {
		return false
	}

	if e.arg.Deadline.Sub(arg.Deadline) > time.Second || arg.Deadline.Sub(e.arg.Deadline) > time.Second {
		return false
	}

	if len(arg.ID) == 0 {
		return false
	}

	return true
}

func (e eqCreateTaskParamsMatcher) String() string {
	return fmt.Sprintf("arg: %v", e.arg)
}

func EqCreateTaskParams(arg store.CreateTaskParams) gomock.Matcher {
	return eqCreateTaskParamsMatcher{arg}
}

func requireBodyMatchTask(t *testing.T, body *bytes.Buffer, task store.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask store.Task
	err = json.Unmarshal(data, &gotTask)

	require.NoError(t, err)
	require.Equal(t, gotTask.ID, task.ID)
	require.Equal(t, gotTask.CreatorID, task.CreatorID)
	require.Equal(t, gotTask.Title, task.Title)
	require.Equal(t, gotTask.Description, task.Description)
	require.WithinDuration(t, gotTask.Deadline, task.Deadline, time.Second)
}

func TestCreateTaskHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":       task.Title,
				"description": task.Description.String,
				"deadline":    task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.CreateTaskParams{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       task.Title,
					Description: task.Description,
					Deadline:    task.Deadline,
				}
				storage.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "No description",
			body: gin.H{
				"title":    task.Title,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.CreateTaskParams{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Title:     task.Title,
					Deadline:  task.Deadline,
				}
				storage.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"title":    task.Title,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title":    task.Title,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDeadline",
			body: gin.H{
				"title":    task.Title,
				"deadline": "invalidDeadline",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
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

			url := "/tasks"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

type eqGetTasksParamsMatcher struct {
	arg store.GetTasksParams
}

func (e eqGetTasksParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.GetTasksParams)
	if !ok {
		return false
	}

	if arg.CreatorID != e.arg.CreatorID {
		return false
	}

	if arg.Title.String != e.arg.Title.String {
		return false
	}

	if arg.Description.String != e.arg.Description.String {
		return false
	}

	if arg.Limit != e.arg.Limit {
		return false
	}

	if arg.Offset != e.arg.Offset {
		return false
	}

	if arg.StartDeadline.Time.Sub(e.arg.StartDeadline.Time).Abs() > time.Second {
		return false
	}

	if arg.EndDeadline.Time.Sub(e.arg.EndDeadline.Time).Abs() > time.Second {
		return false
	}

	if arg.Completed.Bool != e.arg.Completed.Bool {
		return false
	}

	return true
}

func (e eqGetTasksParamsMatcher) String() string {
	return fmt.Sprintf("arg: %v", e.arg)
}

func EqGetTasksParams(arg store.GetTasksParams) gomock.Matcher {
	return eqGetTasksParamsMatcher{arg}
}

func requireBodyMatchTasks(t *testing.T, body *bytes.Buffer, tasks []store.GetTasksRow) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rsp getTasksResponse
	err = json.Unmarshal(data, &rsp)
	require.NoError(t, err)
	gotTasks := rsp.Data
	require.Equal(t, len(tasks), int(rsp.Total))
	require.Equal(t, len(tasks), len(gotTasks))

	for i := range tasks {
		require.Equal(t, tasks[i].ID, gotTasks[i].ID)
		require.Equal(t, tasks[i].Title, gotTasks[i].Title)
		require.Equal(t, tasks[i].Description, gotTasks[i].Description)
		require.Equal(t, tasks[i].CreatorID, gotTasks[i].CreatorID)
		require.Equal(t, tasks[i].Completed, gotTasks[i].Completed)
		require.WithinDuration(t, tasks[i].Deadline, gotTasks[i].Deadline, time.Second)
		require.WithinDuration(t, tasks[i].CreatedAt, gotTasks[i].CreatedAt, time.Second)
	}
}

func TestGetTasksHandler(t *testing.T) {
	user, _ := randomUser(t)

	var n int64 = 10
	tasks := make([]store.GetTasksRow, n)
	for i := range n {
		rt := randomTask(t, user.ID)
		tasks[i].ID = rt.ID
		tasks[i].Title = rt.Title
		tasks[i].Description = rt.Description
		tasks[i].CreatorID = rt.CreatorID
		tasks[i].Deadline = rt.Deadline
		tasks[i].Completed = rt.Completed
		tasks[i].CreatedAt = rt.CreatedAt
		tasks[i].Total = n
	}

	type Query struct {
		Title         string
		Description   string
		StartDeadline string
		EndDeadline   string
		Completed     *bool
		Page          *int32
		Limit         *int32
	}

	incomplete, limit, page := new(bool), new(int32), new(int32)
	*incomplete = false
	*limit = int32(n)
	*page = 1

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Page:          page,
				Limit:         limit,
				Title:         "",
				Description:   "",
				StartDeadline: time.Now().Add(-time.Hour).Format(time.RFC3339),
				EndDeadline:   time.Now().Add(time.Hour).Format(time.RFC3339),
				Completed:     incomplete,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.GetTasksParams{
					CreatorID: user.ID,
					Limit:     int32(n),
					Offset:    0,
					Title: pgtype.Text{
						String: "",
						Valid:  true,
					},
					Description: pgtype.Text{
						String: "",
						Valid:  true,
					},
					StartDeadline: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
					EndDeadline: pgtype.Timestamptz{
						Time:  time.Now().Add(+time.Hour),
						Valid: true,
					},
					Completed: pgtype.Bool{
						Bool:  false,
						Valid: true,
					},
				}
				storage.EXPECT().
					GetTasks(gomock.Any(), EqGetTasksParams(arg)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name:  "NoQueryString",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.GetTasksParams{
					CreatorID: user.ID,
					Limit:     5,
					Offset:    0,
				}
				storage.EXPECT().
					GetTasks(gomock.Any(), EqGetTasksParams(arg)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name:  "UnauthorizedUser",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "InternalError",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]store.GetTasksRow{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDeadline",
			query: Query{
				StartDeadline: "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
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

			url := "/tasks"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()

			if len(tc.query.Title) > 0 {
				q.Add("title", tc.query.Title)
			}

			if len(tc.query.Description) > 0 {
				q.Add("description", tc.query.Description)
			}

			if len(tc.query.StartDeadline) > 0 {
				q.Add("start_deadline", tc.query.StartDeadline)
			}

			if len(tc.query.EndDeadline) > 0 {
				q.Add("end_deadline", tc.query.EndDeadline)
			}

			if tc.query.Completed != nil {
				q.Add("completed", strconv.FormatBool(*tc.query.Completed))
			}

			if tc.query.Page != nil {
				q.Add("page", strconv.FormatInt(int64(*tc.query.Page), 10))
			}

			if tc.query.Limit != nil {
				q.Add("limit", strconv.FormatInt(int64(*tc.query.Limit), 10))
			}

			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTaskByIDHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "NotFound",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, store.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID-1, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/tasks/%s", task.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

type eqUpdateTaskParamsMatcher struct {
	arg store.UpdateTaskParams
}

func (e eqUpdateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.UpdateTaskParams)
	if !ok {
		return false
	}

	if e.arg.Title != arg.Title {
		return false
	}

	if e.arg.Description != arg.Description {
		return false
	}

	if e.arg.Deadline.Time.Sub(arg.Deadline.Time) > time.Second || arg.Deadline.Time.Sub(e.arg.Deadline.Time) > time.Second {
		return false
	}

	if len(arg.ID) == 0 {
		return false
	}

	return true
}

func (e eqUpdateTaskParamsMatcher) String() string {
	return fmt.Sprintf("arg: %v", e.arg)
}

func EqUpdateTaskParams(arg store.UpdateTaskParams) gomock.Matcher {
	return eqUpdateTaskParamsMatcher{arg}
}

func TestUpdateHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)
	newTitle := util.RandomPrintableString(50)
	newDescription := util.RandomPrintableString(300)
	newDeadline := time.Now().Add(2 * time.Hour)
	newCompleted := !task.Completed

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "UpdateAllFields",
			body: gin.H{
				"title":       newTitle,
				"description": newDescription,
				"deadline":    newDeadline.Format(time.RFC3339),
				"completed":   newCompleted,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)

				newTask := store.Task{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Title:     newTitle,
					Description: pgtype.Text{
						String: newDescription,
						Valid:  true,
					},
					Deadline:  newDeadline,
					Completed: newCompleted,
				}

				arg := store.UpdateTaskParams{
					ID: task.ID,
					Title: pgtype.Text{
						String: newTask.Title,
						Valid:  true,
					},
					Description: newTask.Description,
					Deadline: pgtype.Timestamptz{
						Time:  newTask.Deadline,
						Valid: true,
					},
					Completed: pgtype.Bool{
						Bool:  newTask.Completed,
						Valid: true,
					},
				}

				storage.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(newTask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, store.Task{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Title:     newTitle,
					Description: pgtype.Text{
						String: newDescription,
						Valid:  true,
					},
					Deadline:  newDeadline,
					Completed: newCompleted,
					CreatedAt: task.CreatedAt,
				})
			},
		},
		{
			name: "UpdateOnlyDeadline",
			body: gin.H{
				"deadline": newDeadline.Format(time.RFC3339),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)

				newTask := store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       task.Title,
					Description: task.Description,
					Deadline:    newDeadline,
					Completed:   task.Completed,
				}

				arg := store.UpdateTaskParams{
					ID: task.ID,
					Deadline: pgtype.Timestamptz{
						Time:  newTask.Deadline,
						Valid: true,
					},
				}

				storage.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(newTask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       task.Title,
					Description: task.Description,
					Deadline:    newDeadline,
					Completed:   task.Completed,
					CreatedAt:   task.CreatedAt,
				})
			},
		},
		{
			name: "UpdateOnlyTitle",
			body: gin.H{
				"title": newTitle,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)

				newTask := store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       newTitle,
					Description: task.Description,
					Deadline:    task.Deadline,
					Completed:   task.Completed,
				}

				arg := store.UpdateTaskParams{
					ID: task.ID,
					Title: pgtype.Text{
						String: newTask.Title,
						Valid:  true,
					},
				}

				storage.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(newTask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       newTitle,
					Description: task.Description,
					Deadline:    task.Deadline,
					Completed:   task.Completed,
					CreatedAt:   task.CreatedAt,
				})
			},
		},
		{
			name: "UpdateOnlyDescription",
			body: gin.H{
				"description": newDescription,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)

				newTask := store.Task{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Title:     task.Title,
					Description: pgtype.Text{
						String: newDescription,
						Valid:  true,
					},
					Deadline:  task.Deadline,
					Completed: task.Completed,
				}

				arg := store.UpdateTaskParams{
					ID:          task.ID,
					Description: newTask.Description,
				}

				storage.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(newTask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, store.Task{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Title:     task.Title,
					Description: pgtype.Text{
						String: newDescription,
						Valid:  true,
					},
					Deadline:  task.Deadline,
					Completed: task.Completed,
					CreatedAt: task.CreatedAt,
				})
			},
		},
		{
			name: "UpdateOnlyCompleted",
			body: gin.H{
				"completed": newCompleted,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)

				newTask := store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       task.Title,
					Description: task.Description,
					Deadline:    task.Deadline,
					Completed:   newCompleted,
				}

				arg := store.UpdateTaskParams{
					ID: task.ID,
					Completed: pgtype.Bool{
						Bool:  newTask.Completed,
						Valid: true,
					},
				}

				storage.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(newTask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, store.Task{
					ID:          task.ID,
					CreatorID:   task.CreatorID,
					Title:       task.Title,
					Description: task.Description,
					Deadline:    task.Deadline,
					Completed:   newCompleted,
					CreatedAt:   task.CreatedAt,
				})
			},
		},
		{
			name: "NotFound",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, store.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID-1, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/tasks/%s", task.ID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
				storage.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, store.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID-1, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/tasks/%s", task.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
