package version

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Build struct {
	CL      int
	Release string
}

type contextKey struct{}

const firstSeasonTwoChangelist = 3785000

var (
	changelistPattern = regexp.MustCompile(`(?:CL-|\d+\.\d+\.\d+-)(\d+)`)
	releasePattern    = regexp.MustCompile(`Release-(\d+(?:\.\d+)?)`)
)

func Parse(userAgent string) Build {
	var build Build
	if m := changelistPattern.FindStringSubmatch(userAgent); m != nil {
		build.CL, _ = strconv.Atoi(m[1])
	}
	if m := releasePattern.FindStringSubmatch(userAgent); m != nil {
		build.Release = m[1]
	}
	return build
}

func (b Build) Season() int {
	if b.Release == "" {
		if b.CL >= firstSeasonTwoChangelist {
			return 2
		}
		return 1
	}
	major, _, _ := strings.Cut(b.Release, ".")
	season, err := strconv.Atoi(major)
	if err != nil || season < 1 {
		return 1
	}
	return season
}

func (b Build) AtLeast(major, minor int) bool {
	majorText, minorText, _ := strings.Cut(b.Release, ".")
	m, err := strconv.Atoi(majorText)
	if err != nil {
		return false
	}
	n, _ := strconv.Atoi(minorText)
	return m > major || m == major && n >= minor
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey{}, Parse(r.UserAgent()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func From(ctx context.Context) Build {
	build, _ := ctx.Value(contextKey{}).(Build)
	return build
}
