import { Calendar } from './components/calendar';
import { MonthInfo } from './components/month_info';
import { TrackHabit } from './components/track_habit';
import ChevronRight from '../../assets/images/chevron-right.svg?react'
import { useCalendar } from './hooks/use_calendar';
import { useNavigate } from 'react-router-dom';
import { HABIT_MANAGER_PAGE_PATH } from '../habit_manager/habit_manager_page';

export const INDEX_PAGE_PATH = "/" 

export function IndexPage() {
    const calendar = useCalendar()
    const nav = useNavigate()
    function toHabitManagerPage(){
        nav(HABIT_MANAGER_PAGE_PATH)
    }
    return <>
        <div className='flex gap-2'>
            <Calendar calendarDays={calendar}/>
            <MonthInfo/>
        </div>
        <TrackHabit/>
        <div className='font-black font-inter text-2xl border border-gray-600 px-4  bg-white hover:bg-blue-200 cursor-pointer group' onClick={toHabitManagerPage}>
            <div className='h-12 flex justify-center select-none'>
                <div className='flex items-center gap-3'>
                    <span>HABIT MANAGER</span>
                    <ChevronRight className='transform transition-transform group-hover:-translate-x-2'/>
                </div>
            </div>
        </div>
    </>
}