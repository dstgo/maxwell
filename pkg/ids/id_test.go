package ids

import "testing"

func testId(n int, f func() string) []string {
	m := make(map[string]int)
	for range n {
		m[f()]++
	}

	var repeat []string
	for k, v := range m {
		if v > 1 {
			repeat = append(repeat, k)
		}
	}
	return repeat
}

func TestULId(t *testing.T) {
	ids := testId(1000_0000, ULID)
	for _, id := range ids {
		t.Log(id)
	}
}

func TestUUId(t *testing.T) {
	ids := testId(1000_0000, UUID)
	for _, id := range ids {
		t.Log(id)
	}
}
