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

func TestGetTaskByID(t *testing.T) {
	task1 := createRandomTask(t)
	task2, err := testStore.GetTaskByID(context.Background(), task1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1, task2)
}

func TestUpdateTaskOnlyName(t *testing.T) {
	oldTask := createRandomTask(t)
	newName := util.RandomPrintableString(100)

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Name: pgtype.Text{
			String: newName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Name, updatedTask.Name)
	require.Equal(t, newName, updatedTask.Name)
	require.Equal(t, oldTask.Deadline, updatedTask.Deadline)
	require.Equal(t, oldTask.Completed, updatedTask.Completed)
}

func TestUpdateTaskOnlyDeadline(t *testing.T) {
	oldTask := createRandomTask(t)
	newDeadline := time.Now().Add(2 * time.Hour)

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Deadline: pgtype.Timestamptz{
			Time:  newDeadline,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Deadline, updatedTask.Deadline)
	require.WithinDuration(t, newDeadline, updatedTask.Deadline, time.Second)
	require.Equal(t, oldTask.Name, updatedTask.Name)
	require.Equal(t, oldTask.Completed, updatedTask.Completed)
}

func TestUpdateTaskOnlyStatus(t *testing.T) {
	oldTask := createRandomTask(t)
	newStatus := !oldTask.Completed

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Completed: pgtype.Bool{
			Bool:  newStatus,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Completed, updatedTask.Completed)
	require.Equal(t, newStatus, updatedTask.Completed)
	require.Equal(t, oldTask.Name, updatedTask.Name)
	require.Equal(t, oldTask.Deadline, updatedTask.Deadline)
}

func TestUpdateTaskAllFields(t *testing.T) {
	oldTask := createRandomTask(t)
	newName := util.RandomPrintableString(100)
	newDeadline := time.Now().Add(2 * time.Hour)
	newStatus := !oldTask.Completed

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Name: pgtype.Text{
			String: newName,
			Valid:  true,
		},
		Deadline: pgtype.Timestamptz{
			Time:  newDeadline,
			Valid: true,
		},
		Completed: pgtype.Bool{
			Bool:  newStatus,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Name, updatedTask.Name)
	require.NotEqual(t, oldTask.Deadline, updatedTask.Deadline)
	require.NotEqual(t, oldTask.Completed, updatedTask.Completed)
	require.Equal(t, newName, updatedTask.Name)
	require.WithinDuration(t, newDeadline, updatedTask.Deadline, time.Second)
	require.Equal(t, newStatus, updatedTask.Completed)
}

func TestDeleteTask(t *testing.T) {
	task1 := createRandomTask(t)
	err := testStore.DeleteTask(context.Background(), task1.ID)
	require.NoError(t, err)
	task2, err := testStore.GetTaskByID(context.Background(), task1.ID)
	require.Error(t, err)
	require.Empty(t, task2)
}
