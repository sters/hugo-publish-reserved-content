package publish

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/morikuni/failure"
	hugocontent "github.com/sters/simple-hugo-content-parse"
)

const targetFile = ".md"

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

	content, err := hugocontent.ParseMarkdownWithYaml(bytes.NewBuffer(rawContent))
	if err != nil {
		return failure.New(ErrFileContentMismatch)
	}

	if _, ok := content.FrontMatter[p.reservedKey]; !ok {
		return failure.New(ErrContentIsNotReserved)
	}
	if d, ok := content.FrontMatter[p.draftKey]; !ok || d != true {
		return failure.New(ErrContentIsReservedButNotDraft)
	}

	t, err := time.Parse(time.RFC3339, content.FrontMatter["date"].(string))
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	now := time.Now()
	if t.UnixNano() > now.UnixNano() {
		return failure.New(ErrContentIsNotTheTimeYet)
	}

	delete(content.FrontMatter, p.reservedKey)
	delete(content.FrontMatter, p.draftKey)

	result, err := content.Dump()
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	if err := writeFile(filepath, []byte(result), os.ModePerm); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
