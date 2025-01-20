import { useState } from "react"
import RestLogo from '../../../assets/images/rest.svg?react';

export function TrackHabit() {
    const [check, setCheck] = useState(0)
    const [amount, setAmount] = useState(0)
    const [unit, setUnit] = useState("Hour(s)")

    function onCheckClick(){
        setCheck((p)=>(p+1)%3)
    }

    return <div className='grow p-2 relative border-gray-600 border bg-slate-100'>
        <div className='w-max font-black font-inter text-2xl border-r border-b border-gray-600 px-4 pt-2 absolute left-0 top-0 bg-white'>
            <div className='h-12 flex justify-center'>
                <div>
                    TODAY'S HABIT
                </div>
            </div>
        </div>
        <div className='flex flex-col mt-14 gap-3'>
            <div className='flex gap-3 items-center'>
                <div className="flex gap-1">
                    <div className="w-6 h-6 border border-gray-600 grid place-content-center group cursor-pointer select-none" onClick={onCheckClick}>
                        { check == 0 && <div className=" w-5 h-5 group-hover:bg-blue-200"/>}
                        { check == 1 && <div className="w-5 h-5 bg-blue-500"/> }
                        { check == 2 && <div className="w-5 h-5 bg-yellow-500 relative">
                            <div className="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 w-6">
                                <RestLogo/>
                            </div>
                        </div> }
                    </div>
                    <div>Gaming</div>
                </div>
                {check == 1 && <>
                    <div className="flex gap-1">
                        <div className="h-6 min-w-6 border border-gray-600 px-2">
                            3
                        </div>
                        <div>Hour(s)</div>
                    </div>
                    <div className="h-2 border border-gray-600 rounded w-20 bg-white overflow-hidden">
                        <div className="bg-blue-500 w-1/2 h-full"/>
                    </div>
                </>}
            </div>
        </div>
    </div>
}