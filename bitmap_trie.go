package trie

type BitmapTrie struct {
	leaves, labelBitmap []uint64
	labels              []byte
	ranks, selects      []int32
}

var _ Trier = (*BitmapTrie)(nil)

func NewBitmapTrie() *BitmapTrie {
	return new(BitmapTrie)
}

func (trie *BitmapTrie) Get(key string) interface{} {
	return nil
}
func (trie *BitmapTrie) Put(key string, value interface{}) bool {
	return false
}
func (trie *BitmapTrie) Delete(key string) bool {
	return false
}
func (trie *BitmapTrie) Iterate(key string, cb Callback) error {
	return nil
}
func (trie *BitmapTrie) IterateAll(cb Callback) error {
	return nil
}

func (trie *BitmapTrie) PrefixSearch(key string) ([]string, error) {
	return nil, nil
}
