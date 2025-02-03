package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	buildVersion        = "dev"
	buildCommit         = "unknown"
	buildTime    string = "0" // Change to string type
)

func getBuildTimeUnix() int64 {
	// Convert buildTime string to int64, with fallback
	timestamp, err := strconv.ParseInt(buildTime, 10, 64)
	if err != nil {
		return time.Now().Unix()
	}
	return timestamp
}

func GetFormattedBuildTime() string {
	buildTime := getBuildTimeUnix()
	if buildTime == 0 {
		return time.Now().Format("2006-01-02 15:04:05 MST")
	}

	buildTimeUTC := time.Unix(buildTime, 0).UTC()
	localTime := buildTimeUTC.Local()
	return localTime.Format("2006-01-02 15:04:05 MST")
}

func GetVersionString() string {
	// For snapshot or dev builds
	if buildVersion == "dev" || strings.Contains(buildVersion, "snapshot") ||
		strings.Contains(buildVersion, "next") {
		return fmt.Sprintf("%s (snapshot build %s)",
			buildVersion,
			time.Now().Local().Format("2006-01-02"),
		)
	}

	// Normal release version formatting
	buildTime := getBuildTimeUnix()
	if buildTime > 0 {
		buildDate := time.Unix(buildTime, 0).UTC().Local().Format("2006-01-02")
		shortCommit := buildCommit
		if len(buildCommit) > 7 {
			shortCommit = buildCommit[:7]
		}
		return fmt.Sprintf("%s (%s %s)", buildVersion, buildDate, shortCommit)
	}

	// Fallback
	return buildVersion
}
