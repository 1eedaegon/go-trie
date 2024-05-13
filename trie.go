package trie

type TrieType int

const (
	TypeRuneMapTrie TrieType = 0
	TypeBitmapTrie  TrieType = 1
)

// Trie has 3 methods
// 1. Get by key
// 2. Put with key and value
// 3. Delete by key

type Trier interface {
	Get(key string) interface{}
	Put(key string, value interface{}) bool
	Delete(key string) bool
	Iterate(key string, cb Callback) error
	IterateAll(cb Callback) error
}

func NewTrie(trieType TrieType) Trier {
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
