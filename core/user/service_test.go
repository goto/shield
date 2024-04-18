package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/goto/shield/core/user"

	"github.com/goto/shield/core/mocks"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "CreateUserWithUpperCase",
			user: user.User{
				Name:  "John Doe",
				Email: "John.Doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				repository.EXPECT().
					Create(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}).
					Return(user.User{
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}, nil).Once()
				return user.NewService(repository)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.Create(context.Background(), tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "UpdateUserWithUpperCase",
			user: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				repository.EXPECT().
					UpdateByID(mock.Anything, user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}).
					Return(user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}, nil).Once()
				return user.NewService(repository)
			},
			want: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "john.doe2@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.UpdateByID(context.Background(), tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "UpdateUserWithUpperCase",
			user: user.User{
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				repository.EXPECT().
					UpdateByEmail(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}).
					Return(user.User{
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}, nil).Once()
				return user.NewService(repository)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "john.doe2@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.UpdateByEmail(context.Background(), tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
