package miniapp_test

import (
	"fmt"
	"testing"

	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/miniapp"
)

func TestVerify(t *testing.T) {
	v := miniapp.Verify(`query_id=AAEsQU88AgAAACxBTzwGzfAV&user=%7B%22id%22%3A5306794284%2C%22first_name%22%3A%22Nikola%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22nickrsg%22%2C%22language_code%22%3A%22en%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1719441507&hash=5000a598a319b9a2856a15cb5dbc99b3c1c3d2178c68eb41119bd18d3b238d1b`)
	fmt.Printf("v: %v\n", v)
	t.FailNow()
}
