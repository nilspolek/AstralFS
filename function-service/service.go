package functionservice

import (
	"github.com/google/uuid"
)

type FunctionService interface {
	CreateFunction(Function) (int, error)
	DeleteFunction(uuid.UUID) error
	GetFunctions() ([]Function, error)
	Close() error
}

type Function struct {
	Image string `json:"image"`
	Port  int    `json:"port"`
	Route string `json:"route"`
	Id    uuid.UUID
}
