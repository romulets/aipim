export class CloudtrailLogMapping {
    defaultRelatedEntities: Array<string> = [];
    defaultActor: string = "";
    sources: Array<MappedSource> = [];
}

export class MappedSource {
    sourceName: string = "";
    events: Array<MappedEvent> = [];
    relatedEntityFields: Array<string> = [];

    constructor(sourceName?: string) {
        if (!!sourceName) {
            this.sourceName = sourceName!
        }
        
    }
}

export class MappedEvent {
    eventName: string = "";
    targetFields: Array<string> = [];
    actorField: string = "";

    constructor(eventName?: string) {
        if (!!eventName) {
            this.eventName = eventName!
        }
        
    }
}