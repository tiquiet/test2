package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "sber/internal/repository/mocks"
)

func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockJSONStorage(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc := NewDataService(logger, mockRepo)

	ctx := context.TODO()
	key := "validKey"
	expectedData := "some data"

	// Sunny cases

	mockRepo.EXPECT().Get(ctx, key).Return(expectedData, nil)

	result, err := svc.GetData(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)

	body := []byte(`{"foo": "bar"}`)
	expTime := "3600"

	mockRepo.EXPECT().Put(ctx, key, string(body), 3600).Return(key, nil)

	resultKey, err := svc.SaveData(ctx, key, body, expTime)
	assert.NoError(t, err)
	assert.Equal(t, key, resultKey)

	// Rainy cases

	mockRepo.EXPECT().Get(ctx, key).Return("", errors.New("some error"))

	result, err = svc.GetData(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, "", result)

	key = "validKey"
	body = []byte(`invalid json`)
	expTime = "3600"

	resultKey, err = svc.SaveData(ctx, key, body, expTime)
	assert.Error(t, err)
	assert.Equal(t, "", resultKey)

	key = "validKey"
	body = []byte(`{"foo": "bar"}`)
	expTime = "invalid"

	mockRepo.EXPECT().Put(ctx, key, string(body), 0).Return(key, nil)

	resultKey, err = svc.SaveData(ctx, key, body, expTime)
	assert.NoError(t, err)
	assert.Equal(t, key, resultKey)
}
