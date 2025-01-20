import { ChangeEvent, HTMLProps, useEffect, useRef, useState } from "react";
import { Input, inputGuardInt } from "../../../components/input";
import CancelLogo from "../../../assets/images/cancel.svg?react";
import SendLogo from "../../../assets/images/send.svg?react";
import { TDHeight } from "../../../components/td_height";
import {domain} from "../../../../wailsjs/go/models"

const unitPlaceholder = [
    [7, "Page(s)"],
    [45, "Minute(s)"],
    [2, "People(s)"]
]

export function HabitEditForm(){
    const [minPerDay, setMinPerDay] = useState("")
    const [unit, setUnit] = useState("")
    const [restDay, setRestDay] = useState("")
    const [restDayMode, setRestDayMode] = useState<domain.RestDayModeEnum>(Object.values(domain.RestDayModeEnum)[0])

    function onRestDayModeChange(e : ChangeEvent<HTMLSelectElement>) {
        const v = e.target.value as domain.RestDayModeEnum

        if(e.target.value != "" 
            && Object.values(domain.RestDayModeEnum).indexOf(v) < 0){
            return
        }

        setRestDayMode(v)
    }
    return <div className="bg-slate-100 border-gray-600 border p-2 pt-3 relative flex flex-col gap-2 w-70">
    <div className="relative">
        <div className="absolute -top-3 -left-2 w-[calc(100%+(0.5rem*2))] bg-blue-500 h-2"/>
        <table>
            <tbody>
                <tr>
                    <td>Min per Day</td>
                    <td>:</td>
                    <td>
                        <div className='flex gap-1'>
                            <div className='min-w-4'>
                                <Input placeholder='5' value={minPerDay} onChange={setMinPerDay} guard={inputGuardInt(0, 100)}/>
                            </div>
                            <div className='grow'>
                                <Input placeholder='Pages' value={unit} onChange={setUnit}/>
                            </div>
                        </div>
                    </td>
                    <TDHeight rowSpan={3} className="pr-5">
                        {(h)=>
                        <div className="flex items-stretch w-10" style={{height : h+"px"}}>
                            <button className="bg-blue-500 active:bg-blue-500 hover:bg-blue-400 text-white grid place-content-center px-1">
                                <SendLogo className="fill-white"/>
                            </button>
                            <button className="border-2 border-blue-500 text-blue-500 hover:border-blue-400 hover:bg-blue-50 active:bg-blue-100 grid place-content-center px-1">
                                <CancelLogo className="fill-blue-500"/>
                            </button>
                        </div>}
                    </TDHeight>
                </tr>
                <tr>
                    <td>Max Rest Day</td>
                    <td>:</td>
                    <td>
                        <div className='flex gap-1'>
                            <div className='min-w-4'>
                                <Input placeholder='2' value={restDay} onChange={setRestDay} guard={inputGuardInt(0, 100)}/>
                            </div>
                            <div>Each</div>
                            <div className="grow">
                                <select className="border border-gray-600 h-full" value={restDayMode} onChange={onRestDayModeChange}>
                                    { Object.values(domain.RestDayModeEnum).map((v)=> <option value={v}>{v[0] + v.toLowerCase().slice(1)}</option>)}
                                </select>
                            </div>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</div>
}