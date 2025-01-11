export namespace domain {
	
	export class CreateHabitDTO {
	    Name: string;
	    Amount: number;
	    Unit: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateHabitDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
	        this.Unit = source["Unit"];
	    }
	}
	export class Habit {
	    ID: number;
	    Name: string;
	    Amount: number;
	    Unit: string;
	    // Go type: time
	    ArchivedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Habit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
	        this.Unit = source["Unit"];
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
	export class SetTrackerDto {
	    HabitID: number;
	    Amount: number;
	
	    static createFrom(source: any = {}) {
	        return new SetTrackerDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HabitID = source["HabitID"];
	        this.Amount = source["Amount"];
	    }
	}
	export class Tracker {
	    HabitID: number;
	    Amount: number;
	    // Go type: time
	    At: any;
	
	    static createFrom(source: any = {}) {
	        return new Tracker(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HabitID = source["HabitID"];
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

