package workflow

import (
	"testing"
)

func TestSendToAuthor(t *testing.T) {
	wk := NewSendToAuthorWorkflow()
	if err := wk.SendToAuthor("@author", "text"); err != nil {
		t.Fatalf("failed %s", err)
	}
}
