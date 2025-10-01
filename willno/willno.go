package willno

const (
	StringLiteral = "string"
	NumberLiteral = "number"
	BoolLiteral   = "boolean"
)

// ---------------------------
// Core Value representation
// ---------------------------
type Value struct {
	Name        string // identifier name (variable, function, etc.)
	Type        string // token type: "variable", "function", "literal", "comment", etc.
	LiteralType string // "string", "number", "boolean", or custom literal type
	Value       any    // actual parsed value
}

// ---------------------------
// ParsedFile: end-user access
// ---------------------------
type ParsedFile struct {
	tokens map[string]map[string]Value // type -> name -> Value
}

// Example QoL access methods (stub implementations that compile)
func (pf *ParsedFile) Get(tokenType, name string) (any, bool) {
	return nil, false
}

func (pf *ParsedFile) GetString(tokenType, name string) (string, bool) {
	val, ok := pf.Get(tokenType, name)
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

func (pf *ParsedFile) GetNumber(tokenType, name string) (float64, bool) {
	val, ok := pf.Get(tokenType, name)
	if !ok {
		return 0, false
	}
	n, ok := val.(float64)
	return n, ok
}

func (pf *ParsedFile) GetOr(tokenType, name string, def any) any {
	if val, ok := pf.Get(tokenType, name); ok {
		return val
	}
	return def
}

func (pf *ParsedFile) All(tokenType string) []Value {
	return nil
}

func (pf *ParsedFile) Filter(tokenType string, predicate func(Value) bool) []Value {
	return nil
}

// ---------------------------
// Language: internal representation of a user-defined language
// ---------------------------
type Language struct {
	Name     string            // language identifier
	FileExts []string          // optional: file extensions
	Keywords map[string]string // logical type -> actual keyword ("variable" -> "let")
	Comments []string          // comment patterns, e.g., "//", "#"
	Literals []string          // supported literal types: "string", "number", "boolean", or custom
}

func NewLanguage() *LanguageBuilder {
	return &LanguageBuilder{
		lang: &Language{},
	}
}

// ---------------------------
// Token: internal representation of each parsed element
// ---------------------------
type Token struct {
	Type        string // "keyword", "literal", "comment", "unknown"
	Name        string // identifier name
	LiteralType string // only set for literals
	Value       string // raw content from file
}

// ---------------------------
// Parser: entry point for UX
// ---------------------------
type Parser struct {
	languages map[string]*Language
}

// ---------------------------
// LanguageBuilder: for chaining API
// ---------------------------
type LanguageBuilder struct {
	lang *Language
}

// ---------------------------
// Optional: Variable handle for fast ergonomic access
// ---------------------------
type VarHandle[T any] struct {
	value T
}

func (v *VarHandle[T]) Set(val T) { v.value = val }
func (v *VarHandle[T]) Get() T    { return v.value }

// token clasification <- kinda nuts
// token := nextToken(file)
// if parser.IsLiteral(token.Value) {
//     token.Type = Literal
//     token.LiteralType = detectLiteralType(token.Value) // "string", "number", etc.
// } else if parser.IsKeyword(token.Value) {
//     token.Type = Keyword
// } else if parser.IsComment(token.Value) {
//     token.Type = Comment
// } else {
//     // Unknown token
//     token.Type = Unknown
// }
//
// for getting values parsed by lang implement Get(type, value)
// 	QUALITY OF LIFE FUNCTIONS
// parsed.GetString("variable", "x")   // returns string
// parsed.GetNumber("variable", "y")   // returns float64/int
// func (pf *ParsedFile) GetOr[T any](tokenType, name string, def T) T {
//     if val, ok := pf.Get[T](tokenType, name); ok {
//         return val
//     }
//     return def
// }
// age := parsed.GetOr[int]("variable", "age", 18)
