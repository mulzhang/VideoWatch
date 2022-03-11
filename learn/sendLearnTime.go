package learn

import (
	"WatchVideo/course"
	"WatchVideo/datasource"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type Learn struct {
	db *sqlx.DB
}

func NewLearn() *Learn {
	return &Learn{
		db: datasource.ConnectMysql(),
	}
}

//开始学习
func (l *Learn) ExecLearn() {
	record := datasource.SelGroupRecord(l.db)
	var wg sync.WaitGroup
	wg.Add(len(record))
	for _, v := range record {
		if v.CourseTime < 900 {
			//relateIdRefModelDataMap := sync.Map{}
			////执行查询
			//relateIdRefModelDataMap.Store(v.CourseNo, v.CourseTime)
			go func(key datasource.CourseTime) {
				defer wg.Done()
				//获取到当前课程对应的学习视频信息
				groupRecord := datasource.SelRecord(l.db, v.CourseNo)
				//定义视频的下标
				nowIndex := int32(0)
				//当前视频播放进度
				nowPlay := float64(0)
				//定时执行处理
				ticker := time.NewTicker(30 * time.Second)
				for _ = range ticker.C {
					l.PlayVideo(groupRecord, &nowIndex, &nowPlay)
				}
			}(v)
		}

	}
	wg.Wait()
}

//播放视频
func (l *Learn) PlayVideo(key []course.RequestStudy, nowIndex *int32, nowPlay *float64) {

}
