package webdav

import "net/url"

func ParseUserForm(v url.Values) []*WebdavUser {
	indexes, _ := v[`webdav_user_index`]
	users, _ := v[`webdav_user`]
	passwords, _ := v[`webdav_pass`]
	roots, _ := v[`webdav_user_root`]
	userWriteables, _ := v[`webdav_user_writeable`]
	var list []*WebdavUser
	for key, index := range indexes {
		if key >= len(users) || key >= len(passwords) || key >= len(roots) || key >= len(userWriteables) {
			continue
		}
		u := &WebdavUser{
			User:     users[key],
			Password: passwords[key],
			Root:     roots[key],
		}
		if len(u.User) == 0 {
			continue
		}
		u.SetWriteable(userWriteables[key])
		readables, _ := v[`webdav_readables[user][`+index+`]`]
		writeables, _ := v[`webdav_writeables[user][`+index+`]`]
		resources, _ := v[`webdav_resources[user][`+index+`]`]
		for pkey, resource := range resources {
			if pkey >= len(readables) || pkey >= len(writeables) || len(resource) == 0 {
				continue
			}
			p := &WebdavPerm{
				Resource: resource,
			}
			p.SetReadable(readables[pkey])
			p.SetWriteable(writeables[pkey])
			u.Perms = append(u.Perms, p)
		}
		list = append(list, u)
	}
	return list
}

func ParseGlobalForm(v url.Values) []*WebdavPerm {
	var list []*WebdavPerm
	readables, _ := v[`webdav_readables[global]`]
	writeables, _ := v[`webdav_writeables[global]`]
	resources, _ := v[`webdav_resources[global]`]
	for pkey, resource := range resources {
		if pkey >= len(readables) || pkey >= len(writeables) || len(resource) == 0 {
			continue
		}
		p := &WebdavPerm{
			Resource: resource,
		}
		p.SetReadable(readables[pkey])
		p.SetWriteable(writeables[pkey])
		list = append(list, p)
	}
	return list
}
