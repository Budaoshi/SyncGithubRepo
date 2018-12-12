package service

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"strings"
	"os/exec"
	"github.com/juju/errors"
)

/**
 *
 * @author BDS
 * @version 0.0.1
 */

type API struct {
	HttpWork
}

type CompareContent []ContentResponse

type CMD struct {
	Dir        string
	LocalRepo  string
	RemoteRepo string
	Token      string
	Path 	   string
}

//GET /repos/:owner/:repo/contents/:path
func (api *API) CompareRepos(cmd *CMD) (err error) {

	local := api.MakeRequest(api.GetClient(), makeRequest(true, "GET", fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", cmd.LocalRepo, cmd.Path), nil))

	if local.StatusCode != 200 {
		return errors.New(fmt.Sprintf("request %s error: %v", cmd.LocalRepo, string(local.Body)))
	}

	remote := api.MakeRequest(api.GetClient(), makeRequest(false, "GET", fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", cmd.RemoteRepo, cmd.Path), nil))

	if remote.StatusCode != 200 {
		return errors.New(fmt.Sprintf("request %s error: %v", cmd.RemoteRepo, string(remote.Body)))
	}

	localContents := CompareContent{}
	remoteContents := CompareContent{}

	json.Unmarshal(local.Body, &localContents)
	json.Unmarshal(remote.Body, &remoteContents)

	contentMap := make(map[string]ContentResponse)

	for i, ele := range localContents {
		contentMap[ele.Name] = localContents[i]
	}

	var contentAddMap []ContentResponse
	var contentDeleteMap []ContentResponse
	contentUpdateMap := make(map[string]ContentResponse)

	for i, ele := range remoteContents {

		value, ok := contentMap[ele.Name]

		if !ok {
			//新增文件
			contentAddMap = append(contentAddMap, remoteContents[i])
			continue
		}

		if value.Sha != ele.Sha {
			//修改文件
			contentUpdateMap[value.Sha] = remoteContents[i]
			continue
		}

	}
	//寻找删除的文件
	for i, ele := range localContents {
		found := false
		for _, item := range remoteContents {
			if ele.Name == item.Name {
				found = true
				break
			}
		}
		if !found {
			//文件被删除
			contentDeleteMap = append(contentDeleteMap, localContents[i])
		}
	}

	//生成输出文档

	i := len(contentAddMap)
	j := len(contentUpdateMap)
	k := len(contentDeleteMap)

	if i == 0 && j == 0 && k == 0 {
		fmt.Println("All files are consistent!!!")
		return
	}

	strBuilder := strings.Builder{}
	strBuilder.WriteString("## 概述\n")
	strBuilder.WriteString("此 Issue 由脚本自动创建 \n")
	strBuilder.WriteString(fmt.Sprintf("**本次扫描结果：新增 %d，修改 %d，删除 %d** \n", i, j, k))

	strBuilder.WriteString("### 文件变化列表 \n")
	strBuilder.WriteString("| 变更类型 | 文件名 | \n")
	strBuilder.WriteString("| :----: | :---- | \n")

	if i > 0 {
		for _, ele := range contentAddMap {
			strBuilder.WriteString(fmt.Sprintf("| %s | [%s](%s) | \n", "ADDED", ele.Name, ele.Html_url))
		}
	}

	if j > 0 {
		for _, ele := range contentUpdateMap {
			strBuilder.WriteString(fmt.Sprintf("| %s | [%s](%s) | \n", "UPDATED", ele.Name, ele.Html_url))
		}
	}

	if k > 0 {
		for _, ele := range contentDeleteMap {
			strBuilder.WriteString(fmt.Sprintf("| %s | [%s](%s) | \n", "DELETED", ele.Name, ele.Html_url))
		}
	}

	if j > 0 {

		strBuilder.WriteString("### 文档具体变更如下(仅展示 UPDATED 内容) \n")

		for k, v := range contentUpdateMap {
			strBuilder.WriteString(fmt.Sprintf("- UPDATED [%s](%s) \n", v.Name, v.Html_url))
			strBuilder.WriteString("```diff \n")

			cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s ; git diff %s %s | cat | sed '1,2d'", cmd.Dir, k, v.Sha))
			output, err := cmd.Output()
			if err != nil {
				return err
			}
			strBuilder.WriteString(string(output))

			strBuilder.WriteString("``` \n")
		}
	}

	api.CreateIssue(cmd.LocalRepo, cmd.Token, strBuilder.String())

	fmt.Println(strBuilder.String())

	return nil

}

//POST /repos/:owner/:repo/issues
func (api *API) CreateIssue(repo string, token string, content string) {

	issueRequest := &IssueRequest{}
	issueRequest.Title = "翻译文档内容不一致"
	issueRequest.Body = content
	issueRequest.Labels = append(issueRequest.Labels, "docs")
	issueRequest.Labels = append(issueRequest.Labels, "inconsistency")
	issueRequest.Labels = append(issueRequest.Labels, "sync robot")

	req, _ := json.Marshal(issueRequest)

	request, _ := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/issues", repo), strings.NewReader(string(req)))
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	result := api.MakeRequest(api.GetClient(), request)

	if result.StatusCode == 201 {
		fmt.Println("create issue success")
		fmt.Println("")
	} else {
		fmt.Printf("create issue[%d], [%v] \n\n", result.StatusCode, string(result.Body))
	}

}

func makeRequest(self bool, method string, url string, body io.Reader) *http.Request {

	var request *http.Request

	if self {
		request, _ = http.NewRequest(method, url, body)
	} else {
		request, _ = http.NewRequest(method, url, body)
	}
	request.Header.Add("Accept", "application/vnd.github.v3+json")

	return request

}

type ContentResponse struct {
	Name     string `json:"name"`
	Sha      string `json:"sha"`
	Content  string `json:"content"`
	Html_url string `json:"html_url"`
}

type IssueRequest struct {
	Title  string   `json:"title"`  //Required. The title of the issue.
	Body   string   `json:"body"`   //The contents of the issue.
	Labels []string `json:"labels"` //Labels to associate with this issue
}
