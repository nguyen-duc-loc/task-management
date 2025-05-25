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
		ID:    id,
		Title: util.RandomPrintableString(50),
		Description: pgtype.Text{
			String: util.RandomPrintableString(300),
			Valid:  true,
		},
		CreatorID: createRandomUser(t).ID,
		Deadline:  time.Now().Add(time.Hour),
	}

	task, err := testStore.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, arg.ID, task.ID)
	require.Equal(t, arg.Title, task.Title)
	require.Equal(t, arg.Description, task.Description)
	require.Equal(t, arg.CreatorID, task.CreatorID)
	require.WithinDuration(t, arg.Deadline, task.Deadline, time.Second)

	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t)
}

func TestGetTasks(t *testing.T) {
	task := createRandomTask(t)
	arg := GetTasksParams{
		CreatorID: task.CreatorID,
		Title: pgtype.Text{
			String: task.Title[12:41],
			Valid:  true,
		},
		Description: pgtype.Text{
			String: task.Description.String[30:70],
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
	require.Positive(t, len(tasks))

	for _, task := range tasks {
		require.NotEmpty(t, task)
	}
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
	require.Positive(t, len(tasks))

	for _, task := range tasks {
		require.NotEmpty(t, task)
	}
}

func TestGetTaskByID(t *testing.T) {
	task1 := createRandomTask(t)
	task2, err := testStore.GetTaskByID(context.Background(), task1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1, task2)
}

func TestUpdateTaskOnlyTitle(t *testing.T) {
	oldTask := createRandomTask(t)
	newTitle := util.RandomPrintableString(50)

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Title: pgtype.Text{
			String: newTitle,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Title, updatedTask.Title)
	require.Equal(t, newTitle, updatedTask.Title)
	require.Equal(t, oldTask.Description, updatedTask.Description)
	require.Equal(t, oldTask.Deadline, updatedTask.Deadline)
	require.Equal(t, oldTask.Completed, updatedTask.Completed)
}

func TestUpdateTaskOnlyDescription(t *testing.T) {
	oldTask := createRandomTask(t)
	newDescription := util.RandomPrintableString(300)

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Description: pgtype.Text{
			String: newDescription,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldTask.Description, updatedTask.Description)
	require.Equal(t, newDescription, updatedTask.Description.String)
	require.Equal(t, oldTask.Title, updatedTask.Title)
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
	require.Equal(t, oldTask.Title, updatedTask.Title)
	require.Equal(t, oldTask.Description, updatedTask.Description)
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
	require.Equal(t, oldTask.Title, updatedTask.Title)
	require.Equal(t, oldTask.Description, updatedTask.Description)
	require.Equal(t, oldTask.Deadline, updatedTask.Deadline)
}

func TestUpdateTaskAllFields(t *testing.T) {
	oldTask := createRandomTask(t)
	newTitle := util.RandomPrintableString(50)
	newDescription := util.RandomPrintableString(300)
	newDeadline := time.Now().Add(2 * time.Hour)
	newStatus := !oldTask.Completed

	updatedTask, err := testStore.UpdateTask(context.Background(), UpdateTaskParams{
		ID: oldTask.ID,
		Title: pgtype.Text{
			String: newTitle,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: newDescription,
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
	require.NotEqual(t, oldTask.Title, updatedTask.Title)
	require.NotEqual(t, oldTask.Description, updatedTask.Description)
	require.NotEqual(t, oldTask.Deadline, updatedTask.Deadline)
	require.NotEqual(t, oldTask.Completed, updatedTask.Completed)
	require.Equal(t, newTitle, updatedTask.Title)
	require.Equal(t, newDescription, updatedTask.Description.String)
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
