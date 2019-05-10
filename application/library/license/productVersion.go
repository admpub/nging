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

package license

type ProductVersion struct {
	Id               uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	ProductId        uint64 `db:"product_id" bson:"product_id" comment:"产品ID" json:"product_id" xml:"product_id"`
	Version          string `db:"version" bson:"version" comment:"版本号(格式1.0.1)" json:"version" xml:"version"`
	Type             string `db:"type" bson:"type" comment:"版本类型(stable-稳定版;beta-公测版;alpha-内测版)" json:"type" xml:"type"`
	Os               string `db:"os" bson:"os" comment:"支持的操作系统(多个用逗号分隔)，留空表示不限制" json:"os" xml:"os"`
	Arch             string `db:"arch" bson:"arch" comment:"硬件架构(多个用逗号分隔)，留空表示不限制" json:"arch" xml:"arch"`
	ReleasedAt       uint   `db:"released_at" bson:"released_at" comment:"发布时间" json:"released_at" xml:"released_at"`
	Created          uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated          uint   `db:"updated" bson:"updated" comment:"修改时间" json:"updated" xml:"updated"`
	Disabled         string `db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
	Audited          string `db:"audited" bson:"audited" comment:"是否已审核" json:"audited" xml:"audited"`
	ForceUpgrade     string `db:"force_upgrade" bson:"force_upgrade" comment:"是否强行升级为此版本" json:"force_upgrade" xml:"force_upgrade"`
	Description      string `db:"description" bson:"description" comment:"发布说明" json:"description" xml:"description"`
	Remark           string `db:"remark" bson:"remark" comment:"备注" json:"remark" xml:"remark"`
	DownloadUrl      string `db:"download_url" bson:"download_url" comment:"下载网址" json:"download_url" xml:"download_url"`
	Sign             string `db:"sign" bson:"sign" comment:"下载后验证签名(多个签名之间用逗号分隔)" json:"sign" xml:"sign"`
	DownloadUrlOther string `db:"download_url_other" bson:"download_url_other" comment:"备用下载网址" json:"download_url_other" xml:"download_url_other"`
}
