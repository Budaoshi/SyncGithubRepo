package main

import (
	"flag"
	"os"
	"service"
	"fmt"
)

/**
 *
 * @author BDS
 * @version 0.0.1
 */

func main() {


	cmd := &service.CMD{}

	flag.StringVar(&cmd.Dir, "dir", "", "[--dir -d] workspace[~/your/local/path] synchronized already with remote's[与远程同步后的工作区]")
	flag.StringVar(&cmd.Dir, "d", "", "[--dir -d] workspace[~/your/local/path] synchronized already with remote's[与远程同步后的工作区]")
	flag.StringVar(&cmd.LocalRepo, "local", "", "[--local -l] repo[:owner/:repo] need be synchronized [需要同步的本地仓库]")
	flag.StringVar(&cmd.LocalRepo, "l", "", "[--local -l] repo[:owner/:repo] need be synchronized [需要同步的本地仓库]")
	flag.StringVar(&cmd.RemoteRepo, "remote", "", "[--remote -r] repo[:owner/:repo] will be synchronized[将要同步的远程仓库]")
	flag.StringVar(&cmd.RemoteRepo, "r", "", "[--remote -r] repo[:owner/:repo] will be synchronized[将要同步的远程仓库]")
	flag.StringVar(&cmd.Path, "path", "", "[--path -p] file path or directory that will be synchronized[需要同步的文件路径或目录]")
	flag.StringVar(&cmd.Path, "p", "", "[--path -p] file path or directory that will be synchronized[需要同步的文件路径或目录]")
	flag.StringVar(&cmd.Token, "token", "", "[--token -t] token[--local specified repo] requested by creating issue[创建 issue 需要的 token]")
	flag.StringVar(&cmd.Token, "t", "", "[--token -t] token[--local specified repo] requested by creating issue[创建 issue 需要的 token")


	flag.Parse()

	if cmd.Dir == "" || cmd.LocalRepo == "" || cmd.RemoteRepo == ""  || cmd.Token == ""{
		flag.Usage()
		os.Exit(-1)
	}

	var hubApi = &service.API{}
	err := hubApi.CompareRepos(cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}



