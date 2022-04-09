package graph

import (
	"fmt"

	"github.com/encero/reciper/api"
	"github.com/encero/reciper/gql/graph/model"
)

func statusToResult(status string) (*model.Result, error) {
	switch status {
	case api.StatusSuccess:
		return &model.Result{Status: model.StatusSuccess}, nil
	case api.StatusNotFound:
		return &model.Result{Status: model.StatusNotFound}, nil
	case api.StatusError:
		return &model.Result{Status: model.StatusError}, nil
	default:
		return nil, fmt.Errorf("unknown status: %s", status)
	}
}
