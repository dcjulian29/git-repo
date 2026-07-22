/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shared

import (
	"fmt"
	"time"
)

const (
	ageDay   = 24 * time.Hour
	ageMonth = 30 * ageDay
	ageYear  = 365 * ageDay
)

// Age renders a duration as a short, human-readable age such as "3d", "5mo12d",
// or "1y2mo". Negative durations are treated as zero and durations under a
// minute render as "0m".
func Age(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	switch {
	case d >= ageYear:
		years := int(d / ageYear)
		months := int((d % ageYear) / ageMonth)

		if months == 0 {
			return fmt.Sprintf("%dy", years)
		}

		return fmt.Sprintf("%dy%dmo", years, months)
	case d >= ageMonth:
		months := int(d / ageMonth)
		days := int((d % ageMonth) / ageDay)

		if days == 0 {
			return fmt.Sprintf("%dmo", months)
		}

		return fmt.Sprintf("%dmo%dd", months, days)
	case d >= ageDay:
		return fmt.Sprintf("%dd", int(d/ageDay))
	case d >= time.Hour:
		return fmt.Sprintf("%dh", int(d/time.Hour))
	default:
		return fmt.Sprintf("%dm", int(d/time.Minute))
	}
}
