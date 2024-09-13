package request

import (
	"github.com/linbaozhong/model-gen/pkg/utils"
	"testing"
)

func TestPostMoney(t *testing.T) {
	var i = "0.21"
	t.Log(utils.YuanString2Fen(i))
}
