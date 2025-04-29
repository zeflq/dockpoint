package build_savepoint

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeflq/dockpoint/src/domain"
)

// ✅ Basic mocks
type mockWriter struct{}
func (m *mockWriter) Write(lines []string, savepoint string) (string, error) {
	return "/tmp/fake.Dockerfile", nil
}

type mockBuilder struct{}
func (m *mockBuilder) Build(ctx context.Context, dockerfilePath, contextDir, tag string) error {
	return nil
}

type mockChecker struct{}
func (m *mockChecker) TagExists(tag string) (bool, error) {
	return false, nil
}

type mockPusher struct{}
func (m *mockPusher) Push(tag string) error {
	return nil
}

type mockValidator struct{}
func (m *mockValidator) Validate(savepoints []domain.Savepoint) error {
	return nil
}

type mockHasher struct{}
func (m *mockHasher) Hash(content string) (string, error) {
	return "mock-hash", nil
}

type mockTagBuilder struct{}
func (m *mockTagBuilder) BuildFinalTag(ctx domain.TagContext) (string, error) {
	return fmt.Sprintf("%s:%s", ctx.Repo, ctx.FinalImageTag), nil
}

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

// ✅ Parser variants
type mockParser struct {
	dockerfile []string
}

func (m *mockParser) Parse(path string) ([]domain.Savepoint, error) {
	if len(m.dockerfile) == 0 {
		return nil, errors.New("empty dockerfile")
	}
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

type emptyParser struct{}
func (p *emptyParser) Parse(string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{}, nil
}

type parserReturningOne struct{}
func (p *parserReturningOne) Parse(path string) ([]domain.Savepoint, error) {
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

type parserSpy struct {
	Captured string
}

func (p *parserSpy) Parse(path string) ([]domain.Savepoint, error) {
	p.Captured = path
	return []domain.Savepoint{{Name: "base", StartLine: 0, EndLine: 1}}, nil
}

type parserWith struct {
	savepoints []domain.Savepoint
}

func (m *parserWith) Parse(path string) ([]domain.Savepoint, error) {
	return m.savepoints, nil
}

// ✅ Checker variants
type alwaysTagExists struct{}
func (m *alwaysTagExists) TagExists(tag string) (bool, error) { 
	return true, nil 
}

type neverTagExists struct{}
func (m *neverTagExists) TagExists(tag string) (bool, error) { 
	return false, nil 
}

// ✅ Slicer variants
type mockSlicer struct {
	content []string
}

func (m *mockSlicer) Slice(lines []string, sp domain.Savepoint) ([]string, error) {
	return []string{"FROM node:20"}, nil
}

type simpleSlicer struct{}
func (s *simpleSlicer) Slice([]string, domain.Savepoint) ([]string, error) {
	return []string{"FROM alpine"}, nil
}

type dummySlicer struct{}
func (s *dummySlicer) Slice(lines []string, sp domain.Savepoint) ([]string, error) {
	return lines, nil
}