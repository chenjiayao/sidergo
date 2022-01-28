package list

import (
	"testing"
)

//两个接口值相等仅当它们都是nil值或者它们的动态类型相同并且动态值也根据这个动态类型的＝=操作相等。
func TestList_InsertTail(t *testing.T) {
	l := MakeList()
	l.InsertTail(1)
	l.InsertTail(2)
	l.InsertTail(3)
}

func TestList_Exist(t *testing.T) {
	l := MakeList()
	l.InsertTail(1)
	l.InsertTail(2)
	l.InsertTail(3)

	exist := l.Exist(3)
	if !exist {
		t.Errorf("l.Exist(3) = %v, but want true", exist)
	}

	exist = l.Exist(5)
	if exist {
		t.Errorf("l.Exist(5) = %v, but want false", exist)
	}

	exist = l.Exist(1)
	if !exist {
		t.Errorf("l.Exist(1) = %v, but want true", exist)
	}
}

func TestList_RemoveNode(t *testing.T) {
	l := MakeList()
	l.InsertTail(1)
	l.InsertTail(2)
	l.InsertTail(3)

	exist := l.Exist(2)
	if !exist {
		t.Errorf("l.Exist(2) = %v, but want true", exist)
	}

	l.RemoveNode(2)
	exist = l.Exist(2)
	if exist {
		t.Errorf("l.RemoveNode(2) = falied")
	}

	l.RemoveNode(1)
	exist = l.Exist(1)
	if exist {
		t.Errorf("l.RemoveNode(1) = falied")
	}
	l.RemoveNode(6)
}

func TestList_GetElementByIndex(t *testing.T) {
	l := MakeList() // 4 3 2 1
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4)

	v := l.GetElementByIndex(0)
	got := v.(int)
	if got != 4 {
		t.Errorf("l.GetElementByIndex(0) = %d, want 4", got)
	}

	v = l.GetElementByIndex(9)
	if v != nil {
		t.Errorf("l.GetElementByIndex(9) = %v, want nil", got)
	}

	v = l.GetElementByIndex(3)
	got = v.(int)
	if got != 1 {
		t.Errorf("l.GetElementByIndex(3) = %d, want 1", got)
	}

	v = l.GetElementByIndex(-1)
	got = v.(int)
	if got != 1 {
		t.Errorf("l.GetElementByIndex(-1) = %d, want 1", got)
	}

	v = l.GetElementByIndex(-4)
	got = v.(int)
	if got != 4 {
		t.Errorf("l.GetElementByIndex(-4) = %d, want 4", got)
	}

	v = l.GetElementByIndex(-10)
	if v != nil {
		t.Errorf("l.GetElementByIndex(-10) = %d, want 1", got)
	}
}

func TestList_Range(t *testing.T) {
	l := MakeList()
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4) // 4 3 2 1

	v := l.Range(0, 0)
	v0 := v[0].(int)
	if v0 != 4 {
		t.Errorf(" l.Range(0, 0) = [%d], want [4]", v0)
	}
	//////////////
	v = l.Range(0, 4)
	got := make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want := []int{4, 3, 2, 1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(0, 4) = %v, want %v", got, want)
	}

	////////////////
	v = l.Range(1, 4)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{3, 2, 1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(1, 4) = %v, want %v", got, want)
	}

	////////////////
	v = l.Range(-1, 4)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(-1, 4) = %v, want %v", got, want)
	}

	////////////////
	v = l.Range(-1, -4)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(-1, -4) = %v, want %v", got, want)
	}

	/////////
	v = l.Range(3, -1)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(3, -1) = %v, want %v", got, want)
	}

	/////////
	v = l.Range(3, 1)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(3, -1) = %v, want %v", got, want)
	}

	/////////
	v = l.Range(3, -3)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(3, -3) = %v, want %v", got, want)
	}

	/////////
	v = l.Range(-3, 3)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{3, 2, 1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(-3, 3) = %v, want %v", got, want)
	}

	/////////
	v = l.Range(-3, 30)
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{3, 2, 1}
	if !SliceEqual(got, want) {
		t.Errorf("l.Range(-3, 30) = %v, want %v", got, want)
	}
}

func SliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestList_PopFromTail(t *testing.T) {
	l := MakeList() // 4 3 2 1

	got := l.PopFromTail()
	if got != nil {
		t.Errorf("l.PopFromTail() want null")
	}

	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4)
	got = l.PopFromTail()
	if got == nil {
		t.Errorf("l.PopFromTail() got null")
	}

	v := got.(int)
	want := 1
	if v != 1 {
		t.Errorf("l.PopFromTail() = %d, want %d", v, want)
	}

	////
	got = l.PopFromTail()
	if got == nil {
		t.Errorf("l.PopFromTail() got null")
	}

	v = got.(int)
	want = 2
	if v != want {
		t.Errorf("l.PopFromTail() = %d, want %d", v, want)
	}
}

func TestList_Trim(t *testing.T) {
	l := MakeList()
	l.InsertHead(0)
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4) // 4 3 2 1 0

	l.Trim(1, 2)
	v := l.Range(0, l.Len())
	got := make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want := []int{
		3, 2,
	}
	if !SliceEqual(got, want) {
		t.Errorf("l.Trim(1, 2) failed, got =%v", got)
	}

	///////////
	l = MakeList()
	l.InsertHead(0)
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4) // 4 3 2 1 0

	l.Trim(10, 2)
	v = l.Range(0, l.Len())
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Trim(10, 2) failed, got =%v", got)
	}

	///////////
	l = MakeList()
	l.InsertHead(0)
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4) // 4 3 2 1 0

	l.Trim(-1, 2)
	v = l.Range(0, l.Len())
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Trim(-1, 2) failed, got =%v", got)
	}

	///////////
	l = MakeList()
	l.InsertHead(0)
	l.InsertHead(1)
	l.InsertHead(2)
	l.InsertHead(3)
	l.InsertHead(4) // 4 3 2 1 0

	l.Trim(-1, -2)
	v = l.Range(0, l.Len())
	got = make([]int, len(v))
	for i := 0; i < len(v); i++ {
		got[i] = v[i].(int)
	}
	want = []int{}
	if !SliceEqual(got, want) {
		t.Errorf("l.Trim(-1, -2) failed, got =%v", got)
	}
}
