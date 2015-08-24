package jsonmask

import "testing"
import "reflect"

var patterns = map[string]Tree{
	"a": Tree{
		TreeNode{
			Field:  "a",
			Childs: Tree{},
		},
	},
	"a,b": Tree{
		TreeNode{
			Field:  "a",
			Childs: Tree{},
		},
		TreeNode{
			Field:  "b",
			Childs: Tree{},
		},
	},
	"a/b": Tree{
		TreeNode{
			Field: "a",
			Childs: Tree{
				TreeNode{
					Field:  "b",
					Childs: Tree{},
				},
			},
		},
	},
	"a(b,c)": Tree{
		TreeNode{
			Field: "a",
			Childs: Tree{
				TreeNode{
					Field:  "b",
					Childs: Tree{},
				},
				TreeNode{
					Field:  "c",
					Childs: Tree{},
				},
			},
		},
	},
	"a,b/c,d(e,f/g)": Tree{
		TreeNode{
			Field:  "a",
			Childs: Tree{},
		},
		TreeNode{
			Field: "b",
			Childs: Tree{
				TreeNode{
					Field:  "c",
					Childs: Tree{},
				},
			},
		},
		TreeNode{
			Field: "d",
			Childs: Tree{
				TreeNode{
					Field:  "e",
					Childs: Tree{},
				},
				TreeNode{
					Field: "f",
					Childs: Tree{
						TreeNode{
							Field:  "g",
							Childs: Tree{},
						},
					},
				},
			},
		},
	},
}

func TestPatterns(t *testing.T) {

	for pattern, wanted := range patterns {

		tree := Parse(pattern)

		if reflect.DeepEqual(tree, wanted) == false {
			t.Errorf(`Error while parsing %q, got %+v instead of %+v`, pattern, tree, wanted)
		}

	}
}

func TestInvalidTokens(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Should have paniced with invalid token`)
		}
	}()

	// This should die
	parseTokens([]token{token{tag: 'I'}}, nil)
}
