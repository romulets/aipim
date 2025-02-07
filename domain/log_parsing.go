package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	sourceNamePattern = regexp.MustCompile(`// -+ (\w+) Events -+`)
	addFieldPattern   = regexp.MustCompile(`\s*addField\(.+, "(.+)"\)`)
	eventNamePattern  = regexp.MustCompile(`if \(eventName == "(.+)"\)`)
	actorPattern      = regexp.MustCompile(`enrichCtx.actor = field\("(.+)"\)`)
	fieldPattern      = regexp.MustCompile(`field\("(.+)"\)`)
	forEachPattern    = regexp.MustCompile(`f -> addValue\(.+, f\.(.+)\)\);`)
)

func (clm *CloudtrailLogMapping) scan(painless string) {
	var inDefinitions, inSetup bool
	var source *MappedSource
	var event *MappedEvent

	lines := strings.Split(painless, "\n")
	i := 0
	for i < len(lines) {
		line := lines[i]
		if !inDefinitions && strings.Contains(line, "- FUNCTIONS DEFINITIONS -") {
			inDefinitions = true
			i += 1
			continue
		}
		if !inSetup && strings.Contains(line, "- BASIC SETUP -") {
			inSetup = true
			inDefinitions = false
			i += 1
			continue
		}
		if inSetup && strings.Contains(line, "- FUNCTIONS CALLS -") {
			break
		}
		if strings.TrimSpace(line) == "" {
			i += 1
			continue
		}

		if inDefinitions {
			if strings.Contains(line, " Events -") {
				x := mustMatch(sourceNamePattern, line, 1)
				source = &MappedSource{
					SourceName:          x,
					Events:              []MappedEvent{},
					RelatedEntityFields: []string{},
				}
				i += 5
				continue
			}
			if strings.HasPrefix(line, "  addField") {
				// parse line and add to related on source-lvl
				x := mustMatch(addFieldPattern, line, 1)
				source.RelatedEntityFields = append(source.RelatedEntityFields, x)
				i += 1
				continue
			}
			if strings.HasPrefix(line, "  }") && event != nil {
				// finalize event, add to source
				source.Events = append(source.Events, *event)
				event = nil
			}
			if strings.Contains(line, "eventName ==") {
				// go to eventScope
				x := mustMatch(eventNamePattern, line, 1)
				event = &MappedEvent{
					EventName:    x,
					TargetFields: []string{},
					ActorField:   nil,
				}
				i += 1
				continue
			}
			if strings.HasPrefix(line, "    addField") {
				// parse line and add to related on event-lvl
				x := mustMatch(addFieldPattern, line, 1)
				event.TargetFields = append(event.TargetFields, x)
				i += 1
				continue
			}
			if strings.Contains(line, "ArrayList()") {
				prefix := mustMatch(fieldPattern, line, 1)
				suffix := mustMatch(forEachPattern, lines[i+1], 1)
				event.TargetFields = append(
					event.TargetFields,
					fmt.Sprintf("%s[].%s", prefix, suffix),
				)
				i += 2
				continue
			}
			if strings.HasPrefix(line, "}") {
				// finalize source, add to clm
				clm.Sources = append(clm.Sources, *source)
				source = nil
				i += 1
				continue
			}
		} else if inSetup {
			if strings.HasPrefix(line, "enrichCtx.actor =") {
				clm.DefaultActor = mustMatch(actorPattern, line, 1)
				i += 1
				continue
			}
			if strings.HasPrefix(line, "addField(") {
				x := mustMatch(addFieldPattern, line, 1)
				clm.DefaultRelatedEntities = append(clm.DefaultRelatedEntities, x)
				i += 1
				continue
			}
			if strings.Contains(line, "ArrayList()") {
				prefix := mustMatch(fieldPattern, line, 1)
				suffix := mustMatch(forEachPattern, lines[i+1], 1)
				clm.DefaultRelatedEntities = append(
					clm.DefaultRelatedEntities,
					fmt.Sprintf("%s[].%s", prefix, suffix),
				)
				i += 2
				continue
			}
		}
		i += 1
		continue
	}
}

func mustMatch(pattern *regexp.Regexp, s string, match int) string {
	matches := pattern.FindStringSubmatch(s)
	if len(matches) <= match {
		panic(fmt.Sprintf(
			"Could not match %q with %q (group %d)",
			s, pattern, match,
		))
	}
	return matches[match]
}
