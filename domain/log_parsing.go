package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var sourceNamePattern = regexp.MustCompile(`// -+ (\w+) Events -+`)
var addFieldPattern = regexp.MustCompile(`\s+addField\(.+, "(.+)"\)`)
var eventNamePattern = regexp.MustCompile(`if \(eventName == "(.+)"\)`)

func (ms *mappedSource) addRelatedEntityField(f string) {
	ms.relatedEntityFields = append(ms.relatedEntityFields, f)
}

func (clm *cloudtrailLogMapping) scan(painless string) {
	var inDefinitions, inSetup bool
	// var indentLevel int
	// var sourceScope string
	// var eventScope string

	// namedSources := map[string]*mappedSource{}
	// namedEvents := map[string]*mappedEvent{}
	var source *mappedSource
	var event *mappedEvent

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
				matches := sourceNamePattern.FindStringSubmatch(line)
				if len(matches) != 2 {
					panic(fmt.Sprintf(
						"! Should find two elements in %d\n%s\n%s",
						i, line, sourceNamePattern,
					))
				}
				source = &mappedSource{
					sourceName:          matches[1],
					events:              []mappedEvent{},
					relatedEntityFields: []string{},
				}
				// sourceScope = matches[1]
				i += 5
				continue
			}
			if strings.HasPrefix(line, "  addField") {
				// parse line and add to related on source-lvl
				matches := addFieldPattern.FindStringSubmatch(line)
				if len(matches) != 2 {
					panic(fmt.Sprintf(
						"! Should find two elements in %d\n%s\n%s",
						i, line, addFieldPattern,
					))
				}
				source.relatedEntityFields = append(source.relatedEntityFields, matches[1])
				i += 1
				continue
			}
			if strings.HasPrefix(line, "  }") && event != nil {
				// finalize event, add to source
				source.events = append(source.events, *event)
				event = nil
			}
			if strings.Contains(line, "eventName ==") {
				// go to eventScope
				matches := eventNamePattern.FindStringSubmatch(line)
				if len(matches) != 2 {
					panic(fmt.Sprintf(
						"! Should find two elements in %d\n%s\n%s",
						i, line, addFieldPattern,
					))
				}
				event = &mappedEvent{
					eventName:    matches[1],
					targetFields: []string{},
					actorField:   nil,
				}
				// eventScope = matches[1]
				i += 1
				continue
			}
			if strings.HasPrefix(line, "    addField") {
				// parse line and add to related on event-lvl
				matches := addFieldPattern.FindStringSubmatch(line)
				if len(matches) != 2 {
					panic(fmt.Sprintf(
						"! Should find two elements in %d\n%s\n%s",
						i, line, addFieldPattern,
					))
				}
				event.targetFields = append(event.targetFields, matches[1])
				i += 1
				continue
			}
			if strings.HasPrefix(line, "}") {
				// finalize source, add to clm
				clm.sources = append(clm.sources, *source)
				source = nil
				i += 1
				continue
			}
		} else {
			i += 1
			continue
		}
	}
}
