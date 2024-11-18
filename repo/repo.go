package repo

import (
	"github.com/google/uuid"
	functionservice "github.com/nilspolek/AstralFS/function-service"
)

type Repo interface {
	InsertFunction(functionservice.Function) error
	DeleteFunction(uuid.UUID) error
	DeleteAllFunctions() error
	GetFunctions() ([]functionservice.Function, error)
}

type RepoFunctionService struct {
	repo *Repo
	next *functionservice.FunctionService
}

func New(repo *Repo, next functionservice.FunctionService) (functionservice.FunctionService, error) {
	return &RepoFunctionService{
		repo: repo,
		next: &next,
	}, nil
}

func (rfs *RepoFunctionService) CreateFunction(fn functionservice.Function) (int, error) {
	(*rfs.repo).InsertFunction(fn)
	return (*rfs.next).CreateFunction(fn)
}

func (rfs *RepoFunctionService) DeleteFunction(id uuid.UUID) error {
	(*rfs.repo).DeleteFunction(id)
	return (*rfs.next).DeleteFunction(id)
}

func (rfs *RepoFunctionService) GetFunctions() ([]functionservice.Function, error) {
	return (*rfs.next).GetFunctions()
}

func (rfs *RepoFunctionService) Close() error {
	return (*rfs.next).Close()
}
