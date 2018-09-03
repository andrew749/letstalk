package query

import (
	"context"
	"fmt"

	"letstalk/server/core/search"
	"letstalk/server/data"

	"github.com/getsentry/raven-go"
	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

func withEsBackgroundContext(
	es *elastic.Client,
	actionStr string,
	fn func(c *search.ClientWithContext) error,
) {
	if es != nil {
		searchClient := search.NewClientWithContext(es, context.Background())
		err := fn(searchClient)
		if err != nil {
			raven.CaptureError(err, nil)
			rlog.Error(err)
		}
	} else {
		rlog.Warn(fmt.Sprintf("Not %s since no es provided", actionStr))
	}
}

func indexSimpleTraitMultiTrait(es *elastic.Client, trait *data.UserSimpleTrait) {
	withEsBackgroundContext(
		es,
		fmt.Sprintf("indexing simple trait multi trait %s", trait.SimpleTraitName),
		func(c *search.ClientWithContext) error {
			return c.IndexSimpleTraitMultiTrait(trait)
		},
	)
}

func indexCohortMultiTrait(es *elastic.Client, cohort *data.UserCohort) {
	if cohort.Cohort == nil {
		rlog.Warn(
			fmt.Sprintf("Cohort was not provided when adding cohort with id %d", cohort.CohortId),
		)
		return
	}
	var sequenceName string
	if cohort.Cohort.SequenceName != nil {
		sequenceName = *cohort.Cohort.SequenceName
	} else {
		sequenceName = "Non-coop"
	}

	withEsBackgroundContext(
		es,
		fmt.Sprintf(
			"indexing cohort multi trait \"%s %d %s\"",
			cohort.Cohort.ProgramName,
			cohort.Cohort.GradYear,
			sequenceName,
		),
		func(c *search.ClientWithContext) error {
			return c.IndexCohortMultiTrait(cohort)
		},
	)
}

func indexPositionMultiTrait(es *elastic.Client, pos *data.UserPosition) {
	withEsBackgroundContext(
		es,
		fmt.Sprintf("indexing position multi trait \"%s at %s\"", pos.RoleName, pos.OrganizationName),
		func(c *search.ClientWithContext) error {
			return c.IndexPositionMultiTrait(pos)
		},
	)
}

func indexSimpleTrait(es *elastic.Client, trait data.SimpleTrait) {
	withEsBackgroundContext(
		es,
		fmt.Sprintf("indexing simple trait %s", trait.Name),
		func(c *search.ClientWithContext) error {
			searchTrait := search.NewSimpleTraitFromDataModel(trait)
			return c.IndexSimpleTrait(searchTrait)
		},
	)
}

func indexRole(es *elastic.Client, role data.Role) {
	withEsBackgroundContext(
		es,
		fmt.Sprintf("indexing role %s", role.Name),
		func(c *search.ClientWithContext) error {
			searchRole := search.NewRoleFromDataModel(role)
			return c.IndexRole(searchRole)
		},
	)
}

func indexOrganization(es *elastic.Client, organization data.Organization) {
	withEsBackgroundContext(
		es,
		fmt.Sprintf("indexing organization %s", organization.Name),
		func(c *search.ClientWithContext) error {
			searchOrganization := search.NewOrganizationFromDataModel(organization)
			return c.IndexOrganization(searchOrganization)
		},
	)
}
