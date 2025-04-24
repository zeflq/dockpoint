package build_savepoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeflq/dockpoint/src/domain"
)

// ✅ Spy mocks

type spyWriter struct {
	Called bool
	Lines  []string
}

func (s *spyWriter) Write(lines []string, savepoint string) (string, error) {
	s.Called = true
	s.Lines = lines
	return "/tmp/fake.Dockerfile", nil
}

type spyBuilder struct {
	WasCalled bool
}

func (b *spyBuilder) Build(ctx context.Context, dockerfilePath, contextDir, tag string) error {
	b.WasCalled = true
	return nil
}

type spyPusher struct {
	WasCalled bool
}

func (p *spyPusher) Push(tag string) error {
	p.WasCalled = true
	return nil
}

type alwaysFalseChecker struct{}

func (c *alwaysFalseChecker) TagExists(tag string) (bool, error) {
	return false, nil
}

type parserReturningOne struct{}

func (p *parserReturningOne) Parse(path string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{
		{Name: "base", StartLine: 0, EndLine: 1},
	}, nil
}

type simpleSlicer struct{}

func (s *simpleSlicer) Slice([]string, domain.Savepoint) ([]string, error) {
	return []string{"FROM alpine"}, nil
}

func TestExecBuildPushPath(t *testing.T) {
	tmp := t.TempDir()
	dockerfile := filepath.Join(tmp, "Dockerfile")
	_ = os.WriteFile(dockerfile, []byte("# savepoint: base\nFROM alpine"), 0644)
	_ = os.Chdir(tmp)

	writer := &spyWriter{}
	builder := &spyBuilder{}
	pusher := &spyPusher{}

	uc := NewBuildSavepointUseCase(
		&parserReturningOne{},
		&simpleSlicer{},
		writer,
		builder,
		&alwaysFalseChecker{},
		pusher,
		&mockConfig{},
	)

	result, err := uc.Execute(context.Background(), BuildSavepointRequest{
		FilePath:  dockerfile,
		Savepoint: "base",
		Force:     true,
		Push:      true,
		DryRun:    false,
	})

	assert.NoError(t, err)
	assert.False(t, result.Skipped)
	assert.True(t, writer.Called)
	assert.True(t, builder.WasCalled)
	assert.True(t, pusher.WasCalled)
}