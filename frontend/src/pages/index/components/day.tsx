import { CalendarDay } from "../../../classes/calendar_day";

export function Day(props : {
    calenderDay : CalendarDay
}) {
    return <div className={"w-12 h-12 flex items-center justify-center " 
        + (props.calenderDay.activeMonth 
            ? (props.calenderDay.activeDate 
                ? " bg-blue-500 text-white"
                : "bg-slate-200 ") 
            : "bg-slate-600 text-slate-300 ")
        }>
        <div>
            {props.calenderDay.date}
        </div>
    </div>
}