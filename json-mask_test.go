package jsonmask

import "testing"
import "reflect"

var patterns = map[string]Tree{
	"a": []TreeNode{
		TreeNode{
			Field:  "a",
			Childs: []TreeNode{},
		},
	},
	"a,b": []TreeNode{
		TreeNode{
			Field:  "a",
			Childs: []TreeNode{},
		},
		TreeNode{
			Field:  "b",
			Childs: []TreeNode{},
		},
	},
	"a/b": []TreeNode{
		TreeNode{
			Field: "a",
			Childs: []TreeNode{
				TreeNode{
					Field:  "b",
					Childs: []TreeNode{},
				},
			},
		},
	},
	"a(b,c)": []TreeNode{
		TreeNode{
			Field: "a",
			Childs: []TreeNode{
				TreeNode{
					Field:  "b",
					Childs: []TreeNode{},
				},
				TreeNode{
					Field:  "c",
					Childs: []TreeNode{},
				},
			},
		},
	},
	"a,b/c,d(e,f/g)": []TreeNode{
		TreeNode{
			Field:  "a",
			Childs: []TreeNode{},
		},
		TreeNode{
			Field: "b",
			Childs: []TreeNode{
				TreeNode{
					Field:  "c",
					Childs: []TreeNode{},
				},
			},
		},
		TreeNode{
			Field: "d",
			Childs: []TreeNode{
				TreeNode{
					Field:  "e",
					Childs: []TreeNode{},
				},
				TreeNode{
					Field: "f",
					Childs: []TreeNode{
						TreeNode{
							Field:  "g",
							Childs: []TreeNode{},
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
