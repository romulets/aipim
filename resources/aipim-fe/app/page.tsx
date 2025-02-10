'use client'

import { useState, useCallback, ChangeEvent } from "react"
import Image from "next/image";
import CodeEditor from '@uiw/react-textarea-code-editor';
import { CloudtrailLogMapping, MappedSource, MappedEvent } from "./domain"

// const serverUrl = "http://localhost:8777"
const serverUrl = ""

class Toggles {
  default: boolean = false;
  sources: number = -1;
}

class Stats {
  sources: number = 0
  events: number = 0
  targets: number = 0
  relatedEntities: number = 0
}

export default function Home() {
  const [cloudtrailMapping, setMapping] = useState(new CloudtrailLogMapping())
  const [stats, setStats] = useState(new Stats())
  const [toggles, setToggles] = useState(new Toggles());
  const [painlessCode, setPainlessCode] = useState("");

  const toggleDefault = useCallback(() => {
    setToggles(() => ({ default: !toggles.default, sources: -1 } as Toggles))
  }, [toggles])

  const updateStats = async (mapping: CloudtrailLogMapping) => {
    const sts = new Stats()
    sts.relatedEntities += mapping.defaultRelatedEntities.length
    sts.sources = mapping.sources.length
    mapping.sources.forEach(source => {
      sts.events += source.events.length
      sts.relatedEntities += source.relatedEntityFields.length
      source.events.forEach(event => {
        sts.targets += event.targetFields.length
      })
    })

    setStats(sts)
  }

  const toggleSource = (idx: number) => () => {
    if (toggles.sources == idx) {
      setToggles(() => (new Toggles()));
      return
    }

    setToggles(() => ({ default: false, sources: idx } as Toggles))
  }

  const arrowClassName = ((open: boolean) => open ? "collapse-arrow-open" : "collapse-arrow-closed")
  const collapseBody = ((open: boolean) => open ? "collapse-body-open" : "collapse-body-closed")

  let toPainlessAbort: AbortController | null = null;
  const updateCloudtrailMapping = (newMapping: CloudtrailLogMapping) => {
    if (toPainlessAbort != null) {
      toPainlessAbort!.abort()
    }

    toPainlessAbort = new AbortController()

    const mapping = { ...cloudtrailMapping, ...newMapping }
    setMapping(mapping);
    updateStats(mapping);

    (async () => {
      const resp = await fetch(`${serverUrl}/api/mapping/to-painless`, {
        method: 'POST',
        mode: "cors",
        body: JSON.stringify(mapping),
        headers: {
          "Content-Type": "application/json",
        },
        signal: toPainlessAbort!.signal
      })

      const painless = await resp.text();
      setPainlessCode(painless)

      toPainlessAbort = null
    })()
  }

  let fromPainlessAbort: AbortController | null = null;
  const getMappingBack = (code: string) => {
    if (fromPainlessAbort != null) {
      fromPainlessAbort!.abort()
    }

    fromPainlessAbort = new AbortController()
    setPainlessCode(code);

    (async () => {
      const resp = await fetch(`${serverUrl}/api/mapping/from-painless`, {
        method: 'POST',
        mode: "cors",
        body: code,
        signal: fromPainlessAbort!.signal
      })

      if (resp.status != 200) {
        const err = await resp.text()
        alert(err)
        return
      }

      const mapping = await resp.json();
      setMapping(mapping as CloudtrailLogMapping)
      updateStats(mapping);

      fromPainlessAbort = null
    })()
  }

  const addEventSource = () => {
    const idx = cloudtrailMapping.sources.length // no need to decrement because we are increasing anyway
    updateCloudtrailMapping({ sources: [...cloudtrailMapping.sources, new MappedSource("aws_source")] } as CloudtrailLogMapping)
    toggleSource(idx)()
  }

  return (
    <div className="h-screen">
      {/* HEADER */}
      <div className="header flex h-[7dvh] p-2">
        <div className="flex flex-1">
          <Image
            className="w-8 flex-none"
            src="/aipim-background.svg"
            alt="Aipim Logo"
            width={30}
            height={30}
            priority />

          <h1 className="flex-none p-2 mt-1"><b>Aipim</b></h1>
        </div>

        <div className="flex flex-2 mt-1">
          <b className="p-1 pt-2"> Sources:</b> 
          <span className="p-1 pt-2">{stats.sources}</span>
          <b className="p-1 pt-2"> Events:</b> 
          <span className="p-1 pt-2">{stats.events}</span>
          <b className="p-1 pt-2"> Related Entities:</b> 
          <span className="p-1 pt-2">{stats.relatedEntities}</span>
          <b className="p-1 pt-2"> Targets:</b> 
          <span className="p-1 pt-2">{stats.targets}</span>
        </div>
      </div>

      {/* MAIN */}
      <main className="flex h-[89dvh]">

        {/* FORM */}
        <section className="flex-1 p-5 overflow-y-auto">

          <div className="form-section mb-2 flex flex-col">
            <button onClick={toggleDefault} className="p-2 flex-1 text-left">
              <Image
                className={`inline mr-3 collapse-arrow ${arrowClassName(toggles.default)}`}
                src="arrow.svg"
                alt="icon collapse"
                width={5}
                height={5}
                priority />
              <b>Default Values</b>
            </button>

            <div className={`collapse-body ${collapseBody(toggles.default)}`}>

              <div className="flex p-2">
                <label className="mr-2 p-1 form-label" htmlFor="default-actor">Default actor</label>
                <input className="text-input p-1  w-4/5"
                  type="text" id="default-actor"
                  name="default-actor"
                  onChange={(e) => updateCloudtrailMapping({ defaultActor: e.target.value } as CloudtrailLogMapping)}
                  value={cloudtrailMapping.defaultActor}
                />
              </div>

              <div className="inc-list-header flex">
                <div className="flex-1"><label className="title">Default Related Entities</label></div>
                <button className="flex-2 main-btn add-btn" onClick={() => updateCloudtrailMapping({ defaultRelatedEntities: [...cloudtrailMapping.defaultRelatedEntities, ""] } as CloudtrailLogMapping)}>+ Add Related Entity</button>
              </div>

              <div className="inc-list p-4 pr-2 m-2 mt-5 mr-0">
                {cloudtrailMapping.defaultRelatedEntities.length == 0
                  ? <p>No default related entities</p>

                  : cloudtrailMapping.defaultRelatedEntities.map((related, idx) => {
                    const updateEntity = (e: ChangeEvent<HTMLInputElement>) => {
                      const entities = [...cloudtrailMapping.defaultRelatedEntities]
                      entities[idx] = e.target.value
                      updateCloudtrailMapping({ defaultRelatedEntities: entities } as CloudtrailLogMapping)
                    }

                    const removeEntity = () => {
                      const entities = [...cloudtrailMapping.defaultRelatedEntities]
                      entities.splice(idx, 1)
                      updateCloudtrailMapping({ defaultRelatedEntities: entities } as CloudtrailLogMapping)
                    }

                    return (<div className="" key={`default-related-${idx}`}>
                      <input className="text-input p-1"
                        type="text"
                        onChange={updateEntity}
                        value={cloudtrailMapping.defaultRelatedEntities[idx]} />

                      <button className="col-span-1 main-btn remove-btn" onClick={removeEntity}>-</button>

                    </div>)
                  })}
              </div>
            </div>
          </div>

          <div className="form-section mb-2 p-2 flex contrast">
            <h2 className="font-bold flex-1">Event Sources</h2>
            <button
              onClick={addEventSource}
              className="main-btn flex-3 sm">+ Add Event Source</button>
          </div>

          {cloudtrailMapping.sources.map((source, idx) => {
            const shrinkAll = toggles.sources == idx

            const removeSource = () => {
              const sources = [...cloudtrailMapping.sources]
              sources.splice(idx, 1)
              updateCloudtrailMapping({ sources: sources } as CloudtrailLogMapping)

              if (shrinkAll) {
                setToggles(new Toggles())
              }
            }

            const updateSource = (source: MappedSource) => {
              cloudtrailMapping.sources[idx] = { ...cloudtrailMapping.sources[idx], ...source }
              updateCloudtrailMapping({ ...cloudtrailMapping })
            }

            const addEvent = () => {
              source.events = [...cloudtrailMapping.sources[idx].events, new MappedEvent(`${source.sourceName}Event`)]
              updateSource(source)
            }

            return (
              <div className="form-section mb-2" key={idx}>
                <div className="flex">
                  <button onClick={toggleSource(idx)} className="p-2 flex-1 text-left">
                    <Image
                      className={`inline mr-3 collapse-arrow ${arrowClassName(toggles.sources == idx)}`}
                      src="arrow.svg"
                      alt="icon collapse"
                      width={5}
                      height={5}
                      priority />
                    <b>{source.sourceName}</b> <small>event source</small>
                  </button>

                  <div className="flex-3 content-center mr-1">
                    <button className="main-btn sm" onClick={removeSource}>- Remove Source</button>
                  </div>
                </div>

                <div className={`collapse-body ${collapseBody(toggles.sources == idx)}`}>

                  <div className="flex p-2">
                    <label className="mr-2 p-1 form-label" htmlFor={`source-name-${idx}`}>Source name</label>
                    <input className="text-input p-1"
                      type="text"
                      id={`source-name-${idx}`}
                      name={`source-name-${idx}`}
                      onChange={(e) => updateSource({ sourceName: e.target.value } as MappedSource)}
                      value={cloudtrailMapping.sources[idx].sourceName}
                    />
                  </div>


                  <div className="inc-list-header flex">
                    <div className="flex-1"><label className="title">Related Entities</label></div>
                    <button className="flex-2 main-btn add-btn" onClick={() => updateSource({ relatedEntityFields: [...source.relatedEntityFields, ""] } as MappedSource)}>+ Add Related Entity</button>
                  </div>

                  <div className="inc-list p-4 pr-2 m-2 mt-5 mr-0">

                    {source.relatedEntityFields.length == 0
                      ? <p>No specific related entities</p>

                      : source.relatedEntityFields.map((related, idx) => {
                        const updateEntity = (e: ChangeEvent<HTMLInputElement>) => {
                          const entities = [...source.relatedEntityFields]
                          entities[idx] = e.target.value
                          updateSource({ relatedEntityFields: entities } as MappedSource)
                        }

                        const removeEntity = () => {
                          const entities = [...source.relatedEntityFields]
                          entities.splice(idx, 1)
                          updateSource({ relatedEntityFields: entities } as MappedSource)
                        }

                        return (<div className="" key={`default-related-${idx}`}>
                          <input className="text-input p-1"
                            type="text"
                            onChange={updateEntity}
                            value={source.relatedEntityFields[idx]} />

                          <button className="col-span-1 main-btn remove-btn" onClick={removeEntity}>-</button>
                        </div>)
                      })}
                  </div>

                  <div className="ml-3 pl-1 pb-1 mt-4 flex title-divider">
                    <h2 className="flex-1">Events</h2>
                    <button className="main-btn sm flex-2" onClick={addEvent}>+ Add Event</button>
                  </div>

                  <div className="ml-3 events-holder">
                    {source.events.map((event, eventIdx) => {
                      const updateEvent = (event: MappedEvent) => {
                        source.events[eventIdx] = { ...source.events[eventIdx], ...event }
                        updateSource({ ...source })
                      }

                      const removeEvent = () => {
                        const events = [...source.events]
                        events.splice(idx, 1)
                        updateSource({ events: events } as MappedSource)
                      }

                      return (
                        <div key={`${idx}-${eventIdx}`} className="pl-5 pb-1 pt-4">
                          <div className="pb-1 flex title-sub-divider">
                            <h2 className="flex-1 font-bold">{event.eventName}</h2>
                            <button className="main-btn sm flex-2" onClick={removeEvent}>- Remove Event</button>
                          </div>

                          <div className="grid grid-cols-6 pl-2 pt-2">
                            <label className="mr-2 p-1 form-label" htmlFor={`event-name-${idx}-${eventIdx}`}>Event Name</label>
                            <input className="col-span-5 text-input p-1"
                              type="text"
                              id={`event-name-${idx}-${eventIdx}`}
                              name={`event-name-${idx}-${eventIdx}`}
                              onChange={(e) => updateEvent({ eventName: e.target.value } as MappedEvent)}
                              value={event.eventName}
                            />
                          </div>

                          <div className="grid grid-cols-6 pl-2 pt-2 pb-2">
                            <label className="mr-2 p-1 form-label" htmlFor={`actor${idx}-${eventIdx}`}>Specific Actor</label>
                            <input className="col-span-5 text-input p-1"
                              type="text"
                              id={`actor${idx}-${eventIdx}`}
                              name={`actor${idx}-${eventIdx}`}
                              onChange={(e) => updateEvent({ actorField: e.target.value } as MappedEvent)}
                              value={event.actorField}
                            />
                          </div>

                          <div className="inc-list-header flex">
                            <div className="flex-1"><label className="title">Targets</label></div>
                            <button className="flex-2 main-btn add-btn" onClick={() => updateEvent({ targetFields: [...event.targetFields, ""] } as MappedEvent)}>+ Add Target</button>
                          </div>
                          <div className="inc-list p-4 pr-2 m-2 mt-5 mr-0">

                            {event.targetFields.length == 0
                              ? <p>No targets</p>

                              : event.targetFields.map((related, idx) => {
                                const updateTarget = (e: ChangeEvent<HTMLInputElement>) => {
                                  const targets = [...event.targetFields]
                                  targets[idx] = e.target.value
                                  updateEvent({ targetFields: targets } as MappedEvent)
                                }

                                const removeTarget = () => {
                                  const targets = [...event.targetFields]
                                  targets.splice(idx, 1)
                                  updateEvent({ targetFields: targets } as MappedEvent)
                                }

                                return (<div className="" key={`${idx}`}>

                                  <input className="text-input p-1"
                                    type="text"
                                    onChange={updateTarget}
                                    value={event.targetFields[idx]} />

                                  <button className="col-span-1 main-btn remove-btn" onClick={removeTarget}>-</button>

                                </div>)
                              })}
                          </div>

                        </div>
                      )

                    })}
                  </div>

                </div>
              </div>
            )
          })}


        </section>

        <CodeEditor className="flex-1"
          value={painlessCode}
          language="dart"
          placeholder="// code here"
          onChange={(e) => getMappingBack(e.target.value)}
          padding={15}
          data-color-mode="dark"
          style={{
            fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
            overflowY: 'scroll'
          }}
        />
      </main>

      {/* FOOTER */}
      <footer className="footer h-[4dvh] text-xs  p-2 text-center">Aipim - <b>AWS Interfaced Painless Identifier Mapping</b> - Made by @romulets and @kubasobon (Elastic 2025)</footer>
    </div >
  );
}
