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

func parseTokens(tokens []token) (Tree, error) {

	// Parenthesis counter
	deep := 0

	var _parseTokens func(parent *token) (Tree, error)
	_parseTokens = func(parent *token) (Tree, error) {

		tree := make(Tree, 0)

		for len(tokens) > 0 {
			var err error
			t := tokens[0]
			tokens = tokens[1:]

			switch t.tag {
			case 'S':
				var childs Tree

				childs, err = _parseTokens(&t)

				if err != nil {
					return nil, err
				}

				tree = append(tree, TreeNode{Field: t.value, Childs: childs})

				if parent != nil && parent.tag == '/' {
					// return if parent element was a /
					return tree, nil
				}

			case '/':
				if parent == nil || parent.tag != 'S' {
					// Error: invalid parents
					return nil, fmt.Errorf("error while parsing")
				}

				tree, err = _parseTokens(&t)

				if err != nil {
					return nil, err
				}
				if len(tree) == 0 {
					// Error: no childs
					return nil, fmt.Errorf("error while parsing")
				}

				return tree, nil

			case '(':
				if parent == nil || parent.tag != 'S' {
					// Error: invalid parents
					return nil, fmt.Errorf("error while parsing")
				}

				deep++
				tree, err = _parseTokens(&t)

				if err != nil {
					return nil, err
				}
				if len(tree) == 0 {
					// Error: invalid childs
					return nil, fmt.Errorf("error while parsing")
				}

				return tree, nil

			case ')':
				if parent == nil || parent.tag == ',' || parent.tag == '/' || parent.tag == '(' {
					// Error: invalid parents
					return nil, fmt.Errorf("error while parsing")
				}
				if deep == 0 {
					// Error: parentheses are not balanced
					return nil, fmt.Errorf("error while parsing")
				}
				if len(tokens) > 0 && tokens[0].tag == 'S' {
					// Error: no childs
					return nil, fmt.Errorf("error while parsing")

				}
				deep--
				return tree, nil

			case ',':
				if parent == nil || parent.tag == ',' || parent.tag == '/' || parent.tag == '(' {
					// Error: invalid parents
					return nil, fmt.Errorf("error while parsing")
				}
				if len(tokens) == 0 {
					// Error: nothing more to parse
					return nil, fmt.Errorf("error while parsing")
				}
				return tree, nil

			default:
				return nil, fmt.Errorf("error while parsing")

			}
		}

		return tree, nil
	}

	tree, err := _parseTokens(nil)

	if deep != 0 {
		// Error: parentheses are not balanced
		return nil, fmt.Errorf("error while parsing")
	}

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
