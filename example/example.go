package main

import (
	"fmt"
	"os"

	"github.com/f2prateek/go-circle"
)

func main() {
	circle := circle.New(os.Getenv("CIRCLE_TOKEN"))

	fmt.Println(circle.Me())
	fmt.Println(circle.Projects())
	fmt.Println(circle.RecentBuilds())
	fmt.Println(circle.RecentBuildsForProject("segmentio", "analytics-android"))
	fmt.Println(circle.RecentBuildsForProjectBranch("segmentio", "analytics-android", "pull/382"))
	fmt.Println(circle.BuildSummary("segmentio", "analytics-android", 345))
	fmt.Println(circle.Artifacts("segmentio", "analytics-android", 345))
	fmt.Println(circle.Retry("segmentio", "analytics-android", 346))
}
