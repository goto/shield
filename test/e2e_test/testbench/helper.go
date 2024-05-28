package testbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/db"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	OrgAdminEmail      = "admin1-group1-org1@gotocompany.com"
	DefaultSystemEmail = "shield-service@gotocompany.com"
	IdentityHeader     = "X-Shield-Email"
	userIDHeaderKey    = "X-Shield-UserID"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func createConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	return grpc.DialContext(ctx, host, opts...)
}

func CreateClient(ctx context.Context, host string) (shieldv1beta1.ShieldServiceClient, func(), error) {
	conn, err := createConnection(ctx, host)
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		conn.Close()
	}

	client := shieldv1beta1.NewShieldServiceClient(conn)
	return client, cancel, nil
}

func CreatePublicClient(ctx context.Context, host string) (shieldv1beta1.ShieldPublicServiceClient, func(), error) {
	conn, err := createConnection(ctx, host)
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		conn.Close()
	}

	client := shieldv1beta1.NewShieldPublicServiceClient(conn)
	return client, cancel, nil
}

func BootstrapUser(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-user.json")
	if err != nil {
		return err
	}

	var data []*shieldv1beta1.UserRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	for _, d := range data {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			IdentityHeader: creatorEmail,
		}))
		if _, err := cl.CreateUser(ctx, &shieldv1beta1.CreateUserRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func BootstrapMetadataKey(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-metadata-key.json")
	if err != nil {
		return err
	}

	var data []*shieldv1beta1.MetadataKeyRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	for _, d := range data {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			IdentityHeader: creatorEmail,
		}))
		if _, err := cl.CreateMetadataKey(ctx, &shieldv1beta1.CreateMetadataKeyRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func BootstrapOrganization(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-organization.json")
	if err != nil {
		return err
	}

	var data []*shieldv1beta1.OrganizationRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	for _, d := range data {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			IdentityHeader: creatorEmail,
		}))
		if _, err := cl.CreateOrganization(ctx, &shieldv1beta1.CreateOrganizationRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func BootstrapProject(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-project.json")
	if err != nil {
		return err
	}

	orgResp, err := cl.ListOrganizations(ctx, &shieldv1beta1.ListOrganizationsRequest{})
	if err != nil {
		return err
	}

	if len(orgResp.GetOrganizations()) < 1 {
		return errors.New("no organization found")
	}

	var data []*shieldv1beta1.ProjectRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	data[0].OrgId = orgResp.GetOrganizations()[0].GetId()

	for _, d := range data {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			IdentityHeader: creatorEmail,
		}))
		if _, err := cl.CreateProject(ctx, &shieldv1beta1.CreateProjectRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func BootstrapGroup(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-group.json")
	if err != nil {
		return err
	}

	orgResp, err := cl.ListOrganizations(ctx, &shieldv1beta1.ListOrganizationsRequest{})
	if err != nil {
		return err
	}

	if len(orgResp.GetOrganizations()) < 1 {
		return errors.New("no organization found")
	}

	var data []*shieldv1beta1.GroupRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	data[0].OrgId = orgResp.GetOrganizations()[0].GetId()
	data[1].OrgId = orgResp.GetOrganizations()[0].GetId()
	data[2].OrgId = orgResp.GetOrganizations()[0].GetId()

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		IdentityHeader: creatorEmail,
	}))

	for _, d := range data {
		if _, err := cl.CreateGroup(ctx, &shieldv1beta1.CreateGroupRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func AssignGroupManager(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string) error {
	groupsResp, err := cl.ListGroups(ctx, &shieldv1beta1.ListGroupsRequest{})
	if err != nil {
		return err
	}

	if len(groupsResp.GetGroups()) < 1 {
		return errors.New("no groups found")
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		IdentityHeader: creatorEmail,
	}))

	usr, err := cl.GetCurrentUser(ctx, &shieldv1beta1.GetCurrentUserRequest{})
	if err != nil {
		return err
	}

	for _, grp := range groupsResp.GetGroups() {
		// assign admin to group
		_, err = cl.CreateRelation(ctx, &shieldv1beta1.CreateRelationRequest{
			Body: &shieldv1beta1.RelationRequestBody{
				ObjectId:        grp.GetId(),
				ObjectNamespace: schema.GroupNamespace,
				Subject:         fmt.Sprintf("%s:%s", schema.UserPrincipal, usr.GetUser().GetId()),
				RoleName:        schema.ManagerRole,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func BootstrapResource(ctx context.Context, cl shieldv1beta1.ShieldServiceClient, creatorEmail string, testDataPath string) error {
	testFixtureJSON, err := os.ReadFile(testDataPath + "/mocks/mock-resource.json")
	if err != nil {
		return err
	}

	var data []*shieldv1beta1.ResourceRequestBody
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	projResp, err := cl.ListProjects(ctx, &shieldv1beta1.ListProjectsRequest{})
	if err != nil {
		return err
	}

	if len(projResp.GetProjects()) < 1 {
		return errors.New("no projects found")
	}

	data[0].ProjectId = projResp.GetProjects()[0].GetId()
	data[1].ProjectId = projResp.GetProjects()[0].GetId()

	for _, d := range data {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			IdentityHeader: creatorEmail,
		}))
		if _, err := cl.CreateResource(ctx, &shieldv1beta1.CreateResourceRequest{
			Body: d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func SetupDB(cfg db.Config) (dbc *db.Client, err error) {
	dbc, err = db.New(cfg)
	if err != nil {
		err = fmt.Errorf("failed to setup db: %w", err)
		return
	}

	return
}
