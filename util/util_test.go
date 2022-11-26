package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FlowingSPDG/gotv-plus-go/util"
)

func TestParseToken(t *testing.T) {
	asserts := assert.New(t)
	for _, td := range []struct {
		title    string
		input    string
		steamid  string
		unixTime int64
	}{
		{
			title:    "通常の解析パターン1",
			input:    "s90152525936315402t1635312048",
			steamid:  "90152525936315402",
			unixTime: 1635312048,
		},
	} {
		t.Run("Parse:"+td.title, func(t *testing.T) {
			sid, time, err := util.ParseToken(td.input)
			asserts.NoError(err)
			asserts.Equal(td.steamid, sid)
			asserts.Equal(td.unixTime, time.Unix())
		})
	}
}
