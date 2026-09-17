export namespace app {
	
	export class FileEntry {
	    name: string;
	    path: string;
	    isDir: boolean;
	    size: number;
	    modTime: string;
	    ext: string;
	
	    static createFrom(source: any = {}) {
	        return new FileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.modTime = source["modTime"];
	        this.ext = source["ext"];
	    }
	}
	export class CleanupPreview {
	    total: number;
	    size: number;
	    files: FileEntry[];
	
	    static createFrom(source: any = {}) {
	        return new CleanupPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.size = source["size"];
	        this.files = this.convertValues(source["files"], FileEntry);
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
	export class CleanupRequest {
	    paths: string[];
	    keyword: string;
	    exts: string[];
	    before: string;
	    minSize: number;
	    maxSize: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanupRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paths = source["paths"];
	        this.keyword = source["keyword"];
	        this.exts = source["exts"];
	        this.before = source["before"];
	        this.minSize = source["minSize"];
	        this.maxSize = source["maxSize"];
	    }
	}
	export class CleanupResult {
	    deleted: number;
	    failed: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deleted = source["deleted"];
	        this.failed = source["failed"];
	    }
	}
	export class Config {
	    port: number;
	    saveDir: string;
	    autoFolder: boolean;
	    ip: string;
	    wifiSsid: string;
	    wifiPass: string;
	    wifiSec: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.saveDir = source["saveDir"];
	        this.autoFolder = source["autoFolder"];
	        this.ip = source["ip"];
	        this.wifiSsid = source["wifiSsid"];
	        this.wifiPass = source["wifiPass"];
	        this.wifiSec = source["wifiSec"];
	    }
	}
	
	export class FilePreview {
	    previewable: boolean;
	    kind: string;
	    mime: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FilePreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.previewable = source["previewable"];
	        this.kind = source["kind"];
	        this.mime = source["mime"];
	        this.content = source["content"];
	    }
	}
	export class ServerStatus {
	    running: boolean;
	    ips: string[];
	    ip: string;
	    port: number;
	    url: string;
	    qrDataUri: string;
	    saveDir: string;
	    uploadSub: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.ips = source["ips"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.url = source["url"];
	        this.qrDataUri = source["qrDataUri"];
	        this.saveDir = source["saveDir"];
	        this.uploadSub = source["uploadSub"];
	    }
	}
	export class TrashInfo {
	    count: number;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new TrashInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.size = source["size"];
	    }
	}

}

