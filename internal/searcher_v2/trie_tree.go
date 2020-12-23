package searcher_v2

type trieNode struct {
	children map[rune]*trieNode
	startIndex []int
	isWordEnd bool
}

type TrieTree struct {
	root *trieNode
}

func InitTrieTree() *TrieTree {
	return &TrieTree{
		root: InitTrieNode(),
	}
}

func InitTrieNode() *trieNode {
	return &trieNode{
		children: make(map[rune]*trieNode),
	}
}

func (t *TrieTree) insert(word string, index int) {
	current := t.root
	for _, char := range word {
		if current.children[char] == nil {
			current.children[char] = InitTrieNode()
		}
		current = current.children[char]
	}

	current.isWordEnd = true
	current.startIndex = append(current.startIndex, index)
}

func (t *TrieTree) Find(word string) []int {
	current := t.root
	for _, char := range word{
		child, ok := current.children[char]
		if  !ok {
			return []int{-1}
		}
		current = child
	}
	if current.isWordEnd {
		return current.startIndex
	}
	return []int{-1}
}