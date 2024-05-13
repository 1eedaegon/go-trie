package trie

type TrieType int

const (
	TypeRuneMapTrie TrieType = 0
	TypeBitTrie     TrieType = 1
)

// Trie has 3 methods
// 1. Get by key
// 2. Put with key and value
// 3. Delete by key

type Trier interface {
	Get(key string) interface{}
	Put(key string, value interface{}) bool
	Delete(key string) bool
}

func NewTrie(trieType TrieType) Trier {
	switch trieType {
	case TypeBitTrie:
		return NewBitTrie()
	case TypeRuneMapTrie:
		return NewRuneMapTrie()
	default:
		return NewRuneMapTrie()
	}
}

type Callback func(key string, value interface{}) error
