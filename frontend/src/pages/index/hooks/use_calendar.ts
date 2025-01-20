import { useMemo } from "react";
import { CalendarDay } from "../../../classes/calendar_day";

export function useCalendar() : CalendarDay[] {
    const now = useMemo(()=>new Date(Date.now()), [])
    const calendar = useMemo(()=>{
        let c = [] as CalendarDay[]

        // used to iterate one month*
        // * : from sunday to saturday (5 weeks)
        let d = new Date(now.getFullYear(), now.getMonth());
        // set initial date to sunday
        if(d.getDay() !== 0){
            d.setDate(d.getDate() - d.getDay())
        }
        //prev month and this month iteration
        while(d.getMonth() != ((now.getMonth() + 1) % 12)) {
            c.push(new CalendarDay(d, now.getMonth(), now.getDate()))
            d.setDate(d.getDate() + 1);
        }
        //next month iteration
        for(let i = 0; i<7 && d.getDate() != 0 ; i++){
            c.push(new CalendarDay(d, now.getMonth(), now.getDate()))
            if(d.getDay() == 6){
                break
            }
        }
        return c
    }, [now])
    return calendar
}