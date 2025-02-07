package domain

import (
	_ "embed" // embed templates
	"strings"
)

type CloudtrailLogMapping struct {
	DefaultRelatedEntities []string
	DefaultActor           string
	Sources                []MappedSource
}

type MappedSource struct {
	SourceName          string
	Events              []MappedEvent
	RelatedEntityFields []string
}

type MappedEvent struct {
	EventName    string
	TargetFields []string
	ActorField   *string
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

func (m *CloudtrailLogMapping) toString() string {
	functions := make([]string, 0, len(m.Sources))
	calls := make([]string, 0, len(m.Sources))

	for _, source := range m.Sources {
		fn, call := source.buildFn()
		functions = append(functions, fn)
		calls = append(calls, call)
	}

	script := strings.Replace(scriptTemplate, functionDefPH, strings.Join(functions, "\n\n"), 1)
	script = strings.Replace(script, functionCallPH, strings.Join(calls, "\n"), 1)
	script = strings.Replace(script, defaultActorPH, m.DefaultActor, 1)

	defaultRelated := make([]string, 0, len(m.DefaultRelatedEntities))
	for _, fieldName := range m.DefaultRelatedEntities {
		defaultRelated = append(defaultRelated, addFieldCall(contextRelated, fieldName))
	}

	script = strings.Replace(script, defaultRelatedEntitiesPH, strings.Join(defaultRelated, "\n"), 1)

	return script
}

// Returns function definition and call
func (m *MappedSource) buildFn() (string, string) {
	function := strings.Replace(functionTemplate, sourceNamePH, m.SourceName, 3)
	related := make([]string, 0, len(m.RelatedEntityFields))

	for idx, fieldName := range m.RelatedEntityFields {
		ident := ""
		if idx > 0 {
			ident = "  "
		}
		related = append(related, ident+addFieldCall(contextRelated, fieldName))
	}

	events := make([]string, 0, len(m.Events))

	for _, event := range m.Events {
		events = append(events, event.buildIfCase())
	}

	body := strings.Join(related, "\n")
	body += "\n  \n  "
	body += strings.ReplaceAll(strings.Join(events, " else "), "\n", "\n  ")
	function = strings.Replace(function, functionBodyPH, body, 1)

	call := strings.Replace(functionCallTemplate, sourceNamePH, m.SourceName, 1)

	return function, call
}

func (e *MappedEvent) buildIfCase() string {
	ifCase := strings.Replace(ifEventTemplate, eventNamePH, e.EventName, 1)
	targetFields := make([]string, 0, len(e.TargetFields))

	for _, fieldName := range e.TargetFields {
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
