package config

func NewSQLCollection() *SQLCollection {
	return &SQLCollection{
		Install:    map[string][]string{},
		Insert:     map[string][]string{},
		Preupgrade: map[string]map[string][]string{},
	}
}

type SQLCollection struct {
	Install    map[string][]string            //{ project:[sql-content] }
	Insert     map[string][]string            //{ project:[sql-content] }
	Preupgrade map[string]map[string][]string //{ project:{ version:[sql-content] } }
}

func (s *SQLCollection) RegisterInstall(project, installSQL string) *SQLCollection {
	if _, ok := s.Install[project]; !ok {
		s.Install[project] = []string{installSQL}
		return s
	}
	s.Install[project] = append(s.Install[project], installSQL)
	return s
}

func (s *SQLCollection) RegisterInsert(project string, insertSQL string) *SQLCollection {
	if _, ok := s.Insert[project]; !ok {
		s.Insert[project] = []string{insertSQL}
		return s
	}
	s.Insert[project] = append(s.Insert[project], insertSQL)
	return s
}

func (s *SQLCollection) RegisterPreupgrade(project string, version, preupgradeSQL string) *SQLCollection {
	if _, ok := s.Preupgrade[project]; !ok {
		s.Preupgrade[project] = map[string][]string{
			version: {preupgradeSQL},
		}
		return s
	}
	if _, ok := s.Preupgrade[project][version]; !ok {
		s.Preupgrade[project][version] = []string{preupgradeSQL}
		return s
	}
	s.Preupgrade[project][version] = append(s.Preupgrade[project][version], preupgradeSQL)
	return s
}
