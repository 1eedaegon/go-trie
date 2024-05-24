package trie

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

//	func TestTableByTrieScenarios(t *testing.T) {
//		for scenario, fn := range map[string]func(t *testing.T, trie Trier){
//			"test trie":         testTrie,
//			"test nil behavier": testNilBehavior,
//		} {
//			t.Run(scenario, func(t *testing.T) {
//				trie := NewTrie()
//				fn(t, trie)
//			})
//		}
//	}

func TestTrie(t *testing.T) {
	trie := NewTrie()
	testTrie(t, trie)
}

func TestPrefixSearchOperations(t *testing.T) {
	trie := NewTrie()
	testPrefixSearchOperations(t, trie)
}

func TestMarshalling(t *testing.T) {
	trie := NewTrie()
	testMarshallingOperations(t, trie)
}

func testMarshallingOperations(t *testing.T, trie Trier) {
	trie.Put("HappyBirthday!", 0)
	trie.Put("Happy!", 0)
	trie.Put("Brithday!", 0)
	trie.Put("pyBr", 0)
	trie.Put("HappyBrithday!", 0)

	jsonData, err := json.Marshal(trie)
	jsonString := "{\"children\":{\"B\":{\"children\":{\"r\":{\"children\":{\"i\":{\"children\":{\"t\":{\"children\":{\"h\":{\"children\":{\"d\":{\"children\":{\"a\":{\"children\":{\"y\":{\"children\":{\"!\":{\"value\":0}}}}}}}}}}}}}}}}},\"H\":{\"children\":{\"a\":{\"children\":{\"p\":{\"children\":{\"p\":{\"children\":{\"y\":{\"children\":{\"!\":{\"value\":0},\"B\":{\"children\":{\"i\":{\"children\":{\"r\":{\"children\":{\"t\":{\"children\":{\"h\":{\"children\":{\"d\":{\"children\":{\"a\":{\"children\":{\"y\":{\"children\":{\"!\":{\"value\":0}}}}}}}}}}}}}}},\"r\":{\"children\":{\"i\":{\"children\":{\"t\":{\"children\":{\"h\":{\"children\":{\"d\":{\"children\":{\"a\":{\"children\":{\"y\":{\"children\":{\"!\":{\"value\":0}}}}}}}}}}}}}}}}}}}}}}}}}}},\"p\":{\"children\":{\"y\":{\"children\":{\"B\":{\"children\":{\"r\":{\"value\":0}}}}}}}}}"
	require.NoError(t, err)
	require.Equal(t, jsonString, string(jsonData))

	newTrie := &RuneMapTrie{}
	err = json.Unmarshal(jsonData, newTrie)
	// JSON 데이터를 RuneMapTrie로 언마샬링
	require.NoError(t, err)
	require.Equal(t, trie, newTrie)
	require.True(t, reflect.DeepEqual(trie, newTrie))
}

func testPrefixSearchOperations(t *testing.T, trie Trier) {
	trie.Put("HappyBrithday!", 0)
	trie.Put("Happy!", 0)
	trie.Put("Brithday!", 0)
	trie.Put("pyBr", 0)
	trie.Put("HappyBrithday!", 0)
	result, err := trie.PrefixSearch("Happy")
	require.NoError(t, err)
	require.Contains(t, result, "HappyBrithday!")
	require.Contains(t, result, "Happy!")

}

func testTrie(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	cases := []struct {
		key   string
		value interface{}
	}{
		{"fish", 0},
		{"/cat", 1},
		{"/dog", 2},
		{"/cats", 3},
		{"/caterpillar", 4},
		{"/cat/gideon", 5},
		{"/cat/giddy", 6},
	}

	// get missing keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be missing, found value %v", c.key, value)
		}
	}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, firstPutValue); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// subsequent put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); isNew {
			t.Errorf("expected key %s to have a value already", c.key)
		}
	}

	// get
	for _, c := range cases {
		if value := trie.Get(c.key); value != c.value {
			t.Errorf("expected key %s to have value %v, got %v", c.key, c.value, value)
		}
	}

	// delete, expect Delete to return true indicating a node was nil'd
	for _, c := range cases {
		if deleted := trie.Delete(c.key); !deleted {
			t.Errorf("expected key %s to be deleted", c.key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, c := range cases {
		if deleted := trie.Delete(string(c.key[0])); deleted {
			t.Errorf("expected key %s to be cleaned by delete", string(c.key[0]))
		}
	}

	// get deleted keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be deleted, got value %v", c.key, value)
		}
	}
}

func testNilBehavior(t *testing.T, trie Trier) {
	cases := []struct {
		key   string
		value interface{}
	}{
		{"/cat", 1},
		{"/catamaran", 2},
		{"/caterpillar", nil},
	}
	expectNilValues := []string{"/", "/c", "/ca", "/caterpillar", "/other"}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// get nil
	for _, key := range expectNilValues {
		if value := trie.Get(key); value != nil {
			t.Errorf("expected key %s to have value nil, got %v", key, value)
		}
	}
}

func testTrieRoot(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	const putValue = "value"

	if value := trie.Get(""); value != nil {
		t.Errorf("expected key '' to be missing, found value %v", value)
	}
	if !trie.Put("", firstPutValue) {
		t.Error("expected key '' to be missing")
	}
	if trie.Put("", putValue) {
		t.Error("expected key '' to have a value already")
	}
	if value := trie.Get(""); value != putValue {
		t.Errorf("expected key '' to have value %v, got %v", putValue, value)
	}
	if !trie.Delete("") {
		t.Error("expected key '' to be deleted")
	}
	if value := trie.Get(""); value != nil {
		t.Errorf("expected key '' to be deleted, got value %v", value)
	}
}

func testTrieIterate(t *testing.T, trie Trier) {
	table := map[string]interface{}{
		"":             -1,
		"fish":         0,
		"/cat":         1,
		"/dog":         2,
		"/cats":        3,
		"/caterpillar": 4,
		"/notes":       30,
		"/notes/new":   31,
		"/notes/:id":   32,
	}
	// key -> times iterated
	iterated := make(map[string]int)
	for key := range table {
		iterated[key] = 0
	}

	for key, value := range table {
		if isNew := trie.Put(key, value); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	iterator := func(key string, value interface{}) error {
		// value for each iterated key is correct
		if value != table[key] {
			t.Errorf("expected key %s to have value %v, got %v", key, table[key], value)
		}
		iterated[key]++
		return nil
	}
	if err := trie.IterateAll(iterator); err != nil {
		t.Errorf("expected error nil, got %v", err)
	}

	// each key/value iterated exactly once
	for key, iteratedCount := range iterated {
		if iteratedCount != 1 {
			t.Errorf("expected key %s to be iterated exactly once, got %v", key, iteratedCount)
		}
	}
}

func testTrieIterateAllError(t *testing.T, trie Trier) {
	table := map[string]interface{}{
		"/L1/L2A":        1,
		"/L1/L2B/L3A":    2,
		"/L1/L2B/L3B/L4": 42,
		"/L1/L2B/L3C":    4,
		"/L1/L2C":        5,
	}

	iteratorError := errors.New("iterator error")
	iterated := 0

	for key, value := range table {
		trie.Put(key, value)
	}
	iterator := func(key string, value interface{}) error {
		if value == 42 {
			return iteratorError
		}
		iterated++
		return nil
	}
	if err := trie.IterateAll(iterator); err != iteratorError {
		t.Errorf("expected iterator error, got %v", err)
	}
	if len(table) == iterated {
		t.Errorf("expected nodes iterated < %d, got %d", len(table), iterated)
	}
}

func testTrieIterateAll(t *testing.T, trie Trier) {
	table := map[string]interface{}{
		"fish":             0,
		"/cat":             1,
		"/dog":             2,
		"/cats":            3,
		"/caterpillar":     4,
		"/notes":           30,
		"/notes/new":       31,
		"/notes/new/noise": 32,
	}
	// key -> times iterated
	iterated := make(map[string]int)
	for key := range table {
		iterated[key] = 0
	}

	for key, value := range table {
		if isNew := trie.Put(key, value); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	iterator := func(key string, value interface{}) error {
		// value for each iterated key is correct
		if value != table[key] {
			t.Errorf("expected key %s to have value %v, got %v", key, table[key], value)
		}
		iterated[key]++
		return nil
	}
	if err := trie.Iterate("/notes/new/noise", iterator); err != nil {
		t.Errorf("expected error nil, got %v", err)
	}

	// expect each key/value in path iterated exactly once, and not other keys
	for key, iteratedCount := range iterated {
		switch key {
		case "/notes", "/notes/new", "/notes/new/noise":
			if iteratedCount != 1 {
				t.Errorf("expected key %s to be iterated exactly once, got %v", key, iteratedCount)
			}
		default:
			if iteratedCount != 0 {
				t.Errorf("expected key %s to not be iterated, got %v", key, iteratedCount)
			}
		}
	}

	for key := range table {
		iterated[key] = 0
	}
	if err := trie.Iterate("/notes/new/nose", iterator); err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	// expect each key/value in path iterated exactly once, and not other keys
	for key, iteratedCount := range iterated {
		switch key {
		case "/notes", "/notes/new":
			if iteratedCount != 1 {
				t.Errorf("expected key %s to be iterated exactly once, got %v", key, iteratedCount)
			}
		default:
			if iteratedCount != 0 {
				t.Errorf("expected key %s to not be iterated, got %v", key, iteratedCount)
			}
		}
	}

	var foundRoot bool
	trie.Put("", "ROOT")
	trie.Iterate("/notes/new/noise", func(key string, value interface{}) error {
		if key == "" && value == "ROOT" {
			foundRoot = true
		}
		return nil
	})
	if !foundRoot {
		t.Error("did not find root")
	}
}

func testTrieIterateError(t *testing.T, trie Trier) {
	table := map[string]interface{}{
		"/L1/L2A":        1,
		"/L1/L2A/L3B":    99,
		"/L1/L2A/L3B/L4": 2,
		"/L1/L2B/L3B":    3,
		"/L1/L2C":        4,
	}

	iteratorError := errors.New("iterator error")

	iterated := make(map[string]int)
	for key := range table {
		iterated[key] = 0
	}

	for key, value := range table {
		trie.Put(key, value)
	}
	iterator := func(key string, value interface{}) error {
		if value == 99 {
			return iteratorError
		}
		iterated[key]++
		return nil
	}
	if err := trie.Iterate("/L1/L2A/L3B", iterator); err != iteratorError {
		t.Errorf("expected iterator error, got %v", err)
	}
	// expect each key/value in path, up to error and not including key with
	// value 99, iterated exactly once, and not other keys
	var iteratedTotal int
	for key, iteratedCount := range iterated {
		switch key {
		case "/L1/L2A":
			if iteratedCount != 1 {
				t.Errorf("expected key %s to be iterated exactly once, got %v", key, iteratedCount)
			}
			iteratedTotal++
		default:
			if iteratedCount != 0 {
				t.Errorf("expected key %s to not be iterated, got %v", key, iteratedCount)
			}
		}
	}
	if iteratedTotal != 1 {
		t.Errorf("expected 1 nodes iterated, got %d", iteratedTotal)
	}

	rootError := errors.New("error at root")
	trie.Put("", "ROOT")
	err := trie.Iterate("/L1/L2A/L3B/L4", func(key string, value interface{}) error {
		if key == "" && value == "ROOT" {
			return rootError
		}
		return nil
	})
	if err != rootError {
		t.Errorf("expected %s, got %s", rootError, err)
	}
}
