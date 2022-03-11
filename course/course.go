package course

import (
	"WatchVideo/stu"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type ResponseCourse struct {
	Flag     bool     `json:"flag"`
	Msg      string   `json:"msg"`
	MyCourse []Course `json:"mycourse"`
}

type Course struct {
	CheckMod           string `json:"checkMod"`
	Zystatus           string `json:"zystatus"`
	CourseCateporyName string `json:"courseCateporyName"`
	Status             string `json:"status"`
	State              string `json:"state"`
	Score              string `json:"score"`
	CourId             string `json:"courId"`
	CourNo             string `json:"courNo"`
}

type ResponseStudy struct {
	Flag bool   `json:"flag"`
	Msg  string `json:"msg"`
	Data []struct {
		Id       string `json:"id"`
		VideoUrl string `json:"videoUrl"`
		Name     string `json:"name"`
		CourseNo string `json:"courseNo"`
		NodeType string `json:"nodeType"`
	} `json:"data"`
}

//获取课程的结构体
type RequestStudy struct {
	CourseNo   string  `json:"courseNo"`
	CourseName string  `json:"courseName"`
	UserId     int64   `json:"userId" default=201693412300013`
	CourlibNo  int64   `json:"courlibNo"`
	KPointNo   int64   `json:"kPointNo"`
	VideoUrl   string  `json:"videoUrl"`
	VideoTime  float64 `json:"videoTime"`
}

//存储课程信息
type Kpointdetail struct {
	KpointNo  int64
	CourlibNo int64
	CourseNo  int64
}

//获取课程信息
func GetResCourse(re io.ReadCloser) *ResponseCourse {
	bodyText, err := ioutil.ReadAll(re)
	if err != nil {
		log.Fatal(err)
	}
	var r ResponseCourse
	err = json.Unmarshal(bodyText, &r)
	if err != nil {
		log.Fatal(err)
	}
	return &r
}

//获取Kpoint数组信息
func GetKpo(re io.ReadCloser, s *stu.Student) []*Kpointdetail {
	var kpointDetails []*Kpointdetail
	for _, v := range GetResCourse(re).MyCourse {
		no := v.CourNo

		//获取课程号
		bytes := s.SendHttp("", "GET", "http://wljy.whut.edu.cn/web/ucenterdetail.htm?id="+no)
		reader, _ := goquery.NewDocumentFromReader(bytes)
		reader.Find("button[onclick]").Each(func(i int, s *goquery.Selection) {
			val := s.Nodes[0].Attr[2].Val
			index := strings.Index(val, "(")
			i2 := strings.Index(val, ",'")
			content := val[index+1 : i2]
			split := strings.Split(content, ",")
			kpointNo, _ := strconv.Atoi(split[0])
			courlibNo, _ := strconv.Atoi(split[1])
			courseNo, _ := strconv.Atoi(split[2])
			kpointdetail := &Kpointdetail{
				KpointNo:  int64(kpointNo),
				CourlibNo: int64(courlibNo),
				CourseNo:  int64(courseNo),
			}
			kpointDetails = append(kpointDetails, kpointdetail)
		})
	}
	return kpointDetails
}

//按课程分组 获取到每个课程对应的资源信息
func GetGroup(kpo []*Kpointdetail, s *stu.Student) map[string][]*RequestStudy {
	//记录需要可以播放的所有视频
	m := make(map[string][]*RequestStudy)
	for _, k := range kpo {
		formatInt := strconv.FormatInt(k.KpointNo, 10)
		sendHttp := s.SendHttp(`kpointNo=`+formatInt, "POST", "http://wljy.whut.edu.cn/edu/eduCourseKpoint/findKpointDataByNo.ajax")
		bodyText, err := ioutil.ReadAll(sendHttp)
		if err != nil {
			log.Fatal(err)
			continue
		}
		var rs ResponseStudy
		err = json.Unmarshal(bodyText, &rs)
		if err != nil {
			fmt.Println("解析失败,{}", bodyText)
			continue
		}
		for _, d := range rs.Data {
			if d.NodeType != "video" {
				continue
			}
			//组装学习对象
			study := RequestStudy{
				CourseNo:   d.CourseNo,
				CourseName: d.Name,
				CourlibNo:  k.CourlibNo,
				KPointNo:   k.KpointNo,
				VideoUrl:   d.VideoUrl,
			}
			//以课程分组
			capital, ok := m[d.CourseNo]
			if ok {
				capital = append(capital, &study)
				m[d.CourseNo] = capital
			} else {
				m[d.CourseNo] = []*RequestStudy{&study}
			}
		}
	}
	return m
}
