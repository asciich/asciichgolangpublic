package asciichgolangpublic


type ReplaceBetweenMarkersOptions struct {
	WorkingDirPath string
	Verbose        bool
}

func NewReplaceBetweenMarkersOptions() (r *ReplaceBetweenMarkersOptions) {
	return new(ReplaceBetweenMarkersOptions)
}

func (r *ReplaceBetweenMarkersOptions) GetVerbose() (verbose bool, err error) {

	return r.Verbose, nil
}

func (r *ReplaceBetweenMarkersOptions) GetWorkingDirPath() (workingDirPath string, err error) {
	if r.WorkingDirPath == "" {
		return "", TracedErrorf("WorkingDirPath not set")
	}

	return r.WorkingDirPath, nil
}

func (r *ReplaceBetweenMarkersOptions) MustGetVerbose() (verbose bool) {
	verbose, err := r.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (r *ReplaceBetweenMarkersOptions) MustGetWorkingDirPath() (workingDirPath string) {
	workingDirPath, err := r.GetWorkingDirPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return workingDirPath
}

func (r *ReplaceBetweenMarkersOptions) MustSetVerbose(verbose bool) {
	err := r.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *ReplaceBetweenMarkersOptions) MustSetWorkingDirPath(workingDirPath string) {
	err := r.SetWorkingDirPath(workingDirPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *ReplaceBetweenMarkersOptions) SetVerbose(verbose bool) (err error) {
	r.Verbose = verbose

	return nil
}

func (r *ReplaceBetweenMarkersOptions) SetWorkingDirPath(workingDirPath string) (err error) {
	if workingDirPath == "" {
		return TracedErrorf("workingDirPath is empty string")
	}

	r.WorkingDirPath = workingDirPath

	return nil
}
