package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"util"
	//"path/filepath"
)

func TarPkg(topmost, input_path, output_file string, tw *tar.Writer) {
	// file write

	if tw == nil {
		// topmost dir create tar file
		fw, err := os.Create(output_file)
		if err != nil {
			panic(err)
		}
		defer fw.Close()

		// gzip write
		gw := gzip.NewWriter(fw)
		defer gw.Close()

		// tar write
		tw = tar.NewWriter(gw)
		defer tw.Close()
	}

	// 打开文件夹
	dir, err := os.Open(input_path)
	if err != nil {
		panic(input_path)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			//continue
			fmt.Println("go into subdir", dir.Name()+"\\"+fi.Name())
			TarPkg(topmost, dir.Name()+"\\"+fi.Name(), "", tw)
		} else {
			// 打印文件名称
			fmt.Println("deal with file", fi.Name())

			// 打开文件
			fr, err := os.Open(dir.Name() + "/" + fi.Name())
			if err != nil {
				panic(err)
			}
			defer fr.Close()

			// 信息头
			h := new(tar.Header)
			entire_path := dir.Name() + "/" + fi.Name()
			need_path := entire_path[len(topmost)+1:]
			h.Name = need_path
			h.Size = fi.Size()
			h.Mode = int64(fi.Mode())
			h.ModTime = fi.ModTime()

			// 写信息头
			err = tw.WriteHeader(h)
			if err != nil {
				panic(err)
			}

			// 写文件
			_, err = io.Copy(tw, fr)
			if err != nil {
				panic(err)
			}
		}

	}

	fmt.Println("tar.gz ok")
}

func TarJsonAndDmp(ini Config) {
	//TarPkg(ini.TarInputPath, ini.TarInputPath, ini.TarOutputPath, nil)
	TarFile(ini.TarInputPath, ini.TarOutputPath)

}

func TarFile(in, out string) {
	TarPkg(in, in, out, nil)
}

func UntarFile(in, out string) {
	fmt.Println("UntarFile")
	util.PrepareDir(out)
	created_dir := make(map[string]interface{})
	// file read
	fr, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// 显示文件
		fmt.Println(h.Name)

		//创建子目录，if necessary

		index := strings.LastIndex(h.Name, "/")
		if index != -1 {
			_, ok := created_dir[h.Name[:index]]
			if !ok {
				fmt.Println("creat sub dir", h.Name[:index])
				created_dir[h.Name[:index]] = nil
				os.MkdirAll(out+"\\"+h.Name[:index], 0777)
			}
		}

		// 打开文件
		fw, err := os.OpenFile(out+"\\"+h.Name, os.O_CREATE|os.O_WRONLY, 0644 /*os.FileMode(h.Mode)*/)
		if err != nil {
			panic(err)
		}
		defer fw.Close()

		// 写文件
		_, err = io.Copy(fw, tr)
		if err != nil {
			panic(err)
		}

	}

	fmt.Println("un tar.gz ok")
}

func CollectUnitTest() {
	TarFile("L:\\temp", "L:\\temp.tar.gz")
	UntarFile("L:\\collect.tar.gz", "L:\\untar")
}
