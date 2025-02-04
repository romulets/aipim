package domain

import (
	_ "embed" // embed templates
	"strings"
)

type cloudtrailLogMapping struct {
	defaultRelatedEntities []string
	defaultActor           string
	sources                []mappedSource
}

type mappedSource struct {
	sourceName          string
	events              []mappedEvent
	relatedEntityFields []string
}

type mappedEvent struct {
	eventName    string
	targetFields []string
	actorField   *string
}

//go:embed templates/script.painless
var scriptTemplate string

//go:embed templates/function.painless
var functionTemplate string

//go:embed templates/function_call.painless
var functionCallTemplate string

//go:embed templates/add_field.painless
var addFieldTemplate string

//go:embed templates/if_event.painless
var ifEventTemplate string

const (
	functionDefPH            = "%%FUNCTIONS_DEFINITIONS%%"
	functionCallPH           = "%%FUNCTIONS_CALLS%%"
	sourceNamePH             = "%%SOURCE_NAME%%"
	functionBodyPH           = "%%FUNCTION_BODY%%"
	fieldNamePH              = "%%FIELD_NAME%%"
	contextPH                = "%%CONTEXT%%"
	eventNamePH              = "%%EVENT_NAME%%"
	ifBodyPH                 = "%%IF_BODY%%"
	defaultActorPH           = "%%DEFAULT_ACTOR%%"
	defaultRelatedEntitiesPH = "%%DEFAULT_RELATED_ENTITIES%%"

	contextRelated = "related"
	contextTarget  = "target"
)

func (m *cloudtrailLogMapping) toString() string {
	functions := make([]string, 0, len(m.sources))
	calls := make([]string, 0, len(m.sources))

	for _, source := range m.sources {
		fn, call := source.buildFn()
		functions = append(functions, fn)
		calls = append(calls, call)
	}

	script := strings.Replace(scriptTemplate, functionDefPH, strings.Join(functions, "\n\n"), 1)
	script = strings.Replace(script, functionCallPH, strings.Join(calls, "\n"), 1)
	script = strings.Replace(script, defaultActorPH, m.defaultActor, 1)

	defaultRelated := make([]string, 0, len(m.defaultRelatedEntities))
	for _, fieldName := range m.defaultRelatedEntities {
		defaultRelated = append(defaultRelated, addFieldCall(contextRelated, fieldName))
	}

	script = strings.Replace(script, defaultRelatedEntitiesPH, strings.Join(defaultRelated, "\n"), 1)

	return script
}

// Returns function definition and call
func (m *mappedSource) buildFn() (string, string) {
	function := strings.Replace(functionTemplate, sourceNamePH, m.sourceName, 3)
	related := make([]string, 0, len(m.relatedEntityFields))

	for idx, fieldName := range m.relatedEntityFields {
		ident := ""
		if idx > 0 {
			ident = "  "
		}
		related = append(related, ident+addFieldCall(contextRelated, fieldName))
	}

	events := make([]string, 0, len(m.events))

	for _, event := range m.events {
		events = append(events, event.buildIfCase())
	}

	body := strings.Join(related, "\n")
	body += "\n  \n  "
	body += strings.ReplaceAll(strings.Join(events, " else "), "\n", "\n  ")
	function = strings.Replace(function, functionBodyPH, body, 1)

	call := strings.Replace(functionCallTemplate, sourceNamePH, m.sourceName, 1)

	return function, call
}

func (e *mappedEvent) buildIfCase() string {
	ifCase := strings.Replace(ifEventTemplate, eventNamePH, e.eventName, 1)
	targetFields := make([]string, 0, len(e.targetFields))

	for _, fieldName := range e.targetFields {
		targetFields = append(targetFields, addFieldCall(contextTarget, fieldName))
	}

	text := strings.ReplaceAll(strings.Join(targetFields, "\n"), "\n", "\n  ")
	return strings.Replace(ifCase, ifBodyPH, text, 1)
}

func addFieldCall(context string, fieldName string) string {
	if strings.Contains(fieldName, "[]") {
		return addArrayFieldCall(context, fieldName)
	}
	addField := strings.Replace(addFieldTemplate, fieldNamePH, fieldName, 1)
	return strings.Replace(addField, contextPH, context, 1)
}

// Can't handle multilevel arrays
func addArrayFieldCall(context string, fieldName string) string {
	chunks := strings.Split(fieldName, "[]")

	b := strings.Builder{}
	b.WriteString(`field("`)
	b.WriteString(chunks[0])
	b.WriteString(`").get(new ArrayList())
  .stream().forEach(f -> addValue(enrichCtx.`)
	b.WriteString(context)
	b.WriteString(`, f`)
	b.WriteString(chunks[1])
	b.WriteString(`));`)

	return b.String()
}
