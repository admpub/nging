package listener

// UpdaterInfos       [Project][Table][Field]
var UpdaterInfos = map[string]map[string]map[string]UpdaterInfo{
	``: {
		`nging_file`: {
			`view_url`: UpdaterInfo{},
		},
		`nging_file_thumb`: {
			`view_url`: UpdaterInfo{},
		},
	},
}

type UpdaterInfo struct {
	Seperator string
	Embedded bool
}

// RecordUpdaterInfo 记录
func RecordUpdaterInfo(project, table, field, seperator string, embedded bool, sameFields ...string) {
	if _, ok := UpdaterInfos[project]; !ok {
		UpdaterInfos[project] = map[string]map[string]UpdaterInfo{}
	}
	if _, ok := UpdaterInfos[project][table]; !ok {
		UpdaterInfos[project][table] = map[string]UpdaterInfo{}
	}
	if _, ok := UpdaterInfos[project][table][field]; !ok {
		UpdaterInfos[project][table][field] = UpdaterInfo{
			Seperator: seperator,
			Embedded: embedded,
		}
	}
	for _, field := range sameFields {
		if _, ok := UpdaterInfos[project][table][field]; !ok {
			UpdaterInfos[project][table][field] = UpdaterInfo{
				Seperator: seperator,
				Embedded: embedded,
			}
		}
	}
}
