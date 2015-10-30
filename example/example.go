package main

import (
	"fmt"
	"os"

	"github.com/f2prateek/go-circle"
	"github.com/f2prateek/go-pointers"
	"github.com/segmentio/pointer"
)

func main() {
	c := circle.New(os.Getenv("CIRCLE_TOKEN"))

	fmt.Println(c.Me())
	fmt.Println(c.Projects())
	fmt.Println(c.RecentBuilds())
	fmt.Println(c.RecentBuildsForProject("segmentio", "analytics-android"))
	fmt.Println(c.RecentBuildsForProjectBranch("segmentio", "analytics-android", "pull/346", circle.RecentBuildsOptions{
		Limit:  pointers.Int(2),
		Offset: pointers.Int(1),
		Filter: pointer.String("completed"),
	}))
	fmt.Println(c.BuildSummary("segmentio", "analytics-android", 345))
	fmt.Println(c.Artifacts("segmentio", "analytics-android", 345))
	fmt.Println(c.Retry("segmentio", "analytics-android", 346))
}
