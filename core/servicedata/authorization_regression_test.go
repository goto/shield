package servicedata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/servicedata/mocks"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	errorsPkg "github.com/goto/shield/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCrossUserWriteRequiresExistingServiceDataKeyEditPermission(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		allowed               bool
		identityErr, checkErr error
	}{
		{name: "missing-identity", identityErr: user.ErrMissingEmail},
		{name: "ordinary-user-without-key-grant"},
		{name: "permission-lookup-error", checkErr: errors.New("synthetic unavailable")},
		{name: "authorized-delegated-writer", allowed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			projects := mocks.NewProjectService(t)
			users := mocks.NewUserService(t)
			relations := mocks.NewRelationService(t)
			actor := user.User{ID: "actor", Email: "actor@example.test"}
			input := servicedata.ServiceData{NamespaceID: schema.UserPrincipal, EntityID: "different-target-user", Key: servicedata.Key{ProjectID: "project", Name: "metadata-key"}, Value: "synthetic"}
			users.On("FetchCurrentUser", mock.Anything).Return(actor, tc.identityErr).Once()
			if tc.identityErr == nil {
				projects.On("Get", mock.Anything, "project").Return(project.Project{ID: "project", Slug: "project-slug"}, nil).Once()
				key := servicedata.Key{ProjectID: "project", Name: "metadata-key", ResourceID: "key-resource"}
				repo.On("GetKeyByURN", mock.Anything, servicedata.CreateURN("project-slug", "metadata-key")).Return(key, nil).Once()
				relations.On("CheckPermission", mock.Anything, actor, namespace.Namespace{ID: schema.ServiceDataKeyNamespace}, "key-resource", action.Action{ID: schema.EditPermission}).Return(tc.allowed, tc.checkErr).Once()
				if tc.allowed {
					repo.On("Upsert", mock.Anything, mock.MatchedBy(func(sd servicedata.ServiceData) bool {
						return sd.EntityID == "different-target-user" && sd.Key.ResourceID == "key-resource"
					})).Return(input, nil).Once()
				}
			}
			svc := servicedata.NewService(testLogger, repo, nil, relations, projects, users, nil)
			_, err := svc.Upsert(context.Background(), input)
			if tc.allowed {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				repo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
			}
			if tc.identityErr == nil && tc.checkErr == nil && !tc.allowed {
				require.ErrorIs(t, err, errorsPkg.ErrForbidden)
			}
		})
	}
}
