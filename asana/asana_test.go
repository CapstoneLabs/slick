package asana

import (
	"strings"
	"testing"
)

func TestStringID(t *testing.T) {
	var t1 = Tag{ID: 120, Name: "test1"}
	if t1.StringID() != "120" {
		t.Errorf("Function StringID - ERROR")
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("Test", "TestWorkSpace")
	if c.workspace != "TestWorkSpace" && c.key != "Test" {
		t.Errorf("Function NewClient - ERROR")
	}
}

func TestSetWorkspace(t *testing.T) {
	n := NewClient("Test", "TestWorkSpace")
	n.SetWorkspace("New Test")
	if n.workspace != "New Test" {
		t.Errorf("Function SetWorkspace - ERROR")
	}
}

func TestRequest(t *testing.T) {
	n := NewClient("Test", "TestWorkSpace")
	//testing with https://postman-echo.com/put
	b, _ := n.Request("PUT", "https://postman-echo.com/put", "Put method test")
	if !strings.Contains(string(b), "Put method test") {
		t.Errorf("Function request - ERROR")
	}
}
