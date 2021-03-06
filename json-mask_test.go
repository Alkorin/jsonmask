package jsonmask

import "flag"
import "math/rand"
import "testing"
import "reflect"

var fuzzCount = flag.Uint("test.fuzz", 0, "define fuzzing loop count")

var patterns = map[string]Tree{
	"": Tree{},
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
	// Some invalid patterns
	",":         nil,
	"(":         nil,
	")":         nil,
	"/":         nil,
	"a,":        nil,
	"a,,":       nil,
	"a))":       nil,
	"a(":        nil,
	"a/":        nil,
	"a/(":       nil,
	"a(b":       nil,
	"a((":       nil,
	"a,,b":      nil,
	"a(b,)":     nil,
	"a(b,,)":    nil,
	"a(b,,c)":   nil,
	"(a,b)":     nil,
	"a(b,c//d)": nil,
	"a(b)c":     nil,
}

func TestPatterns(t *testing.T) {

	for pattern, wanted := range patterns {

		tree, err := Parse(pattern)

		if wanted == nil && (tree != nil || err == nil) {
			t.Errorf(`Error while parsing %q, got tree=%+v, error=%+v instead of nil, error`, pattern, tree, err)
			continue
		}

		if reflect.DeepEqual(tree, wanted) == false {
			t.Errorf(`Error while parsing %q, got %+v instead of %+v`, pattern, tree, wanted)
			continue
		}
	}
}

func TestOneInvalidTokens(t *testing.T) {

	// This should return an error
	tree, err := parseTokens([]token{token{tag: 'I'}})

	if tree != nil || err == nil {
		t.Errorf(`Should have returned error with invalid token`)
	}

}

func TestDeepInvalidTokens(t *testing.T) {

	// This should return an error
	tree, err := parseTokens([]token{token{tag: 'S', value: "foo"}, token{tag: 'I'}})

	if tree != nil || err == nil {
		t.Errorf(`Should have returned error with invalid token`)
	}

}

func TestFuzzing(t *testing.T) {

	count := *fuzzCount

	if count == 0 {
		t.Skip("Fuzzing test skipped")
	}

	// Valid tokens
	var runes = []rune("a()/,")

	for count > 0 {
		// Random size [4,16)
		size := 4 + rand.Int31n(12)

		s := make([]rune, size)
		for i := range s {
			s[i] = runes[rand.Intn(len(runes))]
		}

		tree, err := Parse(string(s))

		if err == nil {
			t.Errorf("Found a random valid pattern %q, got %+v", string(s), tree)
		}

		count--
	}

}
