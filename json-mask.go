package jsonmask

import "fmt"

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

func _parseTokens(tokens []token, parent *token, deep int) ([]token, Tree, error) {

	tree := make(Tree, 0)

TokenLoop:
	for len(tokens) != 0 {
		var err error
		t := tokens[0]
		tokens = tokens[1:]

		switch t.tag {
		case 'S':
			var childs Tree

			tokens, childs, err = _parseTokens(tokens, &t, deep)

			if err != nil {
				return nil, nil, err
			}

			tree = append(tree, TreeNode{Field: t.value, Childs: childs})

			if parent != nil && parent.tag == '/' {
				// return if parent element was a /
				break TokenLoop
			}

		case '/':
			if parent == nil || parent.tag != 'S' {
				// Error: invalid parents
				return nil, nil, fmt.Errorf("error while parsing")
			}

			tokens, tree, err = _parseTokens(tokens, &t, deep)

			if err != nil {
				return nil, nil, err
			}
			if len(tree) == 0 {
				// Error: no childs
				return nil, nil, fmt.Errorf("error while parsing")
			}

			break TokenLoop

		case '(':
			if parent == nil || parent.tag != 'S' {
				// Error: invalid parents
				return nil, nil, fmt.Errorf("error while parsing")
			}

			tokens, tree, err = _parseTokens(tokens, &t, deep+1)

			if err != nil {
				return nil, nil, err
			}
			if len(tree) == 0 {
				// Error: invalid childs
				return nil, nil, fmt.Errorf("error while parsing")
			}

			break TokenLoop

		case ')':
			if deep == 0 {
				// Error: parentheses are not balanced
				return nil, nil, fmt.Errorf("error while parsing")
			}
			break TokenLoop

		case ',':
			if parent == nil || parent.tag == ',' || parent.tag == '/' {
				// Error: invalid parents
				return nil, nil, fmt.Errorf("error while parsing")
			}
			if len(tokens) == 0 {
				// Error: nothing more to parse
				return nil, nil, fmt.Errorf("error while parsing")
			}
			break TokenLoop

		default:
			return nil, nil, fmt.Errorf("error while parsing")

		}
	}

	return tokens, tree, nil
}

func parseTokens(tokens []token) (Tree, error) {

	_, tree, err := _parseTokens(tokens, nil, 0)
	return tree, err
}

func Parse(s string) (Tree, error) {

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
	return parseTokens(tokens)
}
