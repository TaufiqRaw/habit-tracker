export const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
type Day = typeof days[number];

export class CalendarDay {
    day : Day
    date : number
    activeMonth : boolean
    activeDate : boolean

    constructor(
        date : Date,
        currentMonth : number,
        currentDate : number,
    ){
        this.day = days[date.getDay()]
        this.date = date.getDate()
        this.activeMonth = currentMonth === date.getMonth()
        this.activeDate = this.activeMonth && (currentDate == this.date)
    }
}