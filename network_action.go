package chromy

import (
	"context"
	"regexp"
)

var _ URLPattern = &regexp.Regexp{}

type URLPattern interface {
	MatchString(string) bool
}

func WaitResource(method string, urlpattern URLPattern) Action {
	return NewResource(ResourcePattern(method, urlpattern), ResourceDone())
}

func WaitResourceWithMatcher(matcher func(*Request) bool) Action {
	return NewResource(ResourceMatch(matcher), ResourceDone())
}

func CaptureRequests(macher func(*Request) bool, onRequest func(*Request)) Action {
	return ActionFunc(func(ctx context.Context, t *Target) error {
		stopper := func(r *Request) bool {
			if macher(r) {
				onRequest(r)
			}

			// continually receive request
			return false
		}

		t.domain.Network.watch(stopper)

		return nil
	})
}
