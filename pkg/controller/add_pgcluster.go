package controller

import (
	"github.com/jasonodonnell/example-operator/pkg/controller/pgcluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pgcluster.Add)
}
