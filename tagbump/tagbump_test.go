// Package tagbump provides unit tests for Git tag bumping functionality
//
// tagbump 包提供 Git 标签升级功能的单元测试
package tagbump

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestGetGitTags tests retrieving and sorting Git tags from repository
// Validates tag retrieval functionality using parent DIR as test repository
//
// TestGetGitTags 测试从仓库获取和排序 Git 标签
// 使用父目录作为测试仓库验证标签获取功能
func TestGetGitTags(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()

	tags, err := gcm.SortedGitTags()
	require.NoError(t, err)
	t.Log(neatjsons.S(tags))
}

// setupTestRepo creates a temporary git repository for testing
// Environment setup must succeed, so we use rese/must for all operations
func setupTestRepo() (string, func()) {
	// Create temp DIR - must succeed
	tempDIR := rese.V1(os.MkdirTemp("", "tagbump-test-*"))

	// Initialize git repository using osexec - must succeed
	execConfig := osexec.NewExecConfig().WithPath(tempDIR)

	// Initialize git repository - must succeed
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "config", "user.name", "Test User"))
	rese.V1(execConfig.Exec("git", "config", "user.email", "test@example.com"))

	// Create initial commit to make it a valid repo - must succeed
	testFile := filepath.Join(tempDIR, "README.md")
	must.Done(os.WriteFile(testFile, []byte("# Test Repo\n"), 0644))

	// Add and commit initial file - must succeed
	rese.V1(execConfig.Exec("git", "add", "."))
	rese.V1(execConfig.Exec("git", "commit", "-m", "Initial commit"))

	// Create initial tag - must succeed
	rese.V1(execConfig.Exec("git", "tag", "v0.0.1"))

	// Return cleanup function
	cleanup := func() {
		must.Done(os.RemoveAll(tempDIR))
	}

	return tempDIR, cleanup
}

func TestBumpTag_CurrentTagAtHead(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

	gcm := gitgo.New(tempDIR)

	config := &BumpConfig{
		TagName:     "v0.0.1",
		TagPrefix:   "v",
		VersionBase: 100,
		AutoConfirm: true,
		SkipGitPush: true,
	}

	success, err := BumpTag(gcm, config)
	require.NoError(t, err)
	require.True(t, success)

	tags := rese.C1(gcm.SortedGitTags())
	t.Log(tags)
	require.Contains(t, tags, "refs/tags/v0.0.1")
}

func TestBumpTag_VersionIncrement(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

	gcm := gitgo.New(tempDIR)
	execConfig := osexec.NewExecConfig().WithPath(tempDIR)

	require.True(t, t.Run("Create Tag", func(t *testing.T) {
		// Make a change to move HEAD forward from v0.0.1
		testFile := filepath.Join(tempDIR, "test.txt")
		must.Done(os.WriteFile(testFile, []byte("test content"), 0644))
		rese.V1(execConfig.Exec("git", "add", "."))
		rese.V1(execConfig.Exec("git", "commit", "-m", "Add test file"))

		config := &BumpConfig{
			TagName:     "v0.0.1",
			TagPrefix:   "v",
			VersionBase: 10,
			AutoConfirm: true,
			SkipGitPush: true,
		}

		success, err := BumpTag(gcm, config)
		require.NoError(t, err)
		require.True(t, success)

		tags := rese.C1(gcm.SortedGitTags())
		t.Log(tags)
		require.Contains(t, tags, "refs/tags/v0.0.2")
	}))

	require.True(t, t.Run("Update Tag", func(t *testing.T) {
		// Make a change to move HEAD forward from v0.0.1
		testFile := filepath.Join(tempDIR, "data.txt")
		must.Done(os.WriteFile(testFile, []byte("test content"), 0644))
		rese.V1(execConfig.Exec("git", "add", "."))
		rese.V1(execConfig.Exec("git", "commit", "-m", "Add test file"))

		config := &BumpConfig{
			TagName:     "v0.0.2",
			TagPrefix:   "v",
			VersionBase: 10,
			AutoConfirm: true,
			SkipGitPush: true,
		}

		success, err := BumpTag(gcm, config)
		require.NoError(t, err)
		require.True(t, success)

		tags := rese.C1(gcm.SortedGitTags())
		t.Log(tags)
		require.Contains(t, tags, "refs/tags/v0.0.3")
	}))
}

func TestBumpTag_DifferentVersionBases(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

	gcm := gitgo.New(tempDIR)
	execConfig := osexec.NewExecConfig().WithPath(tempDIR)

	t.Run("Version Base 10", func(t *testing.T) {
		// Make a change to move HEAD forward from v0.0.1
		testFile := filepath.Join(tempDIR, "test.txt")
		must.Done(os.WriteFile(testFile, []byte("test content"), 0644))
		rese.V1(execConfig.Exec("git", "add", "."))
		rese.V1(execConfig.Exec("git", "commit", "-m", "Add test file"))

		config := &BumpConfig{
			TagName:     "v0.0.1",
			TagPrefix:   "v",
			VersionBase: 10,
			AutoConfirm: true,
			SkipGitPush: true,
		}

		success, err := BumpTag(gcm, config)
		require.NoError(t, err)
		require.True(t, success)

		tags := rese.C1(gcm.SortedGitTags())
		t.Log(tags)
		require.Contains(t, tags, "refs/tags/v0.0.2")
	})

	t.Run("Version Base 100", func(t *testing.T) {
		// Make a change to move HEAD forward from v0.0.1
		testFile := filepath.Join(tempDIR, "data.txt")
		must.Done(os.WriteFile(testFile, []byte("test content"), 0644))
		rese.V1(execConfig.Exec("git", "add", "."))
		rese.V1(execConfig.Exec("git", "commit", "-m", "Add test file"))

		config := &BumpConfig{
			TagName:     "v0.0.2",
			TagPrefix:   "v",
			VersionBase: 100,
			AutoConfirm: true,
			SkipGitPush: true,
		}

		success, err := BumpTag(gcm, config)
		require.NoError(t, err)
		require.True(t, success)

		tags := rese.C1(gcm.SortedGitTags())
		t.Log(tags)
		require.Contains(t, tags, "refs/tags/v0.0.3")
	})
}
