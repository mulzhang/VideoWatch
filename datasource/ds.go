package datasource

import (
	"WatchVideo/course"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	userName  string = "root"
	password  string = "1234"
	ipAddrees string = "120.79.159.146"
	port      int    = 3320
	dbName    string = "course"
	charset   string = "utf8"
)

type CourseTime struct {
	CourseNo   string
	CourseTime float32
}

func ConnectMysql() *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	return Db
}

func AddRecord(Db *sqlx.DB, rs []*course.RequestStudy) {
	for _, v := range rs {
		result, err := Db.Exec("insert into my_course_info(user_id, course_no,course_name, courlib_no, k_point_no, video_url, video_time)  values(?,?,?,?,?,?,?)", 201693412300013, v.CourseNo, v.CourseName, v.CourlibNo, v.KPointNo, v.VideoUrl, v.VideoTime)
		if err != nil {
			fmt.Sprintf("data insert faied, error:[%s]", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		fmt.Sprintf("insert success, last id:[%d]\n", id)
	}
}

func SelGroupRecord(Db *sqlx.DB) []CourseTime {
	result, err := Db.Query("select course_no,course_time_vzb from my_course group by course_no,course_time_vzb")
	if err != nil {
		fmt.Sprintf("data SelGroupRecord, error:[%s]", err.Error())
	}
	var courseNo, courseTime interface{}
	var qh []CourseTime
	for result.Next() {
		time := CourseTime{}
		// to do
		result.Scan(&courseNo, &courseTime)
		if courseNo != nil {
			time.CourseNo = courseNo.(string)
		}
		if courseTime != nil {
			time.CourseTime = courseTime.(float32)
		}
		qh = append(qh, time)
	}
	return qh
}

func SelRecord(Db *sqlx.DB, no string) []course.RequestStudy {
	result, err := Db.Query("select user_id, course_no, courlib_no, k_point_no, video_url, video_time, course_name from my_course_info where courlib_no = ?", no)
	if err != nil {
		fmt.Sprintf("data SelRecord, error:[%s]", err.Error())
	}
	var userId, courseNo, courlibNo, kPointNo, videoTime, courseName interface{}
	var qh []course.RequestStudy
	for result.Next() {
		time := course.RequestStudy{}
		// to do
		result.Scan(&userId, &courseNo, &courlibNo, &kPointNo, &videoTime, &courseName)
		if courseNo != nil {
			time.CourseNo = courseNo.(string)
		}
		if userId != nil {
			time.UserId = userId.(int64)
		}
		if courlibNo != nil {
			time.CourlibNo = courlibNo.(int64)
		}
		if kPointNo != nil {
			time.KPointNo = kPointNo.(int64)
		}
		if videoTime != nil {
			time.VideoTime = videoTime.(float64)
		}
		if courseName != nil {
			time.CourseName = courseName.(string)
		}
		qh = append(qh, time)
	}
	return qh
}
