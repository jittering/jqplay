//go:generate statik -f -ns jq -include "*" -src ../bin/jq/windows_amd64/
package jq

import (
	_ "github.com/owenthereal/jqplay/jq/statik"
)

const osBin = "/jq.exe"
