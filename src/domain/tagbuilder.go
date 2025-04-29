package domain

type TagContext struct {
	FileContent string
	Repo        string
	SavepointName string
	FinalImageTag string
	IsLast      bool
}

type TagBuilder interface {
	BuildFinalTag(ctx TagContext) (string, error)
}
