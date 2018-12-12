#### GitHub 本地仓库与远程库同步工具

一个简单的仓库同步工具，不支持多目录结构，Golang 语言实现。

**实现原理**
参照 `GitHub`提供的 [Get contents](https://developer.github.com/v3/repos/contents/#get-contents) 和 [Create an issue](https://developer.github.com/v3/issues/#create-an-issue) 接口，拉去两个仓库的内容，并对各个文件的的 `SHA`做对比。若两个仓库文件有增加、删除、修改等，会在本地仓库创建相应的 `issue`。

**使用方法**
- 为`sys_repo`设置运行权限：
```shell
#: chmod +x sys_repo
```
- 执行运行查看使用方法
```shell
#: ./sys_repo

  -d string
        [--dir -d] workspace[~/your/local/path] synchronized already with remote's[与远程同步后的工作区]
  -dir string
        [--dir -d] workspace[~/your/local/path] synchronized already with remote's[与远程同步后的工作区]
  -l string
        [--local -l] repo[:owner/:repo] need be synchronized [需要同步的本地仓库]
  -local string
        [--local -l] repo[:owner/:repo] need be synchronized [需要同步的本地仓库]
  -p string
        [--path -p] file path or directory that will be synchronized[需要同步的文件路径或目录]
  -path string
        [--path -p] file path or directory that will be synchronized[需要同步的文件路径或目录]
  -r string
        [--remote -r] repo[:owner/:repo] will be synchronized[将要同步的远程仓库]
  -remote string
        [--remote -r] repo[:owner/:repo] will be synchronized[将要同步的远程仓库]
  -t string
        [--token -t] token[--local specified repo] requested by creating issue[创建 issue 需要的 token
  -token string
        [--token -t] token[--local specified repo] requested by creating issue[创建 issue 需要的 token]

```
- 使用示例
```shell
./sys_repo -b ~/repo/fresco -l Budaoshi/fresco-docs-cn -r facebook/fresoc -p docs/_docs -t eW91ciB0b2tlbiB0byBjcmVhdGUgaXNzdWVzCg==
```


