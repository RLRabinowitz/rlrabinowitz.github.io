package github

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetTags(repositoryUrl string) ([]string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--tags", repositoryUrl)
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, nil
	}

	tags := make([]string, 0)
	for _, line := range strings.Split(buf.String(), "\n") {
		if !strings.Contains(line, "refs/tags/") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("module tags in wrong format")
		}
		ref := fields[1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			continue
		}
		tag := strings.TrimPrefix(ref, "refs/tags/")
		if strings.Contains(tag, "^") {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
