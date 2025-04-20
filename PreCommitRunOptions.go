package asciichgolangpublic

type PreCommitRunOptions struct {
	CommitChanges bool
}

func NewPreCommitRunOptions() (p *PreCommitRunOptions) {
	return new(PreCommitRunOptions)
}

func (p *PreCommitRunOptions) GetCommitChanges() (commitChanges bool) {

	return p.CommitChanges
}

func (p *PreCommitRunOptions) SetCommitChanges(commitChanges bool) {
	p.CommitChanges = commitChanges
}
