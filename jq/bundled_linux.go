//go:generate statik -f -ns jq -include "*" -src ../bin/jq/linux_amd64/
package jq

import (
	_ "github.com/owenthereal/jqplay/jq/statik"
)

const osBin = "/jq"
