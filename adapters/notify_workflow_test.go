package adapters

import (
	"testing"
)

func TestSendToAuthor(t *testing.T) {
	sender := CreateSender()
	defer sender.CloseSender()
	err := sender.SendToAuthor("@author", "text")
	if err != nil {
		t.Fatalf("failed %s", err)
	}
}
