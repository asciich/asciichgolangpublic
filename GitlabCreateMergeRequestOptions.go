package asciichgolangpublic

import (
	"sort"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitlabCreateMergeRequestOptions struct {
	SourceBranchName                string
	TargetBranchName                string
	Title                           string
	Description                     string
	Labels                          []string
	SquashEnabled                   bool
	DeleteSourceBranchOnMerge       bool
	Verbose                         bool
	FailIfMergeRequestAlreadyExists bool
	AssignToSelf                    bool
}

func NewGitlabCreateMergeRequestOptions() (g *GitlabCreateMergeRequestOptions) {
	return new(GitlabCreateMergeRequestOptions)
}

func (g *GitlabCreateMergeRequestOptions) GetAssignToSelf() (assignToSelf bool) {

	return g.AssignToSelf
}

func (g *GitlabCreateMergeRequestOptions) GetDeepCopy() (copy *GitlabCreateMergeRequestOptions) {
	copy = NewGitlabCreateMergeRequestOptions()
	*copy = *g
	return copy
}

func (g *GitlabCreateMergeRequestOptions) GetDeleteSourceBranchOnMerge() (deleteSourceBranchOnMerge bool) {

	return g.DeleteSourceBranchOnMerge
}

func (g *GitlabCreateMergeRequestOptions) GetDescription() (description string, err error) {
	if g.Description == "" {
		return "", errors.TracedErrorf("Description not set")
	}

	return g.Description, nil
}

func (g *GitlabCreateMergeRequestOptions) GetDescriptionOrEmptyStringIfUnset() (description string) {
	return g.Description
}

func (g *GitlabCreateMergeRequestOptions) GetFailIfMergeRequestAlreadyExists() (failIfMergeRequestAlreadyExists bool) {

	return g.FailIfMergeRequestAlreadyExists
}

func (g *GitlabCreateMergeRequestOptions) GetLabels() (labels []string, err error) {
	if g.Labels == nil {
		return nil, errors.TracedErrorf("Labels not set")
	}

	if len(g.Labels) <= 0 {
		return nil, errors.TracedErrorf("Labels has no elements")
	}

	return g.Labels, nil
}

func (g *GitlabCreateMergeRequestOptions) GetLabelsOrEmptySliceIfUnset() (labels []string) {
	if g.Labels == nil {
		return []string{}
	}

	sort.Strings(labels)

	return labels
}

func (g *GitlabCreateMergeRequestOptions) GetSourceBranchName() (sourceBranchName string, err error) {
	if g.SourceBranchName == "" {
		return "", errors.TracedErrorf("SourceBranchName not set")
	}

	return g.SourceBranchName, nil
}

func (g *GitlabCreateMergeRequestOptions) GetSquashEnabled() (squashEnabled bool) {

	return g.SquashEnabled
}

func (g *GitlabCreateMergeRequestOptions) GetTargetBranchName() (targetBranchName string, err error) {
	if g.TargetBranchName == "" {
		return "", errors.TracedErrorf("TargetBranchName not set")
	}

	return g.TargetBranchName, nil
}

func (g *GitlabCreateMergeRequestOptions) GetTitle() (title string, err error) {
	if g.Title == "" {
		return "", errors.TracedErrorf("Title not set")
	}

	return g.Title, nil
}

func (g *GitlabCreateMergeRequestOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateMergeRequestOptions) IsTargetBranchSet() (isSet bool) {
	return g.TargetBranchName != ""
}

func (g *GitlabCreateMergeRequestOptions) MustGetDescription() (description string) {
	description, err := g.GetDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return description
}

func (g *GitlabCreateMergeRequestOptions) MustGetLabels() (labels []string) {
	labels, err := g.GetLabels()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return labels
}

func (g *GitlabCreateMergeRequestOptions) MustGetSourceBranchName() (sourceBranchName string) {
	sourceBranchName, err := g.GetSourceBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourceBranchName
}

func (g *GitlabCreateMergeRequestOptions) MustGetTargetBranchName() (targetBranchName string) {
	targetBranchName, err := g.GetTargetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetBranchName
}

func (g *GitlabCreateMergeRequestOptions) MustGetTitle() (title string) {
	title, err := g.GetTitle()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return title
}

func (g *GitlabCreateMergeRequestOptions) MustSetDescription(description string) {
	err := g.SetDescription(description)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetLabels(labels []string) {
	err := g.SetLabels(labels)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetSourceBranchName(sourceBranchName string) {
	err := g.SetSourceBranchName(sourceBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetTargetBranchName(targetBranchName string) {
	err := g.SetTargetBranchName(targetBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetTitle(title string) {
	err := g.SetTitle(title)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) SetAssignToSelf(assignToSelf bool) {
	g.AssignToSelf = assignToSelf
}

func (g *GitlabCreateMergeRequestOptions) SetDeleteSourceBranchOnMerge(deleteSourceBranchOnMerge bool) {
	g.DeleteSourceBranchOnMerge = deleteSourceBranchOnMerge
}

func (g *GitlabCreateMergeRequestOptions) SetDescription(description string) (err error) {
	if description == "" {
		return errors.TracedErrorf("description is empty string")
	}

	g.Description = description

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetFailIfMergeRequestAlreadyExists(failIfMergeRequestAlreadyExists bool) {
	g.FailIfMergeRequestAlreadyExists = failIfMergeRequestAlreadyExists
}

func (g *GitlabCreateMergeRequestOptions) SetLabels(labels []string) (err error) {
	if labels == nil {
		return errors.TracedErrorf("labels is nil")
	}

	if len(labels) <= 0 {
		return errors.TracedErrorf("labels has no elements")
	}

	g.Labels = labels

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetSourceBranchName(sourceBranchName string) (err error) {
	if sourceBranchName == "" {
		return errors.TracedErrorf("sourceBranchName is empty string")
	}

	g.SourceBranchName = sourceBranchName

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetSquashEnabled(squashEnabled bool) {
	g.SquashEnabled = squashEnabled
}

func (g *GitlabCreateMergeRequestOptions) SetTargetBranchName(targetBranchName string) (err error) {
	if targetBranchName == "" {
		return errors.TracedErrorf("targetBranchName is empty string")
	}

	g.TargetBranchName = targetBranchName

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetTitle(title string) (err error) {
	if title == "" {
		return errors.TracedErrorf("title is empty string")
	}

	g.Title = title

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
