package authz

import (
	"context"

	"github.com/raystack/shield/model"

	"github.com/raystack/salt/log"
	"github.com/raystack/shield/config"
	"github.com/raystack/shield/internal/authz/spicedb"
)

type Policy interface {
	AddPolicy(ctx context.Context, schema string) error
}

type Permission interface {
	AddRelation(ctx context.Context, relation model.Relation) error
	DeleteRelation(ctx context.Context, relation model.Relation) error
	CheckRelation(ctx context.Context, relation model.Relation, action model.Action) (bool, error)
	DeleteSubjectRelations(ctx context.Context, resource model.Resource) error
}

type Authz struct {
	Policy
	Permission
}

func New(config *config.Shield, logger log.Logger) *Authz {
	spice, err := spicedb.New(config.SpiceDB, logger)

	if err != nil {
		logger.Fatal(err.Error())
	}

	return &Authz{
		spice.Policy,
		spice.Permission,
	}
}
