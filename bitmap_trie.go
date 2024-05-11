package trie

type BitTrie struct {
	value    interface{}
	children map[rune]*BitTrie
}

var _ Trier = (*BitTrie)(nil)

func (trie *BitTrie) Get(key string) interface{} {
	return nil
}
func (trie *BitTrie) Put(key string, value interface{}) bool {
	return false
}
func (trie *BitTrie) Delete(key string) bool {
	return false
}
