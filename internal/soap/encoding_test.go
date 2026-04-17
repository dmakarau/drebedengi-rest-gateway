package soap

import (
	"strings"
	"testing"
)

func TestEncodeParam_String(t *testing.T) {
	got := encodeParam(Param{Name: "apiId", Value: "demo_api"})
	want := `<apiId xsi:type="xsd:string">demo_api</apiId>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_Int(t *testing.T) {
	got := encodeParam(Param{Name: "count", Value: 42})
	want := `<count xsi:type="xsd:int">42</count>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_Int64(t *testing.T) {
	got := encodeParam(Param{Name: "id", Value: int64(123456)})
	want := `<id xsi:type="xsd:int">123456</id>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_BoolTrue(t *testing.T) {
	got := encodeParam(Param{Name: "flag", Value: true})
	want := `<flag xsi:type="xsd:boolean">true</flag>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_BoolFalse(t *testing.T) {
	got := encodeParam(Param{Name: "flag", Value: false})
	want := `<flag xsi:type="xsd:boolean">false</flag>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_Nil(t *testing.T) {
	got := encodeParam(Param{Name: "idList", Value: nil})
	want := `<idList xsi:nil="true"/>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_Map(t *testing.T) {
	// Use a single-entry map to avoid iteration order issues.
	got := encodeParam(Param{Name: "params", Value: map[string]any{
		"r_period": 8,
	}})
	want := `<params xsi:type="ns2:Map">` +
		`<item><key xsi:type="xsd:string">r_period</key>` +
		`<value xsi:type="xsd:int">8</value></item>` +
		`</params>`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestEncodeParam_MapMultipleKeys(t *testing.T) {
	m := map[string]any{
		"is_report": false,
		"r_period":  8,
	}
	got := encodeParam(Param{Name: "params", Value: m})

	// Can't check exact string due to map iteration order.
	// Check structure instead.
	if !strings.HasPrefix(got, `<params xsi:type="ns2:Map">`) {
		t.Errorf("missing Map prefix: %s", got)
	}
	if !strings.HasSuffix(got, `</params>`) {
		t.Errorf("missing closing tag: %s", got)
	}
	if !strings.Contains(got, `<key xsi:type="xsd:string">is_report</key><value xsi:type="xsd:boolean">false</value>`) {
		t.Errorf("missing is_report entry: %s", got)
	}
	if !strings.Contains(got, `<key xsi:type="xsd:string">r_period</key><value xsi:type="xsd:int">8</value>`) {
		t.Errorf("missing r_period entry: %s", got)
	}
}

func TestEncodeParam_ArrayOfInts(t *testing.T) {
	got := encodeParam(Param{Name: "idList", Value: []int64{10, 20, 30}})
	want := `<idList xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="xsd:int[3]">` +
		`<item xsi:type="xsd:int">10</item>` +
		`<item xsi:type="xsd:int">20</item>` +
		`<item xsi:type="xsd:int">30</item>` +
		`</idList>`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestEncodeParam_ArrayOfStrings(t *testing.T) {
	got := encodeParam(Param{Name: "ids", Value: []string{"a", "b"}})
	want := `<ids xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="xsd:string[2]">` +
		`<item xsi:type="xsd:string">a</item>` +
		`<item xsi:type="xsd:string">b</item>` +
		`</ids>`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestEncodeParam_ArrayOfMaps(t *testing.T) {
	got := encodeParam(Param{Name: "list", Value: []map[string]any{
		{"name": "test"},
	}})
	want := `<list xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="ns2:Map[1]">` +
		`<item xsi:type="ns2:Map">` +
		`<item><key xsi:type="xsd:string">name</key>` +
		`<value xsi:type="xsd:string">test</value></item>` +
		`</item>` +
		`</list>`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestXmlEscape(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"hello", "hello"},
		{"a & b", "a &amp; b"},
		{"<tag>", "&lt;tag&gt;"},
		{`say "hi"`, "say &quot;hi&quot;"},
		{"it's", "it&apos;s"},
		{"a & <b> \"c\" 'd'", "a &amp; &lt;b&gt; &quot;c&quot; &apos;d&apos;"},
	}
	for _, tt := range tests {
		got := xmlEscape(tt.in)
		if got != tt.want {
			t.Errorf("xmlEscape(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEncodeParam_Float64(t *testing.T) {
	got := encodeParam(Param{Name: "amount", Value: float64(3.14)})
	want := `<amount xsi:type="xsd:float">3.14</amount>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeParam_UnknownType(t *testing.T) {
	// Unknown types fall through to xsd:string via fmt.Sprintf.
	got := encodeParam(Param{Name: "x", Value: struct{ V int }{V: 7}})
	if !strings.Contains(got, `xsi:type="xsd:string"`) {
		t.Errorf("expected xsd:string fallback, got: %s", got)
	}
}

func TestEncodeMapValue_AllTypes(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"string", "hello", `<value xsi:type="xsd:string">hello</value>`},
		{"int", 5, `<value xsi:type="xsd:int">5</value>`},
		{"int64", int64(99), `<value xsi:type="xsd:int">99</value>`},
		{"bool_true", true, `<value xsi:type="xsd:boolean">true</value>`},
		{"bool_false", false, `<value xsi:type="xsd:boolean">false</value>`},
		{"float64", 3.14, `<value xsi:type="xsd:float">3.14</value>`},
		{"nil", nil, `<value xsi:nil="true"/>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeMapValue(tt.val)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
