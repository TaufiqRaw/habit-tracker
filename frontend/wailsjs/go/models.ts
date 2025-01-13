export namespace domain {
	
	export enum RestDayModeEnum {
	    WEEKLY = "Weekly",
	    MONTHLY = "Monthly",
	}
	export class CreateHabitDTO {
	    Name: string;
	    Amount: number;
	    Unit: string;
	    RestDay: number;
	    RestDayMode: RestDayModeEnum;
	    LastHabitID?: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateHabitDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
	        this.Unit = source["Unit"];
	        this.RestDay = source["RestDay"];
	        this.RestDayMode = source["RestDayMode"];
	        this.LastHabitID = source["LastHabitID"];
	    }
	}
	export class Habit {
	    Id: number;
	    LastHabitId?: number;
	    Name: string;
	    Amount: number;
	    Unit: string;
	    RestDay: number;
	    RestDayMode: RestDayModeEnum;
	    // Go type: time
	    StartAt: any;
	    // Go type: time
	    ArchivedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Habit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.LastHabitId = source["LastHabitId"];
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
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
	export class HabitNode {
	    Id: number;
	    LastHabitId?: number;
	    Name: string;
	    Amount: number;
	    Unit: string;
	    RestDay: number;
	    RestDayMode: RestDayModeEnum;
	    // Go type: time
	    StartAt: any;
	    // Go type: time
	    ArchivedAt?: any;
	    PreviousHabit?: Habit;
	    NextHabit?: Habit;
	
	    static createFrom(source: any = {}) {
	        return new HabitNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.LastHabitId = source["LastHabitId"];
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
	        this.Unit = source["Unit"];
	        this.RestDay = source["RestDay"];
	        this.RestDayMode = source["RestDayMode"];
	        this.StartAt = this.convertValues(source["StartAt"], null);
	        this.ArchivedAt = this.convertValues(source["ArchivedAt"], null);
	        this.PreviousHabit = this.convertValues(source["PreviousHabit"], Habit);
	        this.NextHabit = this.convertValues(source["NextHabit"], Habit);
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
	export class UpdateHabitDTO {
	    ID: number;
	    Name?: string;
	    Amount?: number;
	    Unit?: string;
	    RestDay?: number;
	    RestDayMode?: RestDayModeEnum;
	
	    static createFrom(source: any = {}) {
	        return new UpdateHabitDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Amount = source["Amount"];
	        this.Unit = source["Unit"];
	        this.RestDay = source["RestDay"];
	        this.RestDayMode = source["RestDayMode"];
	    }
	}

}

