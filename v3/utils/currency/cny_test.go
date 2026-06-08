package currency

import "testing"

func TestCentsToYuan(t *testing.T) {
	cases := []struct {
		cents int64
		want  string
	}{
		{0, "0.00"},
		{1, "0.01"},
		{10, "0.10"},
		{100, "1.00"},
		{101, "1.01"},
		{12345, "123.45"},
		{99, "0.99"},
		{-1, "-0.01"},
		{-100, "-1.00"},
		{-12345, "-123.45"},
	}
	for _, c := range cases {
		if got := CentsToYuan(c.cents); got != c.want {
			t.Errorf("CentsToYuan(%d) = %q, want %q", c.cents, got, c.want)
		}
	}
}
