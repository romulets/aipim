'use client'

import { useState, useCallback, ChangeEvent } from "react"
import Image from "next/image";
import CodeEditor from '@uiw/react-textarea-code-editor';
import { CloudtrailLogMapping } from "./domain"

class Toggles {
  default: boolean = false;
  sources: boolean = false;
}

export default function Home() {
  const [cloudtrailMapping, setMapping] = useState(new CloudtrailLogMapping())

  const [toggles, setToggles] = useState(new Toggles());

  const expandDefaultSection = useCallback(() => {
    setToggles(() => ({ ...toggles, default: !toggles.default } as Toggles))
  }, [toggles])

  const arrowClassName = ((open: boolean) => open ? "collapse-arrow-open" : "collapse-arrow-closed")
  const collapseBody = ((open: boolean) => open ? "collapse-body-open" : "collapse-body-closed")

  const updateCloudtrailMapping = ((newMapping: CloudtrailLogMapping) => setMapping({ ...cloudtrailMapping, ...newMapping }))

  return (
    <div className="h-screen">
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

      <main className="flex h-[89dvh]">

        <section className="flex-1 p-5">

          <div className="form-section mb-2">
            <button onClick={expandDefaultSection} className="p-2">
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
                <button className="inc-list-btn add-btn" onClick={() => updateCloudtrailMapping({ defaultRelatedEntities: [...cloudtrailMapping.defaultRelatedEntities, ""] } as CloudtrailLogMapping)}>Add +</button>

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

                      <button className="inc-list-btn" onClick={removeEntity}>- Remove</button>

                      <input className="text-input p-1 ml-5 w-5/6"
                        type="text"
                        onChange={updateEntity}
                        value={cloudtrailMapping.defaultRelatedEntities[idx]} />

                    </div>)
                  })}
              </div>

            </div>
          </div>

          {/* <div className="form-section mb-2">
            <button onClick={expandDefaultSources} className="p-2">
              <Image
                className={`inline mr-3 collapse-arrow ${arrowClassName(toggles.sources)}`}
                src="arrow.svg"
                alt="icon collapse"
                width={5}
                height={5}
                priority />
              Sources
            </button>

            <div className={`collapse-body ${collapseBody(toggles.sources)}`}>
              Test
            </div>
          </div> */}

        </section>



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

      <footer className="footer h-[4dvh] text-xs  p-2 text-center">Aipim - <b>AWS Interfaced Painless Identifier Mapping</b> - Made by @romulets and @kubasobon (Elastic 2025)</footer>
    </div>
  );
}
