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

func (clm *CloudtrailLogMapping) Scan(painless string) error {
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
				x, err := mustMatch(sourceNamePattern, line, 1)
				if err != nil {
					return err
				}
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
				x, err := mustMatch(addFieldPattern, line, 1)
				if err != nil {
					return err
				}
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
				x, err := mustMatch(eventNamePattern, line, 1)
				if err != nil {
					return err
				}
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
				x, err := mustMatch(addFieldPattern, line, 1)
				if err != nil {
					return err
				}
				event.TargetFields = append(event.TargetFields, x)
				i += 1
				continue
			}
			if strings.Contains(line, "ArrayList()") {
				prefix, err := mustMatch(fieldPattern, line, 1)
				if err != nil {
					return err
				}
				suffix, err := mustMatch(forEachPattern, lines[i+1], 1)
				if err != nil {
					return err
				}
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
				var err error
				clm.DefaultActor, err = mustMatch(actorPattern, line, 1)
				if err != nil {
					return err
				}
				i += 1
				continue
			}
			if strings.HasPrefix(line, "addField(") {
				x, err := mustMatch(addFieldPattern, line, 1)
				if err != nil {
					return err
				}
				clm.DefaultRelatedEntities = append(clm.DefaultRelatedEntities, x)
				i += 1
				continue
			}
			if strings.Contains(line, "ArrayList()") {
				prefix, err := mustMatch(fieldPattern, line, 1)
				if err != nil {
					return err
				}
				suffix, err := mustMatch(forEachPattern, lines[i+1], 1)
				if err != nil {
					return err
				}
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
	return nil
}

func mustMatch(pattern *regexp.Regexp, s string, match int) (string, error) {
	matches := pattern.FindStringSubmatch(s)
	if len(matches) <= match {
		return "", fmt.Errorf(
			"Could not match %q with %q (group %d)",
			s, pattern, match,
		)
	}
	return matches[match], nil
}
