/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package fileupdater

func (f *FileUpdater) SetProject(project string) *FileUpdater {
	f.project = project
	return f
}

func (f *FileUpdater) Set(table string, field string, tableID string) *FileUpdater {
	f.table = table
	f.field = field
	f.tableID = tableID
	return f
}

func (f *FileUpdater) SetTable(table string) *FileUpdater {
	f.table = table
	return f
}

func (f *FileUpdater) SetField(field string) *FileUpdater {
	f.field = field
	return f
}

func (f *FileUpdater) SetTableID(tableID string) *FileUpdater {
	f.tableID = tableID
	return f
}
