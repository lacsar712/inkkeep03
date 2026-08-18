package printtac

import "testing"

func TestTACAccept(t *testing.T) {
	if err := Enforce("WO-1", "C=80 M=60 Y=55 K=70", []string{"C"}); err != nil {
		t.Fatal(err)
	}
}

func TestTACReject(t *testing.T) {
	err := Enforce("WO-1", "C=100 M=100 Y=100 K=100", []string{"C"})
	if err == nil {
		t.Fatal("expected TAC reject")
	}
}
