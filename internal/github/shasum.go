package github

import (
	"context"
	"fmt"
	"io"
	"strings"
)

func GetShaSum(ctx context.Context, downloadURL string, filename string) (shaSum string, err error) {
	assetContents, assetErr := DownloadAssetContents(ctx, downloadURL)
	if assetErr != nil {
		err = fmt.Errorf("failed to download asset contents: %w", assetErr)
		return
	}

	contents, contentsErr := io.ReadAll(assetContents)
	if err != nil {
		err = fmt.Errorf("failed to read asset contents: %w", contentsErr)
		return
	}

	shaSum = findShaSum(contents, filename, shaSum)

	return
}

func findShaSum(contents []byte, filename string, shaSum string) string {
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if strings.HasSuffix(line, filename) {
			shaSum = strings.Split(line, " ")[0]
			break
		}
	}
	return shaSum
}
