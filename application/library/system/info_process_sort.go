package system

// - cpu

type ProcessListSortByCPUPercent []*Process

func (s ProcessListSortByCPUPercent) Len() int { return len(s) }
func (s ProcessListSortByCPUPercent) Less(i, j int) bool {
	return s[i].CPUPercent < s[j].CPUPercent
}
func (s ProcessListSortByCPUPercent) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByCPUPercentReverse []*Process

func (s ProcessListSortByCPUPercentReverse) Len() int { return len(s) }
func (s ProcessListSortByCPUPercentReverse) Less(i, j int) bool {
	return s[i].CPUPercent > s[j].CPUPercent
}
func (s ProcessListSortByCPUPercentReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// - mem

type ProcessListSortByMemPercent []*Process

func (s ProcessListSortByMemPercent) Len() int { return len(s) }
func (s ProcessListSortByMemPercent) Less(i, j int) bool {
	return s[i].MemPercent < s[j].MemPercent
}
func (s ProcessListSortByMemPercent) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByMemPercentReverse []*Process

func (s ProcessListSortByMemPercentReverse) Len() int { return len(s) }
func (s ProcessListSortByMemPercentReverse) Less(i, j int) bool {
	return s[i].MemPercent > s[j].MemPercent
}
func (s ProcessListSortByMemPercentReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// - thread

type ProcessListSortByNumThreads []*Process

func (s ProcessListSortByNumThreads) Len() int { return len(s) }
func (s ProcessListSortByNumThreads) Less(i, j int) bool {
	return s[i].NumThreads < s[j].NumThreads
}
func (s ProcessListSortByNumThreads) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByNumThreadsReverse []*Process

func (s ProcessListSortByNumThreadsReverse) Len() int { return len(s) }
func (s ProcessListSortByNumThreadsReverse) Less(i, j int) bool {
	return s[i].NumThreads > s[j].NumThreads
}
func (s ProcessListSortByNumThreadsReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// - fd

type ProcessListSortByNumFDs []*Process

func (s ProcessListSortByNumFDs) Len() int { return len(s) }
func (s ProcessListSortByNumFDs) Less(i, j int) bool {
	return s[i].NumFDs < s[j].NumFDs
}
func (s ProcessListSortByNumFDs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByNumFDsReverse []*Process

func (s ProcessListSortByNumFDsReverse) Len() int { return len(s) }
func (s ProcessListSortByNumFDsReverse) Less(i, j int) bool {
	return s[i].NumFDs > s[j].NumFDs
}
func (s ProcessListSortByNumFDsReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// - pid

type ProcessListSortByPid []*Process

func (s ProcessListSortByPid) Len() int { return len(s) }
func (s ProcessListSortByPid) Less(i, j int) bool {
	return s[i].Pid < s[j].Pid
}
func (s ProcessListSortByPid) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByPidReverse []*Process

func (s ProcessListSortByPidReverse) Len() int { return len(s) }
func (s ProcessListSortByPidReverse) Less(i, j int) bool {
	return s[i].Pid > s[j].Pid
}
func (s ProcessListSortByPidReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// - created

type ProcessListSortByCreated []*Process

func (s ProcessListSortByCreated) Len() int { return len(s) }
func (s ProcessListSortByCreated) Less(i, j int) bool {
	return s[i].created < s[j].created
}
func (s ProcessListSortByCreated) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ProcessListSortByCreatedReverse []*Process

func (s ProcessListSortByCreatedReverse) Len() int { return len(s) }
func (s ProcessListSortByCreatedReverse) Less(i, j int) bool {
	return s[i].created > s[j].created
}
func (s ProcessListSortByCreatedReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
