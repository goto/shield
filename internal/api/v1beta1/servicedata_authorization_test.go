package v1beta1

import (
	"context"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	errorsPkg "github.com/goto/shield/pkg/errors"
	pb "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestCrossUserServiceDataDenialReachesAPI(t *testing.T) {
	users := new(mocks.UserService)
	sd := new(mocks.ServiceDataService)
	defer users.AssertExpectations(t)
	defer sd.AssertExpectations(t)
	users.On("Get", mock.Anything, "target-user").Return(user.User{ID: "target-user"}, nil).Once()
	sd.On("Upsert", mock.Anything, mock.MatchedBy(func(value servicedata.ServiceData) bool {
		return value.EntityID == "target-user" && value.Key.Name == "metadata-key"
	})).Return(servicedata.ServiceData{}, errorsPkg.ErrForbidden).Once()
	h := Handler{userService: users, serviceDataService: sd, serviceDataConfig: ServiceDataConfig{MaxUpsert: 1}}
	response, err := h.UpsertUserServiceData(context.Background(), &pb.UpsertUserServiceDataRequest{UserId: "target-user", Body: &pb.UpsertServiceDataRequestBody{Project: "project", Data: &structpb.Struct{Fields: map[string]*structpb.Value{"metadata-key": structpb.NewStringValue("synthetic")}}}})
	require.Nil(t, response)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}
