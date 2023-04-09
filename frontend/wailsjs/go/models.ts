export namespace database {
	
	export class Conversation {
	    id: number;
	    conversationSettingsID: number;
	    title: string;
	    // Go type: time
	    lastMessageTime: any;
	    generating: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Conversation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.conversationSettingsID = source["conversationSettingsID"];
	        this.title = source["title"];
	        this.lastMessageTime = this.convertValues(source["lastMessageTime"], null);
	        this.generating = source["generating"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConversationSetting {
	    id: number;
	    systemPromptTemplate: string;
	    toolsEnabled: string[];
	
	    static createFrom(source: any = {}) {
	        return new ConversationSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.systemPromptTemplate = source["systemPromptTemplate"];
	        this.toolsEnabled = source["toolsEnabled"];
	    }
	}
	export class Message {
	    id: number;
	    conversationID: number;
	    content: string;
	    author: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.conversationID = source["conversationID"];
	        this.content = source["content"];
	        this.author = source["author"];
	    }
	}
	export class Settings {
	    openAiApiKey: string;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.openAiApiKey = source["openAiApiKey"];
	        this.model = source["model"];
	    }
	}
	export class UpdateConversationSettingsParams {
	    systemPromptTemplate: string;
	    toolsEnabled: string[];
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConversationSettingsParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.systemPromptTemplate = source["systemPromptTemplate"];
	        this.toolsEnabled = source["toolsEnabled"];
	        this.id = source["id"];
	    }
	}

}

