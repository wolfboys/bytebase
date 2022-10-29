package mysql

// Framework code is generated by the generator.

import (
	"testing"

	"github.com/bytebase/bytebase/plugin/advisor"
)

func TestStatementDmlDryRun(t *testing.T) {
	tests := []advisor.TestCase{
		{
			Statement: `INSERT INTO tech_book values(1, 'a')`,
			Want: []advisor.Advice{
				{
					Status:  advisor.Success,
					Code:    advisor.Ok,
					Title:   "OK",
					Content: "",
				},
			},
		},
		{
			Statement: `DELETE FROM tech_book`,
			Want: []advisor.Advice{
				{
					Status:  advisor.Success,
					Code:    advisor.Ok,
					Title:   "OK",
					Content: "",
				},
			},
		},
	}

	advisor.RunSQLReviewRuleTests(t, tests, &StatementDmlDryRunAdvisor{}, &advisor.SQLReviewRule{
		Type:    advisor.SchemaRuleStatementDMLDryRun,
		Level:   advisor.SchemaRuleLevelWarning,
		Payload: "",
	}, advisor.MockMySQLDatabase)
}