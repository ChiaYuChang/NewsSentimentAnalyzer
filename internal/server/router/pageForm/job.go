package pageform

import (
	"fmt"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
)

type JobPager struct {
	JStatusStr string  `form:"jstatus" validate:"required"`
	FromJId    int32   `form:"fjid" validate:"min=0"`
	ToJId      int32   `form:"tjid" validate:"min=0"`
	Page       int     `form:"page" validate:"min=0"`
	Direction  bool    `form:"direction"`
	JIdsStr    string  `form:"jids"`
	JIds       []int32 `form:"-"`
}

func (jp *JobPager) ParseJIds() error {
	if jp.JIdsStr == "" {
		return nil
	}

	jids := strings.Split(jp.JIdsStr, ",")
	jp.JIds = make([]int32, 0, len(jids))
	for _, j := range jids {
		i, err := convert.StrTo(strings.TrimSpace(j)).Int32()
		if err != nil {
			return fmt.Errorf("error while parsing jids %w", err)
		}
		jp.JIds = append(jp.JIds, i)
	}
	return nil
}
