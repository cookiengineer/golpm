package lpm

import "testing"

func assertTrieSearch(t *testing.T, trie Trie, value string, want_prefix uint8, want_address string) {

	t.Helper()

	got := trie.Search(value)

	if got == nil {

		t.Errorf("Search(%q) returned nil, expected %s/%d", value, want_address, want_prefix)
		return

	}

	if got.Prefix != want_prefix {
		t.Errorf("Search(%q).Prefix = %d, expected %d", value, got.Prefix, want_prefix)
	}

	if got.Address != want_address {
		t.Errorf("Search(%q).Address = %q, expected %q", value, got.Address, want_address)
	}

}

func TestTrieInsertWrongOrderIPv4(t *testing.T) {

	trie := NewTrie()

	trie.Insert("1.2.3.0/24")
	trie.Insert("1.2.3.0/16")
	trie.Insert("1.2.3.0/8")

	root := trie.Rootv4

	if len(root.Children) != 1 {
		t.Fatalf("root has %d children, expected 1", len(root.Children))
	}

	node8 := root.Children[0]

	if node8.Prefix != 8 || node8.Address != "1.2.3.0" {
		t.Errorf("root child = %s/%d, expected 1.2.3.0/8", node8.Address, node8.Prefix)
	}

	if len(node8.Children) != 1 {
		t.Fatalf("1.2.3.0/8 has %d children, expected 1", len(node8.Children))
	}

	node16 := node8.Children[0]

	if node16.Prefix != 16 || node16.Address != "1.2.3.0" {
		t.Errorf("1.2.3.0/8 child = %s/%d, expected 1.2.3.0/16", node16.Address, node16.Prefix)
	}

	if len(node16.Children) != 1 {
		t.Fatalf("1.2.3.0/16 has %d children, expected 1", len(node16.Children))
	}

	node24 := node16.Children[0]

	if node24.Prefix != 24 || node24.Address != "1.2.3.0" {
		t.Errorf("1.2.3.0/16 child = %s/%d, expected 1.2.3.0/24", node24.Address, node24.Prefix)
	}

	assertTrieSearch(t, trie, "1.2.3.4", 24, "1.2.3.0")

}

func TestTrieInsertCorrectOrderIPv4(t *testing.T) {

	trie := NewTrie()

	trie.Insert("1.2.3.0/8")
	trie.Insert("1.2.3.0/16")
	trie.Insert("1.2.3.0/24")

	root := trie.Rootv4

	if len(root.Children) != 1 {
		t.Fatalf("root has %d children, expected 1", len(root.Children))
	}

	node8 := root.Children[0]
	node16 := node8.Children[0]
	node24 := node16.Children[0]

	if node8.Prefix != 8 || node16.Prefix != 16 || node24.Prefix != 24 {
		t.Errorf("tree nesting is wrong: %d -> %d -> %d", node8.Prefix, node16.Prefix, node24.Prefix)
	}

	assertTrieSearch(t, trie, "1.2.3.4", 24, "1.2.3.0")

}

func TestTrieInsertReversesDeepTree(t *testing.T) {

	trie := NewTrie()

	trie.Insert("1.2.3.4/32")
	trie.Insert("1.2.3.0/24")
	trie.Insert("1.2.3.0/16")
	trie.Insert("1.2.3.0/8")

	root := trie.Rootv4

	if len(root.Children) != 1 {
		t.Fatalf("root has %d children, expected 1", len(root.Children))
	}

	node8 := root.Children[0]

	if len(node8.Children) != 1 {
		t.Fatalf("1.2.3.0/8 has %d children, expected 1", len(node8.Children))
	}

	node16 := node8.Children[0]

	if len(node16.Children) != 1 {
		t.Fatalf("1.2.3.0/16 has %d children, expected 1", len(node16.Children))
	}

	node24 := node16.Children[0]

	if len(node24.Children) != 1 {
		t.Fatalf("1.2.3.0/24 has %d children, expected 1", len(node24.Children))
	}

	node32 := node24.Children[0]

	if node32.Prefix != 32 || node32.Address != "1.2.3.4" {
		t.Errorf("1.2.3.0/24 child = %s/%d, expected 1.2.3.4/32", node32.Address, node32.Prefix)
	}

	assertTrieSearch(t, trie, "1.2.3.4", 32, "1.2.3.4")

}

func TestTrieSearchLongestPrefixMatch(t *testing.T) {

	trie := NewTrie()

	trie.Insert("192.169.0.0/16")
	trie.Insert("192.168.0.0/24")
	trie.Insert("192.169.128.0/17")
	trie.Insert("192.169.0.0/24")

	assertTrieSearch(t, trie, "192.168.0.123", 24, "192.168.0.0")
	assertTrieSearch(t, trie, "192.169.0.123", 24, "192.169.0.0")
	assertTrieSearch(t, trie, "192.169.128.123", 17, "192.169.128.0")
	assertTrieSearch(t, trie, "192.169.64.123", 16, "192.169.0.0")
	assertTrieSearch(t, trie, "8.8.8.8", 0, "0.0.0.0")

}

func TestTrieInsertWrongOrderIPv6(t *testing.T) {

	trie := NewTrie()

	trie.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/64")
	trie.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/48")
	trie.Insert("2001:0db8:0000:0000:0000:0000:0000:0000/32")

	root := trie.Rootv6

	if len(root.Children) != 1 {
		t.Fatalf("root has %d children, expected 1", len(root.Children))
	}

	node32 := root.Children[0]

	if node32.Prefix != 32 {
		t.Errorf("root child prefix = %d, expected 32", node32.Prefix)
	}

	if len(node32.Children) != 1 {
		t.Fatalf("/32 node has %d children, expected 1", len(node32.Children))
	}

	node48 := node32.Children[0]

	if node48.Prefix != 48 {
		t.Errorf("/32 child prefix = %d, expected 48", node48.Prefix)
	}

	if len(node48.Children) != 1 {
		t.Fatalf("/48 node has %d children, expected 1", len(node48.Children))
	}

	node64 := node48.Children[0]

	if node64.Prefix != 64 {
		t.Errorf("/48 child prefix = %d, expected 64", node64.Prefix)
	}

}

func TestTrieSetNodesUnsorted(t *testing.T) {

	trie := NewTrie()

	values := []Node{
		ToNode("1.2.3.0/24"),
		ToNode("1.2.3.0/8"),
		ToNode("1.2.3.0/16"),
	}

	if trie.SetNodes(values) != true {
		t.Fatalf("SetNodes returned false, expected true")
	}

	root := trie.Rootv4

	if len(root.Children) != 1 {
		t.Fatalf("root has %d children, expected 1", len(root.Children))
	}

	node8 := root.Children[0]

	if node8.Prefix != 8 {
		t.Errorf("root child prefix = %d, expected 8", node8.Prefix)
	}

	if len(node8.Children) != 1 {
		t.Fatalf("/8 node has %d children, expected 1", len(node8.Children))
	}

	node16 := node8.Children[0]

	if len(node16.Children) != 1 {
		t.Fatalf("/16 node has %d children, expected 1", len(node16.Children))
	}

	node24 := node16.Children[0]

	if node24.Prefix != 24 {
		t.Errorf("/16 child prefix = %d, expected 24", node24.Prefix)
	}

	assertTrieSearch(t, trie, "1.2.3.4", 24, "1.2.3.0")

}
