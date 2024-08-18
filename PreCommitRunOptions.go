package asciichgolangpublic

type PreCommitRunOptions struct {
	CommitChanges bool
	Verbose       bool
}

func NewPreCommitRunOptions() (p *PreCommitRunOptions) {
	return new(PreCommitRunOptions)
}

func (p *PreCommitRunOptions) GetCommitChanges() (commitChanges bool) {

	return p.CommitChanges
}

func (p *PreCommitRunOptions) GetVerbose() (verbose bool) {

	return p.Verbose
}

func (p *PreCommitRunOptions) SetCommitChanges(commitChanges bool) {
	p.CommitChanges = commitChanges
}

func (p *PreCommitRunOptions) SetVerbose(verbose bool) {
	p.Verbose = verbose
}
