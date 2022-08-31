package mysql

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/pingcap/tidb/parser/ast"

	"github.com/bytebase/bytebase/plugin/advisor"
	"github.com/bytebase/bytebase/plugin/advisor/db"
)

var (
	_ advisor.Advisor = (*IndexNoDuplicateColumnAdvisor)(nil)
	_ ast.Visitor     = (*indexNoDuplicateColumnChecker)(nil)
)

func init() {
	advisor.Register(db.MySQL, advisor.MySQLIndexNoDuplicateColumn, &IndexNoDuplicateColumnAdvisor{})
	advisor.Register(db.TiDB, advisor.MySQLIndexNoDuplicateColumn, &IndexNoDuplicateColumnAdvisor{})
}

// IndexNoDuplicateColumnAdvisor is the advisor checking for no duplicate columns in index.
type IndexNoDuplicateColumnAdvisor struct {
}

// Check checks for no duplicate columns in index.
func (*IndexNoDuplicateColumnAdvisor) Check(ctx advisor.Context, statement string) ([]advisor.Advice, error) {
	stmtList, errAdvice := parseStatement(statement, ctx.Charset, ctx.Collation)
	if errAdvice != nil {
		return errAdvice, nil
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &indexNoDuplicateColumnChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type indexNoDuplicateColumnChecker struct {
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
	text       string
	line       int
}

// Enter implements the ast.Visitor interface.
func (checker *indexNoDuplicateColumnChecker) Enter(in ast.Node) (ast.Node, bool) {
	type duplicateColumn struct {
		table  string
		index  string
		column string
		line   int
		tp     string
	}
	var columnList []duplicateColumn
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range node.Constraints {
			switch constraint.Tp {
			case ast.ConstraintPrimaryKey,
				ast.ConstraintUniq,
				ast.ConstraintUniqIndex,
				ast.ConstraintIndex,
				ast.ConstraintForeignKey:
				if column, duplicate := hasDuplicateColumn(constraint.Keys); duplicate {
					columnList = append(columnList, duplicateColumn{
						tp:     indexTypeString(constraint.Tp),
						table:  node.Table.Name.O,
						index:  constraint.Name,
						column: column,
						line:   constraint.OriginTextPosition(),
					})
				}
			}
		}
	case *ast.CreateIndexStmt:
		if column, duplicate := hasDuplicateColumn(node.IndexPartSpecifications); duplicate {
			columnList = append(columnList, duplicateColumn{
				tp:     "INDEX",
				table:  node.Table.Name.O,
				index:  node.IndexName,
				column: column,
				line:   checker.line,
			})
		}
	case *ast.AlterTableStmt:
		for _, spec := range node.Specs {
			if spec.Tp == ast.AlterTableAddConstraint {
				switch spec.Constraint.Tp {
				case ast.ConstraintPrimaryKey,
					ast.ConstraintUniq,
					ast.ConstraintUniqIndex,
					ast.ConstraintIndex,
					ast.ConstraintForeignKey:
					if column, duplicate := hasDuplicateColumn(spec.Constraint.Keys); duplicate {
						columnList = append(columnList, duplicateColumn{
							tp:     indexTypeString(spec.Constraint.Tp),
							table:  node.Table.Name.O,
							index:  spec.Constraint.Name,
							column: column,
							line:   checker.line,
						})
					}
				}
			}
		}
	}

	for _, column := range columnList {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  checker.level,
			Code:    advisor.DuplicateColumnInIndex,
			Title:   checker.title,
			Content: fmt.Sprintf("%s `%s` has duplicate column `%s`.`%s`", column.tp, column.index, column.table, column.column),
			Line:    column.line,
		})
	}

	return in, false
}

// Leave implements the ast.Visitor interface.
func (*indexNoDuplicateColumnChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func hasDuplicateColumn(keyList []*ast.IndexPartSpecification) (string, bool) {
	checker := make(map[string]bool)
	for _, key := range keyList {
		if key.Expr == nil {
			if _, exists := checker[key.Column.Name.O]; exists {
				return key.Column.Name.O, true
			}
			checker[key.Column.Name.O] = true
		}
	}

	return "", false
}

func indexTypeString(tp ast.ConstraintType) string {
	switch tp {
	case ast.ConstraintPrimaryKey:
		return "PRIMARY KEY"
	case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
		return "UNIQUE KEY"
	case ast.ConstraintForeignKey:
		return "FOREIGN KEY"
	case ast.ConstraintIndex:
		return "INDEX"
	}
	return "INDEX"
}
