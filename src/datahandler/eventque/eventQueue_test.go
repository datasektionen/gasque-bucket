package eventque
import ("testing")

func dummyData(d int) *DisplayContent {
	dummy := new(DisplayContent)
	dummy.Duration = d
	return dummy
}
func assertEqual(t *testing.T, s string, i1 int, i2 int) {
	if i1 != i2 {
		t.Error(s)
	}
}
func TestEventQueue(t *testing.T) {
	e := new(EventQueue)

	if e.Next() != nil {
		t.Error("Expected nil")
	}
	if e.Top() != nil {
		t.Error("Excepted nil")
	}
	if e.Size() != 0 {
		t.Error("Size error")
	}
	e.Push(dummyData(1))
	e.Push(dummyData(2))
	e.Push(dummyData(3))
	assertEqual(t, "Expected other data", e.Top().Duration, 1)
	assertEqual(t, "Expected other data", e.Next().Duration, 1)
	assertEqual(t, "Expected other data", e.Next().Duration, 2)
	assertEqual(t, "Expected other data", e.Next().Duration, 3)
	
	if e.Next() != nil {
		t.Error("Expected nil")
	}
	if e.Size() != 0 {
		t.Error("Expected size 0")
	}
}
