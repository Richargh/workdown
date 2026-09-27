package markdown

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/richargh/workdown/internal/kernel"
)

func Render(item kernel.WorkItem) []byte {
	project := fieldValue(item, "project")
	issueType := fieldValue(item, "issueType")
	var b strings.Builder
	b.WriteString("+++\n")
	writeTOMLString(&b, "remote", item.Remote)
	writeTOMLString(&b, "provider", item.Provider)
	writeTOMLString(&b, "key", item.Key)
	writeTOMLString(&b, "id", item.ID)
	writeTOMLString(&b, "project", project)
	writeTOMLString(&b, "issue_type", issueType)
	writeTOMLString(&b, "title_hash", titleHash(item.Title))
	b.WriteString("+++\n\n")
	b.WriteString("# ")
	b.WriteString(item.Title)
	b.WriteString(" {#wd-field-title}\n")
	return []byte(b.String())
}

func writeTOMLString(b *strings.Builder, key string, value string) {
	_, _ = fmt.Fprintf(b, "%s = %s\n", key, strconv.Quote(value))
}

func titleHash(title string) string {
	sum := sha256.Sum256([]byte(title))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func fieldValue(item kernel.WorkItem, name string) string {
	for _, field := range item.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}
