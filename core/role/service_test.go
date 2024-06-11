package role_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/role/mocks"
	"github.com/goto/shield/core/user"
	"github.com/stretchr/testify/mock"
)

var mockRole = role.Role{
	ID:          "id",
	Name:        "name",
	Types:       []string{"type1", "type2"},
	NamespaceID: "nsid",
}

var mockUser = user.User{
	Name:  "name",
	Email: "email",
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(rr *mocks.Repository)
		want    role.Role
		wantErr bool
	}{
		{
			name: "should call get repository if service being called",
			id:   mockRole.ID,
			setup: func(rr *mocks.Repository) {
				rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mockRole.ID).Return(mockRole, nil)
			},
			want: mockRole,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := new(mocks.Repository)
			tt.setup(rr)
			s := role.NewService(log.NewNoop(), rr, nil, nil)
			got, err := s.Get(context.TODO(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Upsert(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		ctx    context.Context
		rl     role.Role
		setup  func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService)
		ErrStr string
	}{
		{
			name: "should return error if context has no user information",
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if upsert repository return error",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), mockRole).Return("", errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if get repository return error",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), mockRole).Return("", errors.New("some error"))
				rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("string")).Return(role.Role{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should not return error if succeed",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), mockRole).Return(mockRole.ID, nil)
				rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("string")).Return(mockRole, nil)
				as.EXPECT().Log(mock.AnythingOfType("context.withoutCancelCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("activity.Actor"), mock.AnythingOfType("role.LogData")).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := new(mocks.Repository)
			as := new(mocks.ActivityService)
			us := new(mocks.UserService)
			tt.setup(rr, us, as)
			s := role.NewService(log.NewNoop(), rr, us, as)
			_, err := s.Upsert(context.TODO(), tt.rl)
			if err != nil {
				if err.Error() != tt.ErrStr {
					t.Fatalf("got error %s, expected was %s", err.Error(), tt.ErrStr)
				}
			}
		})
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(rr *mocks.Repository)
		want    []role.Role
		wantErr bool
	}{
		{
			name: "should call repository if service being called",
			id:   mockRole.ID,
			setup: func(rr *mocks.Repository) {
				rr.EXPECT().List(mock.AnythingOfType("context.todoCtx")).Return([]role.Role{mockRole}, nil)
			},
			want: []role.Role{mockRole},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := new(mocks.Repository)
			tt.setup(rr)
			s := role.NewService(log.NewNoop(), rr, nil, nil)
			got, err := s.List(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		ctx    context.Context
		rl     role.Role
		setup  func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService)
		ErrStr string
	}{
		{
			name: "should return error if context has no user information",
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if update repository return error",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), mockRole).Return("", errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if get repository return error",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), mockRole).Return("", errors.New("some error"))
				rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("string")).Return(role.Role{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should not return error if succeed",
			rl:   mockRole,
			setup: func(rr *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), mockRole).Return(mockRole.ID, nil)
				rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("string")).Return(mockRole, nil)
				as.EXPECT().Log(mock.AnythingOfType("context.withoutCancelCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("activity.Actor"), mock.AnythingOfType("role.LogData")).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := new(mocks.Repository)
			as := new(mocks.ActivityService)
			us := new(mocks.UserService)
			tt.setup(rr, us, as)
			s := role.NewService(log.NewNoop(), rr, us, as)
			_, err := s.Update(context.TODO(), tt.rl)
			if err != nil {
				if err.Error() != tt.ErrStr {
					t.Fatalf("got error %s, expected was %s", err.Error(), tt.ErrStr)
				}
			}
		})
	}
}
