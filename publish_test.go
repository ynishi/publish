package main

import "testing"

func TestGreet(t *testing.T) {

	testPerson := "world"
	testMsg := greet(testPerson)

	msg :="hello, world!"
	if msg != testMsg {
		t.Fatalf("Greet Message not match.\n want: %q,\n have: %q\n", msg, testMsg)
	}
}
