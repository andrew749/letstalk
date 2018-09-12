package search

import (
	"context"
	"errors"
	"fmt"

	"letstalk/server/core/errs"

	"github.com/olivere/elastic"
)

// Used for indexes that require completion suggestions
type SuggestInput struct {
	Input  []string `json:"input"`
	Weight *int     `json:"weight,omitempty"`
}

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

func NewDefaultEsClient() (*elastic.Client, error) {
	return elastic.NewClient()
}

// Search client to be used within the request context
type ClientWithContext struct {
	client  *elastic.Client
	context context.Context
}

func NewClientWithContext(client *elastic.Client, context context.Context) *ClientWithContext {
	return &ClientWithContext{client, context}
}

func (c *ClientWithContext) CreateEsIndexes() error {
	var compErr *errs.CompositeError = nil
	compErr = errs.AppendNullableError(compErr, c.createSimpleTraitIndex())
	compErr = errs.AppendNullableError(compErr, c.createRoleIndex())
	compErr = errs.AppendNullableError(compErr, c.createOrganizationIndex())
	compErr = errs.AppendNullableError(compErr, c.createMultiTraitIndex())
	return compErr
}

func consolidateBulkResponseErrors(res *elastic.BulkResponse) error {
	if len(res.Failed()) > 0 {
		results := make([]string, len(res.Failed()))
		reasons := make([]string, len(res.Failed()))
		for i, failed := range res.Failed() {
			results[i] = failed.Result
			if failed.Error != nil {
				reasons[i] = failed.Error.Reason
			} else {
				reasons[i] = "unknown reason"
			}
		}

		return errors.New(fmt.Sprintf("More than 0 operations failed: %d, results: %#v, reasons: %#v\n",
			len(reasons), results, reasons))
	}
	return nil
}
