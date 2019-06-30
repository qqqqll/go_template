package main

import (
	"os"
	//"strings"
	"text/template"
	"fmt"
	"log"
	"github.com/360EntSecGroup-Skylar/excelize"
)


const (
	template_str = `package define

	const (
		{{range $i, $main := .Mainds}}
		{{$main.Code}} = {{$main.Command_num}} // {{$main.Command_explain}}
		{{end}}
	)

	{{range $main,$subs:=.SubIds}}
	// {{$main.Command_explain}}的子命令
	const (
		{{range $sub_i,$sub:=$subs}}
		{{$sub.Code}} = {{$sub.Command_num}} // {{$sub.Command_explain}}
		{{end}}
	)
	{{end}}
	`
)

// 用于模板输出的结构体
type CommandsInfo struct {
	// 表示当前excel的sheet内容
	Mainds []Command_xlsx

	SubIds map[Command_xlsx][]Command_xlsx
}

// 表sheet内容的结构体
type Command_xlsx struct {
	Command_num       string //命令编号
	Code              string //命令代码
	Command_explain   string //命令注释

}

func main() {

	// 声明模板输出使用的结构体
	var CmdsInfo CommandsInfo
	CmdsInfo = GetCommandsSlice(CmdsInfo)

	//解析输出模板
	tmpl, err := template.New("test").Parse(template_str)
	if err != nil {
		panic(err)
	}
	// 表数据通过模板输出
	// err = tmpl.Execute(os.Stdout, CmdsInfo)
	// if err != nil {
	// 	panic(err)
	// }

	file, error := os.OpenFile("./error_define.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if error != nil {
		panic(error)
	}
	tmpl.Execute(file, CmdsInfo)
	file.Close()

}

func GetCommandsSlice(commandsInfo CommandsInfo) CommandsInfo {

	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	f, err := excelize.OpenFile("./errorDefine.xlsx")
	if err != nil {
		panic(err)

	}

	SheetMap := f.GetSheetMap()
	sheetName := SheetMap[1]

	//fmt.Printf("01-%v",sheetName)

	// 用于sheet表内容的切片
	command_slice := []Command_xlsx{}

	// 处理单个sheet的数据
	rows := f.GetRows(sheetName)

	var cx Command_xlsx
	fmt.Println("add ",rows[0],rows[1])
	for row_index, row := range rows {

		if row_index > 0 {
			if len(row)<3{
				log.Println("define errorDefine read fail for row !")
			}
			cx.Command_num = row[0]
			cx.Code = row[1]
			cx.Command_explain = row[2]

			command_slice = append(command_slice, cx)
		}
	}
	// 将表sheet内信息组建到输出模板的map内

	commandsInfo.Mainds = command_slice

	return commandsInfo
}
