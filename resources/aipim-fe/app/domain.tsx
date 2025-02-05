export class CloudtrailLogMapping {
    defaultRelatedEntities: Array<string> = [];
    defaultActor: string = "";
    sources: Array<MappedSource> = [];
}

export class MappedSource {
    sourceName: string = "";
    events: Array<MappedEvent> = [];
    relatedEntityFields: Array<string> = [];
}

export class MappedEvent {
    eventName: string = "";
    targetFields: Array<string> = [];
    actorField?: string;
}