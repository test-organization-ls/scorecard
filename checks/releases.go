package checks

import (
	"strings"

	"github.com/dlorenc/scorecard/checker"
	"github.com/google/go-github/v32/github"
)

func init() {
	AllChecks = append(AllChecks, checker.NamedCheck{
		Name: "Signed-Releases",
		Fn:   SignedReleases,
	})
}

func SignedReleases(c checker.Checker) checker.CheckResult {
	releases, _, err := c.Client.Repositories.ListReleases(c.Ctx, c.Owner, c.Repo, &github.ListOptions{})
	if err != nil {
		return checker.RetryResult(err)
	}

	totalReleases := 0
	totalSigned := 0
	for _, r := range releases {
		assets, _, err := c.Client.Repositories.ListReleaseAssets(c.Ctx, c.Owner, c.Repo, r.GetID(), &github.ListOptions{})
		if err != nil {
			return checker.RetryResult(err)
		}
		if len(assets) <= 1 {
			continue
		}
		totalReleases++
		signed := false
		for _, asset := range assets {
			for _, suffix := range []string{".sig", ".minisig"} {
				c.Logf("signed release found: %s", asset)
				if strings.HasSuffix(asset.GetName(), suffix) {
					signed = true
					break
				}
			}
			if signed {
				totalSigned++
				break
			}
		}
	}

	if totalReleases == 0 {
		return checker.InconclusiveResult
	}
	return checker.ProportionalResult(totalSigned, totalReleases, .75)
}
