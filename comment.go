package necovalidate

import (
	"errors"
	"fmt"
	"strings"
)

const ExpressPrefix = "@"

func isExpComment(comment string) bool {
	return strings.HasPrefix(trimComment(comment), ExpressPrefix)
}

func parseComment(comment string) (Expression, error) {
	var expression Expression
	fields := strings.Fields(trimComment(comment))
	if len(fields) == 0 {
		return expression, errors.New("expression can't be empty")
	}

	if !strings.HasPrefix(fields[0], ExpressPrefix) {
		return expression, errors.New(
			fmt.Sprintf("comment [%s]: expression name should start with %s",
				comment,
				ExpressPrefix,
			),
		)
	}

	expression.Name = fields[0]
	expression.Args = fields[1:]
	return expression, nil
}

func trimComment(comment string) string {
	return strings.TrimSpace(strings.TrimLeft(comment, "/"))
}
