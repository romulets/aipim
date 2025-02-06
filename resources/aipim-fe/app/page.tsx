'use client'

import { useState, useCallback, ChangeEvent } from "react"
import Image from "next/image";
import CodeEditor from '@uiw/react-textarea-code-editor';
import { CloudtrailLogMapping, MappedSource } from "./domain"

class Toggles {
  default: boolean = false;
  sources: number = -1;
}

export default function Home() {
  const [cloudtrailMapping, setMapping] = useState(new CloudtrailLogMapping())

  const [toggles, setToggles] = useState(new Toggles());

  const toggleDefault = useCallback(() => {
    setToggles(() => ({ default: !toggles.default, sources: -1 } as Toggles))
  }, [toggles])


  const toggleSource = (idx: number) => () => {
    if (toggles.sources == idx) {
      setToggles(() => (new Toggles()));
      return
    }

    setToggles(() => ({ default: false, sources: idx } as Toggles))
  }

  const arrowClassName = ((open: boolean) => open ? "collapse-arrow-open" : "collapse-arrow-closed")
  const collapseBody = ((open: boolean) => open ? "collapse-body-open" : "collapse-body-closed")

  const updateCloudtrailMapping = ((newMapping: CloudtrailLogMapping) => setMapping({ ...cloudtrailMapping, ...newMapping }))

  return (
    <div className="h-screen">
      {/* HEADER */}
      <div className="header flex h-[7dvh] p-2">
        <Image
          className="w-8 flex-none"
          src="/aipim-background.svg"
          alt="Aipim Logo"
          width={30}
          height={30}
          priority />

        <h1 className="flex-none p-2">Aipim</h1>
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

              <div className="inc-list p-4 m-2 mt-5">
                <label className="title">Default Related Entities</label>
                <button className="main-btn add-btn" onClick={() => updateCloudtrailMapping({ defaultRelatedEntities: [...cloudtrailMapping.defaultRelatedEntities, ""] } as CloudtrailLogMapping)}>Add +</button>

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

                    return (<div className="m-3" key={`default-related-${idx}`}>

                      <button className="main-btn" onClick={removeEntity}>- Remove</button>

                      <input className="text-input p-1 ml-5 w-5/6"
                        type="text"
                        onChange={updateEntity}
                        value={cloudtrailMapping.defaultRelatedEntities[idx]} />

                    </div>)
                  })}
              </div>
            </div>
          </div>

          <div className="form-section mb-2 p-2 flex contrast">
            <h2 className="font-bold flex-1">Event Sources</h2>
            <button
              onClick={() => updateCloudtrailMapping({ sources: [...cloudtrailMapping.sources, new MappedSource("aws_source")] } as CloudtrailLogMapping)}
              className="main-btn flex-3 sm">+ Add</button>
          </div>

          {cloudtrailMapping.sources.map((source, idx) => {
            const removeSource = () => {
              const sources = [...cloudtrailMapping.sources]
              sources.splice(idx, 1)
              updateCloudtrailMapping({ sources: sources } as CloudtrailLogMapping)

              if (toggles.sources == idx) {
                setToggles(new Toggles())
              }
            }

            const updateSource = (source: MappedSource) => {
              cloudtrailMapping.sources[idx] = { ...cloudtrailMapping.sources[idx], ...source }
              console.log(source)
              updateCloudtrailMapping({ ...cloudtrailMapping })
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
                    <button className="main-btn sm" onClick={removeSource}>- Remove</button>
                  </div>
                </div>

                <div className={`collapse-body ${collapseBody(toggles.sources == idx)}`}>

                  <div className="flex p-2">
                    <label className="mr-2 p-1 form-label" htmlFor={`source-name-${idx}`}>Source name</label>
                    <input className="text-input p-1"
                      type="text" id="default-actor"
                      name={`source-name-${idx}`}
                      onChange={(e) => updateSource({ sourceName: e.target.value } as MappedSource)}
                      value={cloudtrailMapping.sources[idx].sourceName}
                    />
                  </div>

                </div>
              </div>
            )
          })}


        </section>


        {/* CODE EDITOR */}
        {/* <CodeEditor className="flex-1"
          value={JSON.stringify(cloudtrailMapping, null, 2)}
          // language="dart"
          language="json"
          placeholder="// code here"
          // onChange={(evn) => setCode(evn.target.value)}
          padding={15}
          data-color-mode="dark"
          style={{
            fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
            overflowY: 'scroll'
          }}
        /> */}

        <CodeEditor className="flex-1"
          // value={JSON.stringify(cloudtrailMapping, null, 2)}
          language="dart"
          // language="json"
          placeholder="// code here"
          // onChange={(evn) => setCode(evn.target.value)}
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
    </div>
  );
}
