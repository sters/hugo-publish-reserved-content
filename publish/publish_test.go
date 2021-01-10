package publish

import (
	"os"
	"testing"

	"github.com/morikuni/failure"
)

func TestPublisher_CheckReservedAndPublish(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantCode failure.StringCode
	}{
		{"empty", "", ErrFileContentMismatch},
		{"missing reserved", `
---
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
draft: true
---
Happy New Year!
---
`, ErrFileContentMismatch},
		{"missing draft", `
---
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
---
Happy New Year!
---
`, ErrFileContentMismatch},
		{"not draft", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: false
---
Happy New Year!
`, ErrContentIsReservedButNotDraft},
		{"is not the time", `
---
date: 2100-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: true
---
Happy New Year!
`, ErrContentIsNotTheTimeYet},
		{"success", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: true
---
Happy New Year!
`, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// fake
			readFile = func(filename string) ([]byte, error) {
				return []byte(tt.content), nil
			}
			writeFile = func(filename string, data []byte, perm os.FileMode) error {
				return nil
			}

			err := New("reserved", "draft").CheckReservedAndPublish("dummy.md")
			if tt.wantCode == "" {
				if err != nil {
					t.Errorf("want no error, got=%+v", err)
				}
				return
			}
			if !failure.Is(err, tt.wantCode) {
				t.Errorf("want error=%+v, got=%+v", tt.wantCode, err)
			}
		})
	}
}
