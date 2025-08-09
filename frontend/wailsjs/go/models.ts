export namespace models {
	
	export class Attachment {
	    title: string;
	    title_link: string;
	    text: string;
	    image_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Attachment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.title_link = source["title_link"];
	        this.text = source["text"];
	        this.image_url = source["image_url"];
	    }
	}
	export class Element {
	    type: string;
	    text?: string;
	    user_id?: string;
	    url?: string;
	    name?: string;
	    elements?: Element[];
	
	    static createFrom(source: any = {}) {
	        return new Element(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.text = source["text"];
	        this.user_id = source["user_id"];
	        this.url = source["url"];
	        this.name = source["name"];
	        this.elements = this.convertValues(source["elements"], Element);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class Block {
	    type: string;
	    block_id: string;
	    elements: Element[];
	
	    static createFrom(source: any = {}) {
	        return new Block(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.block_id = source["block_id"];
	        this.elements = this.convertValues(source["elements"], Element);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ChannelPin {
	    id: string;
	    type: string;
	    created: number;
	    user: string;
	    owner: string;
	
	    static createFrom(source: any = {}) {
	        return new ChannelPin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.created = source["created"];
	        this.user = source["user"];
	        this.owner = source["owner"];
	    }
	}
	export class ChannelPurpose {
	    value: string;
	    creator: string;
	    last_set: number;
	
	    static createFrom(source: any = {}) {
	        return new ChannelPurpose(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.creator = source["creator"];
	        this.last_set = source["last_set"];
	    }
	}
	export class ChannelTopic {
	    value: string;
	    creator: string;
	    last_set: number;
	
	    static createFrom(source: any = {}) {
	        return new ChannelTopic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.creator = source["creator"];
	        this.last_set = source["last_set"];
	    }
	}
	export class Channel {
	    id: string;
	    name: string;
	    created: number;
	    creator: string;
	    is_archived: boolean;
	    is_general: boolean;
	    members: string[];
	    topic: ChannelTopic;
	    purpose: ChannelPurpose;
	    pins?: ChannelPin[];
	
	    static createFrom(source: any = {}) {
	        return new Channel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created = source["created"];
	        this.creator = source["creator"];
	        this.is_archived = source["is_archived"];
	        this.is_general = source["is_general"];
	        this.members = source["members"];
	        this.topic = this.convertValues(source["topic"], ChannelTopic);
	        this.purpose = this.convertValues(source["purpose"], ChannelPurpose);
	        this.pins = this.convertValues(source["pins"], ChannelPin);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	
	
	export class Edited {
	    user: string;
	    ts: string;
	
	    static createFrom(source: any = {}) {
	        return new Edited(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = source["user"];
	        this.ts = source["ts"];
	    }
	}
	
	export class File {
	    id: string;
	    name: string;
	    title: string;
	    mimetype: string;
	    url_private: string;
	    user: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.mimetype = source["mimetype"];
	        this.url_private = source["url_private"];
	        this.user = source["user"];
	    }
	}
	export class Reply {
	    user: string;
	    ts: string;
	
	    static createFrom(source: any = {}) {
	        return new Reply(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = source["user"];
	        this.ts = source["ts"];
	    }
	}
	export class Reaction {
	    name: string;
	    users: string[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new Reaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.users = source["users"];
	        this.count = source["count"];
	    }
	}
	export class UserProfile {
	    first_name: string;
	    real_name: string;
	    image_72: string;
	    team: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UserProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.first_name = source["first_name"];
	        this.real_name = source["real_name"];
	        this.image_72 = source["image_72"];
	        this.team = source["team"];
	        this.name = source["name"];
	    }
	}
	export class Message {
	    subtype?: string;
	    user?: string;
	    text?: string;
	    inviter?: string;
	    type?: string;
	    ts: string;
	    client_msg_id?: string;
	    team?: string;
	    user_profile?: UserProfile;
	    blocks?: Block[];
	    reactions?: Reaction[];
	    thread_ts?: string;
	    parent_user_id?: string;
	    reply_count?: number;
	    replies?: Reply[];
	    edited?: Edited;
	    files?: File[];
	    attachments?: Attachment[];
	    user_name?: string;
	    avatar_url?: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.subtype = source["subtype"];
	        this.user = source["user"];
	        this.text = source["text"];
	        this.inviter = source["inviter"];
	        this.type = source["type"];
	        this.ts = source["ts"];
	        this.client_msg_id = source["client_msg_id"];
	        this.team = source["team"];
	        this.user_profile = this.convertValues(source["user_profile"], UserProfile);
	        this.blocks = this.convertValues(source["blocks"], Block);
	        this.reactions = this.convertValues(source["reactions"], Reaction);
	        this.thread_ts = source["thread_ts"];
	        this.parent_user_id = source["parent_user_id"];
	        this.reply_count = source["reply_count"];
	        this.replies = this.convertValues(source["replies"], Reply);
	        this.edited = this.convertValues(source["edited"], Edited);
	        this.files = this.convertValues(source["files"], File);
	        this.attachments = this.convertValues(source["attachments"], Attachment);
	        this.user_name = source["user_name"];
	        this.avatar_url = source["avatar_url"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	

}

