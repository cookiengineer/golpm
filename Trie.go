package lpm

import "github.com/cookiengineer/lpm/utils"
import "fmt"
import "sort"
import "strings"

func findTrieNode(node *Node, value Node) *Node {

	if node.Type == value.Type && node.Address == value.Address && node.Prefix == value.Prefix {
		return node
	}

	var result *Node = nil

	if node.ContainsNode(value) {

		if len(node.Children) > 0 {

			for c := 0; c < len(node.Children); c++ {

				check := findTrieNode(node.Children[c], value)

				if check != nil {

					if result != nil {

						if check.Prefix > result.Prefix {
							result = check
						}

					} else if result == nil {
						result = check
					}

				}

			}

			if result == nil {
				result = node
			}

		} else {
			result = node
		}

	}

	return result

}

func printTrieNode(node *Node, indent string) {

	fmt.Println(indent + "└─" + node.String())

	if len(node.Children) > 0 {

		for c := 0; c < len(node.Children); c++ {
			printTrieNode(node.Children[c], indent + "  ")
		}

	}

}

type Trie struct {
	Rootv4 *Node `json:"rootv4"`
	Rootv6 *Node `json:"rootv6"`
}

func NewTrie() Trie {

	var trie Trie

	var rootv4 = NewNode("0.0.0.0", 0)
	var rootv6 = NewNode("0000:0000:0000:0000:0000:0000:0000:0000", 0)

	rootv4.Name = "AS0"
	rootv6.Name = "AS0"

	trie.Rootv4 = &rootv4
	trie.Rootv6 = &rootv6

	return trie

}

func (trie *Trie) Insert(value string) bool {

	node := ToNode(value)

	return trie.InsertNode(node)

}

func insertTrieNode(root *Node, node Node) bool {

	var result bool = false

	parent := findTrieNode(root, node)

	if parent != nil {

		if node.Prefix > parent.Prefix {

			children := make([]*Node, 0)
			kept := make([]*Node, 0)

			for c := 0; c < len(parent.Children); c++ {

				child := parent.Children[c]

				if node.ContainsNode(*child) {
					children = append(children, child)
				} else {
					kept = append(kept, child)
				}

			}

			node.Children = children
			parent.Children = kept
			parent.Children = append(parent.Children, &node)

			result = true

		}

	}

	return result

}

func (trie *Trie) InsertNode(value Node) bool {

	var result bool = false
	var node = value

	if node.Type == "ipv6" {
		result = insertTrieNode(trie.Rootv6, node)
	} else if node.Type == "ipv4" {
		result = insertTrieNode(trie.Rootv4, node)
	}

	return result

}

func (trie *Trie) Print() {

	fmt.Println("ipv4")
	printTrieNode(trie.Rootv4, "")

	fmt.Println("")

	fmt.Println("ipv6")
	printTrieNode(trie.Rootv6, "")

}

func (trie *Trie) Search(value string) *Node {

	var result *Node = nil

	if strings.Contains(value, "/") {

		address := value[0:strings.Index(value, "/")]
		prefix := value[strings.Index(value, "/")+1:]

		if utils.IsIPv6(address) {

			address = utils.ToIPv6(address)
			result = trie.SearchNode(ToNode(address + "/" + prefix))

		} else if utils.IsIPv4(address) {

			address = utils.ToIPv4(address)
			result = trie.SearchNode(ToNode(address + "/" + prefix))

		}

	} else {

		if utils.IsIPv6(value) {

			value = utils.ToIPv6(value)
			result = trie.SearchNode(ToNode(value + "/128"))

		} else if utils.IsIPv4(value) {

			value = utils.ToIPv4(value)
			result = trie.SearchNode(ToNode(value + "/32"))

		}

	}

	return result

}

func (trie *Trie) SearchNode(value Node) *Node {

	var result *Node = nil

	if value.Type == "ipv6" {

		node := findTrieNode(trie.Rootv6, value)

		if node != nil {
			result = node
		}

	} else if value.Type == "ipv4" {

		node := findTrieNode(trie.Rootv4, value)

		if node != nil {
			result = node
		}

	}

	return result

}

func (trie *Trie) SetNodes(values []Node) bool {

	var result bool = true

	sort.Slice(values, func(a int, b int) bool {

		if values[a].Prefix == values[b].Prefix {

			if values[a].Type == values[b].Type {

				if values[a].Type == "ipv4" {

					bytes_a := utils.ToIPv4Bytes(values[a].Address, values[a].Prefix)
					bytes_b := utils.ToIPv4Bytes(values[b].Address, values[b].Prefix)

					for b := 0; b < len(bytes_a); b++ {

						if bytes_a[b] == bytes_b[b] {
							continue
						} else {
							return bytes_a[b] < bytes_b[b]
						}

					}

					return false

				} else if values[a].Type == "ipv6" {

					bytes_a := utils.ToIPv6Bytes(values[a].Address, values[a].Prefix)
					bytes_b := utils.ToIPv6Bytes(values[b].Address, values[b].Prefix)

					for b := 0; b < len(bytes_a); b++ {

						if bytes_a[b] == bytes_b[b] {
							continue
						} else {
							return bytes_a[b] < bytes_b[b]
						}

					}

					return false

				} else {
					return false
				}

			} else {
				return values[a].Type < values[b].Type
			}

		} else {
			return values[a].Prefix < values[b].Prefix
		}

	})

	for v := 0; v < len(values); v++ {

		check := trie.InsertNode(values[v])

		if check == false {
			result = false
		}

	}

	return result

}
