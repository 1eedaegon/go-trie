package trie

/*
RuneMapTrie is a trie of runes with string keys and interface{} values.
random search speed is fast, but memory efficiency and insert operation efficiency
are relatively low.
*/
type RuneMapTrie struct {
	value    interface{}
	children map[rune]*RuneMapTrie
}

var _ Trier = (*RuneMapTrie)(nil)

// NewMapTrie allocates and returns a new *RuneMapTrie.
func NewRuneMapTrie() *RuneMapTrie {
	return &RuneMapTrie{}
}

// NewMapTrie allocates and returns a new *map of rune
func newRuneMap() map[rune]*RuneMapTrie {
	return map[rune]*RuneMapTrie{}
}

// Sequential search while traversing a string exactly
func (trie *RuneMapTrie) Get(key string) interface{} {
	node := trie
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return nil
		}
	}
	return node.value
}

/*
Put operation inserts a value into the trie by key, replacing an existing item.
Put returns true if it a new value, or false if it places an existing value.
The default value is nil.
*/
func (trie *RuneMapTrie) Put(key string, value interface{}) bool {
	node := trie
	for _, r := range key {
		if node.children == nil {
			node.children = newRuneMap()
		}
		if node.children[r] == nil {
			node.children[r] = NewRuneMapTrie()
		}
		node = node.children[r]
	}
	isNew := node.value == nil
	node.value = value
	return isNew
}

type nodeRune struct {
	r    rune
	node *RuneMapTrie
}

/*
Delete operation removes the value for a key.
If it is a leaf node, the node is removed from the tree.
And if the parent node is a leaf node, children are nilized.
*/
func (trie *RuneMapTrie) Delete(key string) bool {
	path := make([]nodeRune, len(key)) // Rune 탐색 후 캐싱
	node := trie                       // 탐색하면서 이쪽이 변경된다. 탐색이 끝나면 마지막노드가 된다.
	for i, r := range key {
		path[i] = nodeRune{r: r, node: node}
		node = node.children[r]
		if node == nil { // 다음 rune에 대해 하위를 탐색해야하지만 없다.
			return false
		}
	}
	node.value = nil // 우선 값을 지운다.
	if node.isLeaf() {
		// 탐색이 끝난 마지막 노드가 leaf면, 부모의 children에서 나를 삭제한다
		// 가장 마지막 노드를 가리키기 때문에 nodeRune의 역순으로 탐색한다
		for i := len(key) - 1; i >= 0; i-- {
			parent := path[i].node
			r := path[i].r
			delete(parent.children, r)
			if !parent.isLeaf() { // 부모가 leaf가 아니면 지우는 iteration을 종료한다.
				break
			}
			parent.children = nil
			if parent.value != nil {
				break
			}
		}
	}
	return true
}

func (trie *RuneMapTrie) isLeaf() bool {
	return len(trie.children) == 0
}

/*
Iterates over callback within the given key string.
If value has null, skip iterate with callback
*/
func (trie *RuneMapTrie) Iterate(key string, cb Callback) error {
	if trie.value != nil {
		if err := cb("", trie.value); err != nil {
			return err
		}
	}

	for idx, r := range key {
		if trie = trie.children[r]; trie == nil {
			return nil
		}
		if trie.value != nil {
			prefix := string(key[0 : idx+1])
			if err := cb(prefix, trie.value); err != nil {
				return err
			}
		}
	}

	return nil
}

// Iterate all trie element recursive dfs by key
func (trie *RuneMapTrie) IterateAll(cb Callback) error {
	return trie.dfs("", cb)
}

/*
DFS
Implements dfs by searching the children of a node first.
If value has null, skip iterate with callback
*/
func (trie *RuneMapTrie) dfs(key string, cb Callback) error {
	if trie.value != nil {
		if err := cb(key, trie.value); err != nil {
			return err
		}
	}
	for r, child := range trie.children {
		if err := child.dfs(key+string(r), cb); err != nil {
			return err
		}
	}
	return nil
}

/*
Prefix search by string key
*/
func (trie *RuneMapTrie) PrefixSearch(key string) ([]string, error) {
	keys := []string{}
	searchCb := func(key string, value interface{}) error {
		keys = append(keys, key)
		return nil
	}
	for _, r := range key {
		if trie = trie.children[r]; trie == nil {
			return keys, nil
		}
	}
	for r, child := range trie.children {
		if err := child.dfs(key+string(r), searchCb); err != nil {
			return keys, err
		}
	}
	return keys, nil
}
func (trie *RuneMapTrie) dfsByKey(key string, cb Callback) error {
	if err := cb(key, trie.value); err != nil {
		return err
	}
	for r, child := range trie.children {
		if err := child.dfsByKey(key+string(r), cb); err != nil {
			return err
		}
	}
	return nil
}
