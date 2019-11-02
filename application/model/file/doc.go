package file

// 共分两步：
// 一. 上传文件：上传时往file表记录无主数据
//
// 二. 提交数据：如在文章编辑表单页面点击提交按钮
//
//   提交数据时，包含以下三种情况：
//
//   1. 新增数据：设置file中的table_id/table_name/field_name/project和used_times，
//      使其成为有主记录
//
//   2. 编辑数据：
//      1) 新增图片：与“新增数据”操作相同
//      2) 删除图片：与“删除数据”操作相同
//
//   3. 删除数据：删除文件file_embedded表关联记录，设置file表used_times减1
//
