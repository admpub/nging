# Done

v1.0

1. 支持各种形式的HTTP请求
2. 支持代理模式：socks5， http(s)
3. 支持Cookie持久
4. 支持模拟上传下载文件

v1.1 20180402

1. 解决中国区用户被墙问题。

# Doing List



# Todo List




--------------

代理提交规范:

1. 不能在master和dev开发
2. 只能从dev中拉出代码分支, 拉出分支开发后, 等待dev合并, 然后dev再合并到master上, 并打tag
3. 如果是小分支,可以合并到大分支, 然后删除
4. Tag命名: v1.0.0 前面表示新版本号, 后面表示新特征号,最后面表示bug修复号
5. 大分支命名: release-1.0 前面表示新版本号, 后面表示新特征号

由于我们的项目较小, 4和5暂不实行, 只在master打标签, 标签为方便不取修bug号,即v1.0.从dev拉出的分支命名任意,
合并入dev后, 根据需求可以决定是否删除.

```
git tag 查看tag
git push origin --tags 表示推tag
git tag -d v1.0.0  表示删除本地tag
git push origin tag --delete v1.0.0表示删除远程tag

git branch -a 查看分支
git branch -d rm 删除本地分支
git push origin --delete rm 删除远程分支
	fmt.Printf("%#v", []byte(""))
git diff release-2.0 release-1.0 分支对比
```