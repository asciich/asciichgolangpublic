package bashutils

import _ "embed"

//go:embed files/default_script_base.sh
var DEFAULT_SCRIPT_STRUCTURE string

// Returns a default bash script as string.
// Useful as starting point for new bash scripts.
func GetDefaultScriptStructure() string {
	return DEFAULT_SCRIPT_STRUCTURE
}
