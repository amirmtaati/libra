package test

import (
	"github.com/amirmtaati/libra/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestScanner(t *testing.T) {
	tmpDir := t.TempDir()
	testFiles := []core.File{
		core.NewFile("document.pdf", filepath.Join(tmpDir, "document.pdf"), ".pdf"),
		core.NewFile("config.json", filepath.Join(tmpDir, "config.json"), ".json"),
	}

	for _, file := range testFiles {
		err := os.WriteFile(file.Path, []byte("content"), 0644)
		require.NoError(t, err)
	}

	scanner := core.NewScanner(tmpDir)
	files, err := scanner.ScanDir()
	require.NoError(t, err)
	expectedResult := testFiles

	assert.Equal(t, len(files), len(expectedResult))
	sort.Slice(testFiles, func(i, j int) bool {
		return testFiles[i].Title < testFiles[j].Title
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Title < files[j].Title
	})

	for i, expected := range expectedResult {
		actual := files[i]
		assert.Equal(t, actual.Title, expected.Title)
	}
}
