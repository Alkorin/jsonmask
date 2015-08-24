package jsonmask

type Tree []TreeNode

type TreeNode struct {
	Field  string
	Childs Tree
}

type token struct {
	tag   rune
	value string
}

// Delimiters : ( , ) /
var delimiters = map[rune]bool{
	'(': true,
	')': true,
	'/': true,
	',': true,
}

func parseTokens(tokens []token, parent *token) ([]token, Tree) {

	tree := make(Tree, 0)

	for len(tokens) != 0 {
		t := tokens[0]
		tokens = tokens[1:]

		switch t.tag {
		case 'S':
			var childs Tree

			tokens, childs = parseTokens(tokens, &t)

			tree = append(tree, TreeNode{Field: t.value, Childs: childs})

			if parent != nil && parent.tag == '/' {
				// return if parent element was a /
				return tokens, tree
			}

		case '/':
			return parseTokens(tokens, &t)

		case '(':
			return parseTokens(tokens, &t)

		case ')':
			return tokens, tree

		case ',':
			return tokens, tree

		default:
			panic("Should not happend")

		}
	}

	return tokens, tree
}

func Parse(s string) Tree {

	tokens := make([]token, 0)
	name := make([]rune, 0)

	// First step, tokenize string
	for _, c := range s {
		if delimiters[c] == true {
			if len(name) != 0 {
				tokens = append(tokens, token{tag: 'S', value: string(name)})
				name = make([]rune, 0)
			}
			tokens = append(tokens, token{tag: c})
		} else {
			name = append(name, c)
		}
	}

	// Final name part
	if len(name) != 0 {
		tokens = append(tokens, token{tag: 'S', value: string(name)})
	}

	// Second step, parse tokens
	_, tree := parseTokens(tokens, nil)

	return tree
}
