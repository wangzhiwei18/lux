package bcy

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
)

func TestDownload(t *testing.T) {
	// 1. 激活 httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.ActivateNonDefault(http.DefaultClient)

	// 2. 测试数据
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://bcy.net/item/detail/6558738153367142664",
				Title: "cos正片 命运石之门 牧濑红莉栖 克里斯蒂娜…",
				Size:  13035762, // 需要根据模拟图片大小调整
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 3. 为当前测试注册模拟响应
			setupMockResponses(tt.args.URL)
			
			// 4. 执行测试
			data, err := New().Extract(tt.args.URL, extractors.Options{})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
		
		// 5. 重置模拟，避免测试间干扰
		httpmock.Reset()
	}
}

// setupMockResponses 设置所有需要的模拟响应
func setupMockResponses(bcyURL string) {
	// 模拟 BCY 页面 HTML 响应
	mockHTML := `
<!DOCTYPE html>
<html>
<head>
	<title>cos正片 命运石之门 牧濑红莉栖 克里斯蒂娜… - 半次元 banciyuan - ACG爱好者社区</title>
</head>
<body>
	<script>
		JSON.parse("{\"detail\":{\"post_data\":{\"multi\":[{\"original_path\":\"https://img.example.com/pic1.jpg\"},{\"original_path\":\"https://img.example.com/pic2.jpg\"}]}}}");
	</script>
</body>
</html>`
	
	httpmock.RegisterResponder("GET", bcyURL,
		httpmock.NewStringResponder(200, mockHTML))
	
	// 模拟图片响应（用于 request.Size 调用）
	// 注意：request.Size 通常发送 HEAD 请求获取 Content-Length
	image1URL := "https://img.example.com/pic1.jpg"
	image2URL := "https://img.example.com/pic2.jpg"
	
	// 模拟图片1的 HEAD 请求（获取大小）
	httpmock.RegisterResponder("HEAD", image1URL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Length", "6517881") // 6.5MB
			return resp, nil
		})
	
	// 模拟图片2的 HEAD 请求
	httpmock.RegisterResponder("HEAD", image2URL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Length", "6517881") // 6.5MB
			return resp, nil
		})
	
	// 如果需要，也可以模拟 GET 请求（更完整）
	httpmock.RegisterResponder("GET", image1URL,
		httpmock.NewStringResponder(200, "fake image content"))
	httpmock.RegisterResponder("GET", image2URL,
		httpmock.NewStringResponder(200, "fake image content"))
}
