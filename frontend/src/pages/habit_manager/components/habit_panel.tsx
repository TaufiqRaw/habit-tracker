import { HabitPanelItem } from "./habit_panel_item"
import DblChevronRight from '../../../assets/images/double-chevron-right.svg?react'
import { HabitEditForm } from './habit_edit_form'

export function HabitPanel(){
    return <div className='flex flex-col gap-1 p-2 bg-slate-50 border-gray-300 border relative pt-8'>
        <Habit/>
    </div>
}

function Habit(){
    return <>
        <div className='w-max border-r border-b border-gray-300 px-4 absolute left-0 top-0 bg-white'>
            <div className=' flex justify-center'>
                <div>
                    Gaming
                </div>
            </div>
        </div>
        <div className="flex gap-2 items-center">
            <HabitPanelItem/>
            <DblChevronRight  className='fill-gray-300'/>
            <HabitEditForm/>
        </div>
    </>
}