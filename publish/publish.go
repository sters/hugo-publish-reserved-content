package publish

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/morikuni/failure"
)

const (
	targetFile    = ".md"
	hugoSeparator = "---\n"
)

var (
	// ErrNotTarget on specified filepath
	ErrNotTarget = failure.StringCode("not target file")
	// ErrFileCannotLoad on specified filepath
	ErrFileCannotLoad = failure.StringCode("file cannot load")
	// ErrFileEmpty on specified filepath
	ErrFileEmpty = failure.StringCode("file is empty")
	// ErrFileContentMismatch on specified filepath
	ErrFileContentMismatch = failure.StringCode("file content mismatch")
	// ErrContentIsReservedButNotDraft on specified filepath
	ErrContentIsReservedButNotDraft = failure.StringCode("content is reserved but not draft")
	// ErrContentIsNotReserved on specified filepath
	ErrContentIsNotReserved = failure.StringCode("content is not reserved")
	// ErrContentIsNotTheTimeYet on specified filepath
	ErrContentIsNotTheTimeYet = failure.StringCode("content is not the time yet")

	readFile  = ioutil.ReadFile
	writeFile = ioutil.WriteFile
)

// New is constructor of Publisher
func New(reservedKey string, draftKey string) *Publisher {
	return &Publisher{
		reservedKey: reservedKey,
		draftKey:    draftKey,
	}
}

// Publisher doing publish reserved content
type Publisher struct {
	reservedKey string
	draftKey    string
}

// CheckReservedAndPublish reserved content
func (p *Publisher) CheckReservedAndPublish(filepath string) error {
	if !strings.Contains(filepath, targetFile) {
		return failure.New(ErrNotTarget)
	}

	rawContent, err := readFile(filepath)
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileCannotLoad))
	}

	if len(rawContent) == 0 {
		return failure.New(ErrFileEmpty)
	}

	content := strings.Split(string(rawContent), hugoSeparator)
	if len(content) < 3 {
		return failure.New(ErrFileContentMismatch)
	}

	var v map[string]interface{}
	if err := yaml.Unmarshal([]byte(content[1]), &v); err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	if _, ok := v[p.reservedKey]; !ok {
		return failure.New(ErrContentIsNotReserved)
	}
	if d, ok := v[p.draftKey]; !ok || d != true {
		return failure.New(ErrContentIsReservedButNotDraft)
	}

	t, err := time.Parse(time.RFC3339, v["date"].(string))
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	now := time.Now()
	if t.UnixNano() > now.UnixNano() {
		return failure.New(ErrContentIsNotTheTimeYet)
	}

	delete(v, p.reservedKey)
	delete(v, p.draftKey)

	meta, err := yaml.Marshal(v)
	if err != nil {
		log.Printf("%s: error: %+v", filepath, err)
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	content[1] = string(meta)
	if err := writeFile(filepath, []byte(strings.Join(content, hugoSeparator)), os.ModePerm); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
