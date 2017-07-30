package main

// 程序需求
// 1.win gui实现界面
// 2.配置文件进行配置输入输出文件夹
// 3.输入文件夹进行时间分离(避免重合)
// 4.并发复制文件

import (
	"fmt"
	"time"
	//	"fmt"

	"strings"

	"os"

	"path/filepath"

	"path"

	"io"

	"github.com/Unknwon/goconfig"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// gui主窗口
func main() {
	var inTE, outTE *walk.TextEdit
	// 加载配置
	conf, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		EchoERR("配置文件异常")
		return
	}

	// 获取各项配置
	title, err := conf.GetValue("Window", "title")
	if err != nil {
		EchoERR("[配置]title异常")
		return
	}

	width, err := conf.Int("Window", "width")
	if err != nil {
		EchoERR("[配置]width异常")
		return
	}

	height, err := conf.Int("Window", "height")
	if err != nil {
		EchoERR("[配置]height异常")
		return
	}

	outpath, err := conf.GetValue("Path", "outpath")
	if err != nil {
		EchoERR("[配置]outpath异常")
		return
	}

	inputpath, err := conf.GetValue("Path", "inputpath")
	if err != nil {
		EchoERR("[配置]inputpath异常")
		return
	}

	fmt.Println(outpath, inputpath)

	if PathExists(outpath) == false || PathExists(inputpath) == false {
		EchoERR("[路径]输入输出路径配置错误")
		return
	}

	MainWindow{
		Title:   title,
		MinSize: Size{width, height},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "获取",
				OnClicked: func() {
					//					outTE.SetText(strings.ToUpper(inTE.Text()))
					filesstr := inTE.Text()
					outTE.SetText(DoGetFile(inputpath, outpath, filesstr))
				},
			},
		},
	}.Run()
}

// 判断文件、文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func EchoERR(mag string) {
	fmt.Println(mag)
	time.Sleep(5 * time.Second)
}

// 程序主要进程
func DoGetFile(inputpath, outpath, filesstr string) string {
	// 判断要获取的文件或者文件夹是否存在
	filesarr := strings.Split(filesstr, "\r\n")
	for _, v := range filesarr {
		allPath := inputpath + v
		if !PathExists(allPath) {
			return "文件或文件夹不存在" + allPath
		}
	}
	err := WalkPath(inputpath, outpath, filesarr)
	if err != nil {
		return "运行时出错，请检查文件、文件夹"
	}

	return "成功！"
}

func WalkPath(rootpath, outpath string, path []string) error {
	for _, v := range path {
		fmt.Println(v)
		err := filepath.Walk(rootpath+v, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			_path := strings.Replace(path, rootpath, "", -1) // 获得相对路径
			outpathall := strings.Replace(outpath+_path, `\`, `/`, -1)
			CopyFile(outpathall, path) // 复制
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 复制文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	dir := path.Dir(dstName)
	//	fmt.Println("文件名：", dir, []byte(dstName), PathExists(dir))
	//	fmt.Println("文件名：", dir, dstName, PathExists(dir))
	if !PathExists(dir) {
		os.MkdirAll(dir, 0777) //创建文件夹
		fmt.Println("文件夹不存在:", dir)
	}
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
