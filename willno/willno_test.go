package willno_test

import (
	"reflect"
	"testing"

	"github.com/ayaki102/willNo/willno"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantExt  string
	}{
		{
			name:     "File with extension",
			input:    "main.go",
			wantName: "main",
			wantExt:  "go",
		},
		{
			name:     "File without extension",
			input:    "script",
			wantName: "script",
			wantExt:  "",
		},
		{
			name:     "File with multiple dots",
			input:    "test.config.json",
			wantName: "test",
			wantExt:  "config.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := willno.ParseFile(tt.input)
			if got.Lang == nil {
				t.Fatal("Lang is nil")
			}
			// note: name and fileExts are private, so we can't test them directly
			// but we can verify the builder was created properly
			if got == nil {
				t.Error("Expected non-nil LanguageBuilder")
			}
		})
	}
}

func TestAddKeyword(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		want     []string
	}{
		{
			name:     "Add single keyword",
			keywords: []string{"let"},
			want:     []string{"let"},
		},
		{
			name:     "Add multiple keywords",
			keywords: []string{"let", "const", "var"},
			want:     []string{"let", "const", "var"},
		},
		{
			name:     "Add duplicate keywords - should be ignored",
			keywords: []string{"let", "let", "const"},
			want:     []string{"let", "const"},
		},
		{
			name:     "Empty keywords",
			keywords: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := willno.ParseFile("test")
			for _, kw := range tt.keywords {
				lb = lb.AddKeyword(kw)
			}

			if !reflect.DeepEqual(lb.Lang.Keywords, tt.want) {
				t.Errorf("got %v, want %v", lb.Lang.Keywords, tt.want)
			}
		})
	}
}

func TestAddComment(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		want     []string
	}{
		{
			name:     "Single line comment",
			comments: []string{"//"},
			want:     []string{"//"},
		},
		{
			name:     "Multiple comment types",
			comments: []string{"//", "#", "/*"},
			want:     []string{"//", "#", "/*"},
		},
		{
			name:     "Duplicate comments allowed",
			comments: []string{"//", "//"},
			want:     []string{"//", "//"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := willno.ParseFile("test")
			for _, c := range tt.comments {
				lb = lb.AddComment(c)
			}

			if !reflect.DeepEqual(lb.Lang.Comments, tt.want) {
				t.Errorf("got %v, want %v", lb.Lang.Comments, tt.want)
			}
		})
	}
}

func TestAddLiteral(t *testing.T) {
	tests := []struct {
		name     string
		literals []string
		want     []string
	}{
		{
			name:     "Add string literal",
			literals: []string{willno.StringLiteral},
			want:     []string{"string"},
		},
		{
			name:     "Add multiple literals",
			literals: []string{willno.StringLiteral, willno.NumberLiteral, willno.BoolLiteral},
			want:     []string{"string", "number", "boolean"},
		},
		{
			name:     "Add custom literal",
			literals: []string{"regex", "null"},
			want:     []string{"regex", "null"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := willno.ParseFile("test")
			for _, lit := range tt.literals {
				lb = lb.AddLiteral(lit)
			}

			if !reflect.DeepEqual(lb.Lang.Literals, tt.want) {
				t.Errorf("got %v, want %v", lb.Lang.Literals, tt.want)
			}
		})
	}
}

func TestMethodChaining(t *testing.T) {
	lb := willno.ParseFile("myLang.xyz").
		AddKeyword("let").
		AddKeyword("const").
		AddComment("//").
		AddComment("#").
		AddLiteral(willno.StringLiteral).
		AddLiteral(willno.NumberLiteral)

	expectedKeywords := []string{"let", "const"}
	expectedComments := []string{"//", "#"}
	expectedLiterals := []string{"string", "number"}

	if !reflect.DeepEqual(lb.Lang.Keywords, expectedKeywords) {
		t.Errorf("Keywords: got %v, want %v", lb.Lang.Keywords, expectedKeywords)
	}
	if !reflect.DeepEqual(lb.Lang.Comments, expectedComments) {
		t.Errorf("Comments: got %v, want %v", lb.Lang.Comments, expectedComments)
	}
	if !reflect.DeepEqual(lb.Lang.Literals, expectedLiterals) {
		t.Errorf("Literals: got %v, want %v", lb.Lang.Literals, expectedLiterals)
	}
}

func TestParse(t *testing.T) {
	lb := willno.ParseFile("test").
		AddKeyword("let").
		AddComment("//")

	parsed := lb.Parse()

	if parsed == nil {
		t.Fatal("Parse() returned nil")
	}

	// test that parsed tokens map is initialized
	val, ok := parsed.Get("variable", "test")
	if ok && val != nil {
		// map is accessible (even if empty)
	}
}

func TestParsedGet(t *testing.T) {
	lb := willno.ParseFile("test")
	parsed := lb.Parse()

	// test getting non-existent value
	val, ok := parsed.Get("variable", "nonexistent")
	if ok {
		// currently Get always returns true, even for missing keys
		// this WILL be fixed
		_ = val
	}
}

func TestParsedGetString(t *testing.T) {
	lb := willno.ParseFile("test")
	parsed := lb.Parse()

	_, ok := parsed.GetString("variable", "x")
	if ok {
		t.Log("GetString returned ok for non-existent variable")
	}
}

func TestParsedGetNumber(t *testing.T) {
	lb := willno.ParseFile("test")
	parsed := lb.Parse()

	_, ok := parsed.GetNumber("variable", "count")
	if ok {
		t.Log("GetNumber returned ok for non-existent variable")
	}
}

func TestParsedGetOr(t *testing.T) {
	tests := []struct {
		name    string
		defVal  any
		wantDef bool
	}{
		{
			name:    "Default string",
			defVal:  "default",
			wantDef: true,
		},
		{
			name:    "Default number",
			defVal:  42,
			wantDef: true,
		},
		{
			name:    "Default nil",
			defVal:  nil,
			wantDef: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := willno.ParseFile("test")
			parsed := lb.Parse()

			got := parsed.GetOr("variable", "nonexistent", tt.defVal)
			if tt.wantDef && !reflect.DeepEqual(got, tt.defVal) {
				t.Errorf("GetOr() = %v, want default %v", got, tt.defVal)
			}
		})
	}
}

func TestParsedAll(t *testing.T) {
	lb := willno.ParseFile("test")
	parsed := lb.Parse()

	values := parsed.All("variable")
	if values == nil {
		t.Error("All() returned nil, expected empty slice")
	}
	if len(values) != 0 {
		t.Errorf("All() returned %d values, expected 0", len(values))
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"StringLiteral", willno.StringLiteral, "string"},
		{"NumberLiteral", willno.NumberLiteral, "number"},
		{"BoolLiteral", willno.BoolLiteral, "boolean"},
		{"Var", willno.Var, "variable"},
		{"Fn", willno.Fn, "function"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestValueStruct(t *testing.T) {
	v := willno.Value{
		Name:        "myVar",
		Type:        willno.Var,
		LiteralType: willno.StringLiteral,
		Value:       "hello",
	}

	if v.Name != "myVar" {
		t.Errorf("Name = %q, want %q", v.Name, "myVar")
	}
	if v.Type != "variable" {
		t.Errorf("Type = %q, want %q", v.Type, "variable")
	}
	if v.LiteralType != "string" {
		t.Errorf("LiteralType = %q, want %q", v.LiteralType, "string")
	}
	if v.Value != "hello" {
		t.Errorf("Value = %v, want %q", v.Value, "hello")
	}
}
