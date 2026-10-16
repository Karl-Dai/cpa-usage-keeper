package cpa

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const upgradeCompatibilityTimeout = 3 * time.Second

var cpaReleaseVersion = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

// CheckUpgradeCompatibility 只为未完成的旧库升级诊断最低 CPA 版本，不授予网络可用性保证。
// 未知/dev 版本和暂时无法连接返回未验证；只有正式版本明确低于合同才阻止破坏性迁移。
func (c *Client) CheckUpgradeCompatibility(ctx context.Context) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if c == nil || c.baseURL == "" || c.managementKey == "" {
		return false, nil
	}
	probeCtx, cancel := context.WithTimeout(ctx, upgradeCompatibilityTimeout)
	defer cancel()
	// CPA 8.0.12 保留的轻量只读 debug 接口同样携带版本头，无需枚举凭证。
	// 只发一次请求；不追随重定向，也不读取或输出响应正文。
	client := *c.httpClient
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	probe := func() (string, error) {
		req, err := http.NewRequestWithContext(probeCtx, http.MethodGet, c.baseURL+"/v0/management/debug", nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.managementKey)
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return strings.TrimSpace(resp.Header.Get("X-CPA-VERSION")), nil
	}
	version, err := probe()
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return false, nil
	}
	parts := cpaReleaseVersion.FindStringSubmatch(version)
	if parts == nil {
		return false, nil
	}
	var numbers [3]uint64
	for i := range numbers {
		value, err := strconv.ParseUint(parts[i+1], 10, 64)
		if err != nil {
			return false, nil
		}
		numbers[i] = value
	}
	if numbers[0] < 8 || numbers[0] == 8 && numbers[1] == 0 && numbers[2] < 12 {
		return false, fmt.Errorf("CPA %d.%d.%d does not meet the minimum CPA 8.0.12 requirement; upgrade CPA and verify CPA_BASE_URL before retrying the database upgrade", numbers[0], numbers[1], numbers[2])
	}
	return true, nil
}
