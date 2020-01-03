package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
)

func main() {
	var imgPath string
	var targetName string
	fmt.Println("请输入源图片名称:")
	fmt.Scanln(&imgPath)
	fmt.Println("请输入目标文件名称:")
	fmt.Scanln(&targetName)
	img2txt(imgPath, 200, []string{"@", "#", "*", "%", "+", ",", ".", " "}, "\n", targetName)
}

/**
 * imagePath: 图片路径
 * size: 生成文本后的尺寸(1代表一个像素,1个像素会被替换成1个字符)
 * txts: 将像素处理成字符列表
 * rowend: 换行字符(windows和linux不同)
 * output: 生成文本文件保存路径
 */
func img2txt(imgPath string, size uint, txts []string, rowend string, output string) {
	file, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("打开图片失败: ", err)
		return
	}

	defer file.Close()

	var img image.Image
	ext := path.Ext(imgPath)

	switch ext {
	case ".JPG":
		fallthrough
	case ".JPEG":
		fallthrough
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(file)
	case ".PNG":
	case ".png":
		img, err = png.Decode(file)
	default:
		err = errors.New("目前只支持png或者jpg的解码")
	}

	if err != nil {
		fmt.Println("图片解码失败: ", err)
		return
	}

	var width = size
	var height = (size * uint(img.Bounds().Dy())) / (uint(img.Bounds().Dx()))
	height = height * 6 / 10
	newimg := resize.Resize(width, height, img, resize.Lanczos3) //根据高宽resize图片，并得到新图片的像素值
	dx := newimg.Bounds().Dx()
	dy := newimg.Bounds().Dy()

	//创建一个字节buffer，一会用来保存字符
	textBuffer := bytes.Buffer{}

	//遍历图片每一行每一列像素
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			colorRgb := newimg.At(x, y)
			r, g, b, _ := colorRgb.RGBA()

			//获得三原色的值，算一个平均数出来
			avg := uint8((r + g + b) / 3 >> 8)
			//有多少个用来替换的字符就将256分为多少个等分，然后计算这个像素的平均值趋紧与哪个字符，最后，将这个字符添加到字符buffer里
			num := avg / uint8(256/len(txts))
			textBuffer.WriteString(txts[num])
		}

		textBuffer.WriteString(rowend) //一行结束，换行
	}

	//将字符buffer的数据写入到文本文件里，结束。
	f, err := os.Create(output)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	_, err = f.WriteString(textBuffer.String())
	if err != nil {
		fmt.Println("生成文件出错: ", err)
	} else {
		fmt.Println("文件生成成功")
	}
}
