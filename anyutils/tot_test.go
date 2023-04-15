package anyutils

import "testing"

func ObjectTreeForTesting001_float() any {
	ret := make(map[string]map[string]map[string]any)
	ret["a"] = make(map[string]map[string]any)
	ret["a"]["b"] = make(map[string]any)
	ret["a"]["b"]["c"] = 1.2
	return ret
}

func ObjectTreeForTesting001_string() any {
	ret := make(map[string]map[string]map[string]any)
	ret["a"] = make(map[string]map[string]any)
	ret["a"]["b"] = make(map[string]any)
	ret["a"]["b"]["c"] = "teststring"
	return ret
}

func ObjectTreeForTesting001_int64() any {
	ret := make(map[string]map[string]map[string]any)
	ret["a"] = make(map[string]map[string]any)
	ret["a"]["b"] = make(map[string]any)
	ret["a"]["b"]["c"] = int64(-64)
	return ret
}

func ObjectTreeForTesting001_int32() any {
	ret := make(map[string]map[string]map[string]any)
	ret["a"] = make(map[string]map[string]any)
	ret["a"]["b"] = make(map[string]any)
	ret["a"]["b"]["c"] = int64(-32)
	return ret
}

func ObjectTreeForTesting001_uint64() any {
	ret := make(map[string]map[string]map[string]any)
	ret["a"] = make(map[string]map[string]any)
	ret["a"]["b"] = make(map[string]any)
	ret["a"]["b"]["c"] = uint64(64)
	return ret
}

func TestTraverseObjectTree002_001(t *testing.T) {

	tree := ObjectTreeForTesting001_string()
	value, found, err := TraverseObjectTree002_string(
		tree, true, false, "a", "b", "c",
	)

	if err != nil {
		t.Fail()
		return
	}

	if !found {
		t.Fail()
		return
	}

	if value != "teststring" {
		t.Fail()
		return
	}

	return
}

func TestTraverseObjectTree002_002(t *testing.T) {

	tree := ObjectTreeForTesting001_string()
	_, found, err := TraverseObjectTree002_string(
		tree, true, true, "a", "b", "d",
	)

	if err != nil {
		t.Fail()
		return
	}

	if found {
		t.Fail()
		return
	}

	return
}

func TestTraverseObjectTree002_003(t *testing.T) {

	tree := ObjectTreeForTesting001_string()
	_, found, err := TraverseObjectTree002_string(
		tree, true, true, "a", "z", "d",
	)

	if err != nil {
		t.Fail()
		return
	}

	if found {
		t.Fail()
		return
	}

	return
}

func TestTraverseObjectTree002_004(t *testing.T) {

	tree := ObjectTreeForTesting001_int32()
	value, found, err := TraverseObjectTree002_int64(
		tree, true, false, "a", "b", "c",
	)

	if err != nil {
		t.Error("err:", err)
		return
	}

	if !found {
		t.Error("not found")
		return
	}

	if value != -32 {
		t.Error("value != 32:", value)
		return
	}

	return
}
