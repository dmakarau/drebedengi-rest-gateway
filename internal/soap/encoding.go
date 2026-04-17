package soap

import (
	"fmt"
	"strings"
)

// encodeParam dispatches to the correct encoder based on the Go value type.
func encodeParam(p Param) string {
	if p.Value == nil {
		return fmt.Sprintf(`<%s xsi:nil="true"/>`, p.Name)
	}

	switch v := p.Value.(type) {
	case string:
		return encodeSimple(p.Name, "xsd:string", xmlEscape(v))
	case int:
		return encodeSimple(p.Name, "xsd:int", fmt.Sprintf("%d", v))
	case int64:
		return encodeSimple(p.Name, "xsd:int", fmt.Sprintf("%d", v))
	case bool:
		val := "false"
		if v {
			val = "true"
		}
		return encodeSimple(p.Name, "xsd:boolean", val)
	case float64:
		return encodeSimple(p.Name, "xsd:float", fmt.Sprintf("%g", v))
	case map[string]any:
		return encodeMap(p.Name, v)
	case []map[string]any:
		return encodeArrayOfMaps(p.Name, v)
	case []int64:
		return encodeArrayOfInts(p.Name, v)
	case []string:
		return encodeArrayOfStrings(p.Name, v)
	default:
		return encodeSimple(p.Name, "xsd:string", fmt.Sprintf("%v", v))
	}
}

// encodeSimple emits a single element with xsi:type.
func encodeSimple(name, xsiType, value string) string {
	return fmt.Sprintf(`<%s xsi:type="%s">%s</%s>`, name, xsiType, value, name)
}

// encodeMap emits an Apache XML-SOAP Map (ns2:Map).
func encodeMap(name string, m map[string]any) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<%s xsi:type="ns2:Map">`, name))
	for k, v := range m {
		b.WriteString(`<item>`)
		b.WriteString(fmt.Sprintf(`<key xsi:type="xsd:string">%s</key>`, xmlEscape(k)))
		b.WriteString(encodeMapValue(v))
		b.WriteString(`</item>`)
	}
	b.WriteString(fmt.Sprintf(`</%s>`, name))
	return b.String()
}

// encodeMapValue encodes a value inside a Map item.
func encodeMapValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf(`<value xsi:type="xsd:string">%s</value>`, xmlEscape(val))
	case int:
		return fmt.Sprintf(`<value xsi:type="xsd:int">%d</value>`, val)
	case int64:
		return fmt.Sprintf(`<value xsi:type="xsd:int">%d</value>`, val)
	case bool:
		s := "false"
		if val {
			s = "true"
		}
		return fmt.Sprintf(`<value xsi:type="xsd:boolean">%s</value>`, s)
	case float64:
		return fmt.Sprintf(`<value xsi:type="xsd:float">%g</value>`, val)
	case []string:
		return encodeArrayOfStringsAs("value", val)
	case []int64:
		return encodeArrayOfIntsAs("value", val)
	case nil:
		return `<value xsi:nil="true"/>`
	default:
		return fmt.Sprintf(`<value xsi:type="xsd:string">%v</value>`, val)
	}
}

// encodeArrayOfMaps emits a SOAP-ENC:Array containing ns2:Map items.
func encodeArrayOfMaps(name string, items []map[string]any) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<%s xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="ns2:Map[%d]">`, name, len(items)))
	for _, m := range items {
		b.WriteString(`<item xsi:type="ns2:Map">`)
		for k, v := range m {
			b.WriteString(`<item>`)
			b.WriteString(fmt.Sprintf(`<key xsi:type="xsd:string">%s</key>`, xmlEscape(k)))
			b.WriteString(encodeMapValue(v))
			b.WriteString(`</item>`)
		}
		b.WriteString(`</item>`)
	}
	b.WriteString(fmt.Sprintf(`</%s>`, name))
	return b.String()
}

// encodeArrayOfInts emits a SOAP-ENC:Array of xsd:int.
func encodeArrayOfInts(name string, items []int64) string {
	return encodeArrayOfIntsAs(name, items)
}

func encodeArrayOfIntsAs(name string, items []int64) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<%s xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="xsd:int[%d]">`, name, len(items)))
	for _, v := range items {
		b.WriteString(fmt.Sprintf(`<item xsi:type="xsd:int">%d</item>`, v))
	}
	b.WriteString(fmt.Sprintf(`</%s>`, name))
	return b.String()
}

// encodeArrayOfStrings emits a SOAP-ENC:Array of xsd:string.
func encodeArrayOfStrings(name string, items []string) string {
	return encodeArrayOfStringsAs(name, items)
}

func encodeArrayOfStringsAs(name string, items []string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<%s xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="xsd:string[%d]">`, name, len(items)))
	for _, v := range items {
		b.WriteString(fmt.Sprintf(`<item xsi:type="xsd:string">%s</item>`, xmlEscape(v)))
	}
	b.WriteString(fmt.Sprintf(`</%s>`, name))
	return b.String()
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}
