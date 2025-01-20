import { CalendarDay, days } from "../../../classes/calendar_day";
import { Day } from "./day";

export function Calendar(props : {
    calendarDays : CalendarDay[]
}) {
    return <div className="grid grid-cols-7 gap-2 w-max p-2 pt-0 border-gray-600 border place-self-center">
    { days.map((d, i)=><div className='w-12 h-12 flex items-center justify-center'>
        {d.slice(0, 3)}
    </div>)}
    { props.calendarDays.map((cd, i)=><Day calenderDay={cd}/>)}
    
</div>
}