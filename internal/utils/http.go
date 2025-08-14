package utils

import (
	"fmt"
	"io"
	"net/http"
)

// CheckHTTPResponse 检查HTTP响应状态码，如果不是2xx则返回错误
func CheckHTTPResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
