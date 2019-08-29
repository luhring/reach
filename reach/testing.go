package reach

import "testing"

func diffErrorf(t *testing.T, item string, expected, actual interface{}) {
	t.Helper()
	t.Errorf("'%s' value differed from expected value...\n\nexpected:\n%v\n\nactual:\n%v\n\n", item, expected, actual)
}
