package fstabconfigfile

type Entry struct {
	Device  string
	Dir     string
	Type    string
	Options string
	Dump    string
	Fsck    string
}
