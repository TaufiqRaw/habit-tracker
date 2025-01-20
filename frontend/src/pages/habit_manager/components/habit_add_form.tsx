import { ChangeEvent, useContext, useMemo, useState } from "react";
import { Input, inputGuardInt } from "../../../components/input";
import CancelLogo from "../../../assets/images/cancel.svg?react";
import SendLogo from "../../../assets/images/send.svg?react";
import { TDHeight } from "../../../components/td_height";
import {domain} from "../../../../wailsjs/go/models"
import * as HabitService from "../../../../wailsjs/go/service/HabitService"
import { FlashContext } from "../../../components/app_layout";

const actPlaceholders = [
    ["Read a book", 7, "Page(s)"],
    ["Working out", 45, "Minute(s)"],
    ["Network with", 2, "People(s)"]
] as const

export function HabitAddForm(props : {
    onCancel : ()=>void
}){
    const flashContext = useContext(FlashContext)

    const [name, setName] =  useState("")
    const [minPerDay, setMinPerDay] = useState("")
    const [unit, setUnit] = useState("")
    const [restDay, setRestDay] = useState("")
    const [restDayMode, setRestDayMode] = useState<domain.RestDayModeEnum>(Object.values(domain.RestDayModeEnum)[0])

    const placeholderIndex = useMemo(()=>Math.floor(Math.random() * actPlaceholders.length), [])

    function onRestDayModeChange(e : ChangeEvent<HTMLSelectElement>) {
        const v = e.target.value as domain.RestDayModeEnum

        if(e.target.value != "" 
            && Object.values(domain.RestDayModeEnum).indexOf(v) < 0){
            return
        }

        setRestDayMode(v)
    }
    async function onSubmit() {
        const res = await HabitService.Create(new domain.CreateHabitDTO({
            Name : name,
            Amount : parseInt(minPerDay),
            RestDay : parseInt(restDay),
            RestDayMode : restDayMode,
            Unit : unit
        } as domain.CreateHabitDTO));
        if(res.Err){
            return flashContext.flashMsg(res.Err, "error")
        }else {
            console.log(JSON.stringify(res.Data))
            return flashContext.flashMsg("Success", "success")
        }
    }
    return <div className="bg-slate-50 border-gray-300 border p-2 pt-3 relative flex flex-col gap-2 max-w-2xl">
    <div className="absolute right-0 bottom-full -mb-1 py-2 bg-white border border-gray-300">
        <div className="flex items-stretch justify-center w-20 h-10">
            <button className="bg-blue-500 active:bg-blue-500 hover:bg-blue-400 text-white grid place-content-center px-1" onClick={onSubmit}>
                <SendLogo className="fill-white"/>
            </button>
            <button className="border-2 border-blue-500 text-blue-500 hover:border-blue-400 hover:bg-blue-50 active:bg-blue-100 grid place-content-center px-1" onClick={props.onCancel}>
                <CancelLogo className="fill-blue-500"/>
            </button>
        </div>
    </div>
    <div>
        <table className="w-full">
            <tbody>
                <tr>
                    <td className="w-3/12">Name</td>
                    <td>:</td>
                    <td>
                        <Input placeholder={actPlaceholders[placeholderIndex][0]} value={name} onChange={setName}/>
                    </td>
                </tr>
                <tr>
                    <td>Min each day</td>
                    <td>:</td>
                    <td>
                        <div className='flex gap-1'>
                            <div className='min-w-4'>
                                <Input placeholder={actPlaceholders[placeholderIndex][1]+ ""} value={minPerDay} onChange={setMinPerDay} guard={inputGuardInt(0, 100)}/>
                            </div>
                            <div className='grow'>
                                <Input placeholder={actPlaceholders[placeholderIndex][2]} value={unit} onChange={setUnit}/>
                            </div>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>Max rest day</td>
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