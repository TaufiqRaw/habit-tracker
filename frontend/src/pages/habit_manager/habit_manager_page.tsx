import { HabitPanel } from './components/habit_panel'
import PlusLogo from "../../assets/images/add.svg?react"
import { useState } from 'react'
import { HabitAddForm } from './components/habit_add_form';

export const HABIT_MANAGER_PAGE_PATH = "/habits" 

export function HabitManagerPage(){
    const [addMode, setAddMode] = useState(false);
    
    function onCancelAdd(){
        setAddMode(false)
    }
    function onAddClick(){
        setAddMode(true)
    }
    return <>
        <div className='grow p-2 relative border-gray-600 border bg-slate-100'>
            <div className='w-max font-black font-inter text-2xl border-r border-b border-gray-600 px-4 pt-2 absolute left-0 top-0 bg-white'>
                <div className='h-12 flex justify-center'>
                    <div>
                        ALL HABIT
                    </div>
                </div>
            </div>
            <div className='absolute right-0 top-0 p-2'>
                <button className='bg-blue-500 p-2 hover:bg-blue-400 active:bg-blue-500' onClick={onAddClick}>
                    <PlusLogo className='fill-white'/>
                </button>
            </div>
            <div className='flex flex-col mt-14 gap-3'>
                { addMode && <HabitAddForm onCancel={onCancelAdd}/>}
                <HabitPanel/>
            </div>
        </div>
    </>
}