package store

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nguyen-duc-loc/task-management/backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomTask(t *testing.T) Task {
	id, err := gonanoid.New()
	require.NoError(t, err)
	require.NotZero(t, len(id))

	arg := CreateTaskParams{
		ID:        id,
		Name:      util.RandomPrintableString(100),
		CreatorID: createRandomUser(t).ID,
		Deadline:  time.Now().Add(time.Hour),
	}

	task, err := testStore.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, arg.ID, task.ID)
	require.Equal(t, arg.Name, task.Name)
	require.Equal(t, arg.CreatorID, task.CreatorID)
	require.WithinDuration(t, arg.Deadline, task.Deadline, time.Second)

	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t)
}

func TestGetTasks(t *testing.T) {
	task1 := createRandomTask(t)
	arg := GetTasksParams{
		CreatorID: task1.CreatorID,
		Name: pgtype.Text{
			String: task1.Name[30:70],
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
		Limit:  5,
		Offset: 0,
	}

	tasks, err := testStore.GetTasks(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, len(tasks), 1)

	task2 := tasks[0]
	require.NotEmpty(t, task2)
	require.Equal(t, task1, task2)
}

func TestGetTasksWithNullValue(t *testing.T) {
	task1 := createRandomTask(t)
	arg := GetTasksParams{
		CreatorID: task1.CreatorID,
		Limit:     5,
		Offset:    0,
	}

	tasks, err := testStore.GetTasks(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, len(tasks), 1)

	task2 := tasks[0]
	require.NotEmpty(t, task2)
	require.Equal(t, task1, task2)
}
