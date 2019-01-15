package image_uploader

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestDownloadImage(t *testing.T) {
	tests := []struct {
		url string
	}{
		{"https://gss0.bdstatic.com/5bVWsj_p_tVS5dKfpU_Y_D3/res/r/image/2017-09-27/297f5edb1e984613083a2d3cc0c5bb36.png"},
	}

	for _, test := range tests {
		f, size, err := DownloadImage(test.url)
		if err != nil {
			t.Error(err)
		}
		_, _ = f.Seek(0, io.SeekStart)
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Error(err)
		}
		RemoveFile(f)
		if len(b) != int(size) {
			t.Errorf("len(%d) != Size(%d)", len(b), int(size))
		}
	}
}
