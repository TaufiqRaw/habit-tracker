export namespace domain {
	
	export enum RestDayModeEnum {
	    Week = "WEEK",
	    Month = "MONTH",
	}
	export class CreateHabitDTO {
	    Name: string;
	    MinPerDay: number;
	    Unit: string;
	    RestDay: number;
	    RestDayMode: RestDayModeEnum;
	
	    static createFrom(source: any = {}) {
	        return new CreateHabitDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.MinPerDay = source["MinPerDay"];
	        this.Unit = source["Unit"];
	        this.RestDay = source["RestDay"];
	        this.RestDayMode = source["RestDayMode"];
	    }
	}
	export class HabitNode {
	    Id: number;
	    MinPerDay: number;
	    Unit: string;
	    RestDay: number;
	    RestDayMode: RestDayModeEnum;
	    // Go type: time
	    StartAt: any;
	    // Go type: time
	    ArchivedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new HabitNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.MinPerDay = source["MinPerDay"];
	        this.Unit = source["Unit"];
	        this.RestDay = source["RestDay"];
	        this.RestDayMode = source["RestDayMode"];
	        this.StartAt = this.convertValues(source["StartAt"], null);
	        this.ArchivedAt = this.convertValues(source["ArchivedAt"], null);
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
	export class Habit {
	    Id: number;
	    Name: string;
	    NodeIDs: number[];
	    NodeLength: number;
	    LastNode: HabitNode;
	
	    static createFrom(source: any = {}) {
	        return new Habit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.NodeIDs = source["NodeIDs"];
	        this.NodeLength = source["NodeLength"];
	        this.LastNode = this.convertValues(source["LastNode"], HabitNode);
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
	
	export class SetTrackerDto {
	    HabitId: number;
	    Amount: number;
	
	    static createFrom(source: any = {}) {
	        return new SetTrackerDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HabitId = source["HabitId"];
	        this.Amount = source["Amount"];
	    }
	}
	export class Tracker {
	    HabitId: number;
	    Amount: number;
	    // Go type: time
	    At: any;
	
	    static createFrom(source: any = {}) {
	        return new Tracker(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HabitId = source["HabitId"];
	        this.Amount = source["Amount"];
	        this.At = this.convertValues(source["At"], null);
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

export namespace service {
	
	export class TrackerArrResult {
	    Data: domain.Tracker[];
	    Err?: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackerArrResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], domain.Tracker);
	        this.Err = source["Err"];
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
	export class habitArrResult {
	    Data: domain.Habit[];
	    Err?: string;
	
	    static createFrom(source: any = {}) {
	        return new habitArrResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], domain.Habit);
	        this.Err = source["Err"];
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
	export class habitNodeResult {
	    Data?: domain.HabitNode;
	    Err?: string;
	
	    static createFrom(source: any = {}) {
	        return new habitNodeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], domain.HabitNode);
	        this.Err = source["Err"];
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
	export class habitResult {
	    Data?: domain.Habit;
	    Err?: string;
	
	    static createFrom(source: any = {}) {
	        return new habitResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], domain.Habit);
	        this.Err = source["Err"];
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

