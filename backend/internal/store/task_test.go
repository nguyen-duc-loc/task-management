package store

import (
	"context"
	"testing"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nguyen-duc-loc/task-management/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
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
}
