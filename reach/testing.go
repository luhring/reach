package reach

import "testing"

// DiffErrorf provides a convenient way to output a difference between two values (such as between an expected value and an actual value) that caused a test to fail.
func DiffErrorf(t *testing.T, item string, expected, actual interface{}) {
	t.Helper()
	t.Errorf("'%s' value differed from expected value...\n\nexpected:\n%v\n\nactual:\n%v\n\n", item, expected, actual)
}
