package can

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import "testing"

// ---------------------------------------------------------------------------------------
//  tests
// ---------------------------------------------------------------------------------------

func TestFrame_ToUint64(t *testing.T) {
	msg := Frame{Length: 2, Data: [8]byte{1, 0, 0, 0, 0, 0, 0, 0}}
	if msg.ToUint64() != 1 {
		t.Errorf("ToUint64() should be 1 but was %d", msg.ToUint64())
		return
	}

	msg = Frame{Length: 2, Data: [8]byte{4, 0, 0, 0xff, 0, 0, 0, 0}}
	if msg.ToUint64() != 4 {
		t.Errorf("ToUint64() should be 0 but was %d: length attribut not considered?", msg.ToUint64())
		return
	}

	msg = Frame{Length: 2, Data: [8]byte{4, 1, 0, 0, 0, 0, 0, 0}}
	if msg.ToUint64() != 260 {
		t.Errorf("ToUint64() should be 260 but was %d: length attribut not considered?", msg.ToUint64())
		return
	}
}
