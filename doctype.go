package xgen

import (
	"fmt"
	"regexp"
	"strings"
)

type entityMap map[string]string

var replacements = map[string](string){
	"&amp;":  "&",
	"&#x09;": "\t",
	"&#x20;": " ",
	"&#x22;": `"`,
	"&#x25;": `%`,
	"&#39;":  "'",
	"&#94;":  `^`, // This is 94 dec (5E hex)
	"&#96;":  "`", // This is 96dec (69 hex)
	"&#124;": "|", // This is 124 dec, should be 7c hex
}

var entityPattern = regexp.MustCompile(`<!ENTITY\s+([a-zA-Z0-9_-]+)\s"(.+)"\s*>`)

// parseDocType parses DOCTYPE for ENTITY and extracts them into a map.
func parseDoctype(txt string) (entityMap, bool) {
	if !strings.HasPrefix(txt, "DOCTYPE") {
		return nil, false
	}
	m := make(map[string]string)
	matches := entityPattern.FindAllStringSubmatch(txt, -1)
	for _, match := range matches {
		m[match[1]] = match[2]
	}
	fmt.Printf("initial entities: %d\n", len(m))
	for {
		changes := false
		for k1, v1 := range m {
			esc := "&" + k1 + ";"
			for k2, v2 := range m {
				if strings.Contains(v2, esc) {
					m[k2] = strings.ReplaceAll(v2, esc, v1)
					changes = true
				}
			}
		}
		if !changes {
			break
		}
	}
	for kr, vr := range replacements {
		for k, v := range m {
			if strings.Contains(v, kr) {
				m[k] = strings.ReplaceAll(v, kr, vr)
			}
		}
	}

	/*for k, v := range m {
		fmt.Printf("%q: %q\n", k, v)
	}*/

	return m, true
}
