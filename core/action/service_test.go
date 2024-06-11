package action_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/action/mocks"
	"github.com/goto/shield/core/user"
	"github.com/stretchr/testify/mock"
)

var mockAction = action.Action{
	ID:          "action_id",
	Name:        "name",
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
		setup   func(ar *mocks.Repository)
		want    action.Action
		wantErr bool
	}{
		{
			name: "should call get repository if service being called",
			id:   mockAction.ID,
			setup: func(ar *mocks.Repository) {
				ar.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mockAction.ID).Return(mockAction, nil)
			},
			want: mockAction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := new(mocks.Repository)
			tt.setup(ar)
			s := action.NewService(log.NewNoop(), ar, nil, nil)
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
		act    action.Action
		setup  func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService)
		ErrStr string
	}{
		{
			name: "should return error if context has no user information",
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if repository return error",
			act:  mockAction,
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				ar.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), mockAction).Return(action.Action{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should not return error if succeed",
			act:  mockAction,
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				ar.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), mockAction).Return(mockAction, nil)
				as.EXPECT().Log(mock.AnythingOfType("context.withoutCancelCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("activity.Actor"), mock.AnythingOfType("action.LogData")).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := new(mocks.Repository)
			as := new(mocks.ActivityService)
			us := new(mocks.UserService)
			tt.setup(ar, us, as)
			s := action.NewService(log.NewNoop(), ar, us, as)
			_, err := s.Upsert(context.TODO(), tt.act)
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
		setup   func(ar *mocks.Repository)
		want    []action.Action
		wantErr bool
	}{
		{
			name: "should call repository if service being called",
			id:   mockAction.ID,
			setup: func(ar *mocks.Repository) {
				ar.EXPECT().List(mock.AnythingOfType("context.todoCtx")).Return([]action.Action{mockAction}, nil)
			},
			want: []action.Action{mockAction},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := new(mocks.Repository)
			tt.setup(ar)
			s := action.NewService(log.NewNoop(), ar, nil, nil)
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
		act    action.Action
		setup  func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService)
		ErrStr string
	}{
		{
			name: "should return error if context has no user information",
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should return error if repository return error",
			act:  mockAction,
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				ar.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), mockAction).Return(action.Action{}, errors.New("some error"))
			},
			ErrStr: "some error",
		},
		{
			name: "should not return error if succeed",
			act:  mockAction,
			setup: func(ar *mocks.Repository, us *mocks.UserService, as *mocks.ActivityService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(mockUser, nil)
				ar.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), mockAction).Return(mockAction, nil)
				as.EXPECT().Log(mock.AnythingOfType("context.withoutCancelCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("activity.Actor"), mock.AnythingOfType("action.LogData")).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := new(mocks.Repository)
			as := new(mocks.ActivityService)
			us := new(mocks.UserService)
			tt.setup(ar, us, as)
			s := action.NewService(log.NewNoop(), ar, us, as)
			_, err := s.Update(context.TODO(), tt.act.ID, tt.act)
			if err != nil {
				if err.Error() != tt.ErrStr {
					t.Fatalf("got error %s, expected was %s", err.Error(), tt.ErrStr)
				}
			}
		})
	}
}
