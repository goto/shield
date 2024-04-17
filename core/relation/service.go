package relation

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
	"go.uber.org/zap"
)

const (
	auditKeyRelationCreate        = "relation.create"
	auditKeyRelationSubjectDelete = "relation_subject.delete"
)

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data any) error
}

type Service struct {
	logger          *zap.SugaredLogger
	repository      Repository
	authzRepository AuthzRepository
	userService     UserService
	activityService ActivityService
}

func NewService(logger *zap.SugaredLogger, repository Repository, authzRepository AuthzRepository, userService UserService, activityService ActivityService) *Service {
	return &Service{
		logger:          logger,
		repository:      repository,
		authzRepository: authzRepository,
		userService:     userService,
		activityService: activityService,
	}
}

func (s Service) Get(ctx context.Context, id string) (RelationV2, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, rel RelationV2) (RelationV2, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return RelationV2{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	createdRelation, err := s.repository.Create(ctx, rel)
	if err != nil {
		return RelationV2{}, fmt.Errorf("%w: %s", ErrCreatingRelationInStore, err.Error())
	}

	err = s.authzRepository.AddV2(ctx, createdRelation)
	if err != nil {
		return RelationV2{}, fmt.Errorf("%w: %s", ErrCreatingRelationInAuthzEngine, err.Error())
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		relationLogData := createdRelation.ToRelationLogData()
		if err := s.activityService.Log(ctx, auditKeyRelationCreate, currentUser.ID, relationLogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return createdRelation, nil
}

func (s Service) List(ctx context.Context) ([]RelationV2, error) {
	return s.repository.List(ctx)
}

// TODO: Update & Delete planned for v0.6
// TODO: Audit log
func (s Service) Update(ctx context.Context, toUpdate Relation) (Relation, error) {
	//oldRelation, err := s.repository.Get(ctx, toUpdate.ID)
	//if err != nil {
	//	return Relation{}, err
	//}
	//
	//newRelation, err := s.repository.Update(ctx, toUpdate)
	//if err != nil {
	//	return Relation{}, err
	//}
	//
	//if err = s.authzRepository.Delete(ctx, oldRelation); err != nil {
	//	return Relation{}, err
	//}
	//
	//if err = s.authzRepository.Add(ctx, newRelation); err != nil {
	//	return Relation{}, err
	//}
	//
	//return newRelation, nil
	return Relation{}, nil
}

func (s Service) Delete(ctx context.Context, rel Relation) error {
	//fetchedRel, err := s.repository.GetByFields(ctx, rel)
	//if err != nil {
	//	return err
	//}
	//
	//if err = s.authzRepository.Delete(ctx, rel); err != nil {
	//	return err
	//}
	//
	//return s.repository.DeleteByID(ctx, fetchedRel.ID)
	return nil
}

func (s Service) GetRelationByFields(ctx context.Context, rel RelationV2) (RelationV2, error) {
	fetchedRel, err := s.repository.GetByFields(ctx, rel)
	if err != nil {
		return RelationV2{}, err
	}

	return fetchedRel, nil
}

func (s Service) DeleteV2(ctx context.Context, rel RelationV2) error {
	fetchedRel, err := s.repository.GetByFields(ctx, rel)
	if err != nil {
		return err
	}
	if err := s.authzRepository.DeleteV2(ctx, fetchedRel); err != nil {
		return err
	}

	return s.repository.DeleteByID(ctx, fetchedRel.ID)
}

func (s Service) CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error) {
	return s.authzRepository.Check(ctx, Relation{
		ObjectNamespace:  resourceNS,
		ObjectID:         resourceIdxa,
		SubjectID:        usr.ID,
		SubjectNamespace: namespace.DefinitionUser,
	}, action)
}

func (s Service) DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	err = s.authzRepository.DeleteSubjectRelations(ctx, resourceType, optionalResourceID)
	if err != nil {
		return err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		relationSubjectlogData := ToRelationSubjectLogData(resourceType, optionalResourceID)
		if err := s.activityService.Log(ctx, auditKeyRelationSubjectDelete, currentUser.ID, relationSubjectlogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return nil
}
