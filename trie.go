package trie

type TrieType int

const (
	TypeRuneMapTrie TrieType = 0
	TypeBitmapTrie  TrieType = 1
)

// Trie has 6 methods
// 1. Get by key
// 2. Put with key and value
// 3. Delete by key
// 4. Iterate with callbck by key
// 5. Iterate all element with callback
// 6. Search for a prefix, it returns an slice of string
type Trier interface {
	Get(key string) interface{}
	Put(key string, value interface{}) bool
	Delete(key string) bool
	Iterate(key string, cb Callback) error
	IterateAll(cb Callback) error
	PrefixSearch(key string) ([]string, error)
}

func NewTrie(trieOption ...interface{}) Trier {
	trieType := TypeRuneMapTrie
	if len(trieOption) > 0 {
		t := trieOption[0].(int)
		trieType = TrieType(t)
	}
	switch trieType {
	case TypeBitmapTrie:
		return NewBitmapTrie()
	case TypeRuneMapTrie:
		return NewRuneMapTrie()
	default:
		return NewRuneMapTrie()
	}
}

type Callback func(key string, value interface{}) error
