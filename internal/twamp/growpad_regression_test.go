package twamp

// Regression coverage for GrowPad aliasing: when dst still has spare
// capacity, the grown slice must not share dst's backing array, otherwise
// writes to the new slice pollute the original pad buffer.

import "testing"

func TestGrowPadRegressionSpareCapacity(t *testing.T) {
	// dst has 6 bytes of spare capacity: append(dst, extra) would normally
	// reuse dst's backing array and alias it.
	dst := make([]byte, 2, 8)
	copy(dst, []byte("AB"))

	got := GrowPad(dst, 'C')

	// Content is correct.
	if string(got) != "ABC" {
		t.Fatalf("got %q, want %q", got, "ABC")
	}

	// Mutate the grown slice in place; the original pad buffer must be
	// untouched both in content and in length.
	got[0] = 'X'
	got[1] = 'Y'
	got[2] = 'Z'
	if dst[0] != 'A' || dst[1] != 'B' {
		t.Fatalf("dst polluted by write through: %q", dst)
	}
	if len(dst) != 2 {
		t.Fatalf("dst length changed: %d", len(dst))
	}

	// The grown slice owns its own backing array: appending to it must not
	// touch dst even though dst still has spare capacity.
	more := append(got, 'D')
	dst[0] = 'A' // restore (no-op) and ensure no back-propagation
	if more[0] != 'X' || more[1] != 'Y' || more[2] != 'Z' || more[3] != 'D' {
		t.Fatalf("unexpected grown slice state: %q", more)
	}
}

func TestGrowPadRegressionEmpty(t *testing.T) {
	got := GrowPad(nil, 'Q')
	if len(got) != 1 || got[0] != 'Q' {
		t.Fatalf("got %v, want [Q]", got)
	}
}
