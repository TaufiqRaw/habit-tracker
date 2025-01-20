import React, { createContext, useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';

type FlashStatus = "success" | "warning" | "error"
type FlashContextType = {
    flashMsg(msg : string, status : FlashStatus) : void,
}
export const FlashContext = createContext<FlashContextType>({
    flashMsg : (m)=>{}
})

const MSG_TIMEOUT = 3000
export function AppLayout() {
    const [flashMsg, setFlashMsg] = useState<{
        msg : string, status : FlashStatus
    } | null>(null)
    const [clearMsgCb, setClearMsgCb] = useState<number | null>(null) 

    useEffect(()=>{
        if(flashMsg === null)
            return
        else {
            if(clearMsgCb !== null){
                clearTimeout(clearMsgCb)
            }
            setClearMsgCb(setTimeout(()=>{
                setFlashMsg(null)
                setClearMsgCb(null)
            }, MSG_TIMEOUT))
        }
    }, [flashMsg])

    return (
        <div id='App' className='p-2 flex flex-col gap-2 relative'>
            { flashMsg && 
                <div className='absolute top-0 pt-2 flex flex-col items-center w-[99vw]'>
                    <div className={
                        "relative z-50 w-32 p-2 " + 
                        (flashMsg.status === "success" 
                            ? "bg-green-200 border-green-500 border"
                            : flashMsg.status === "warning"
                                ? "bg-yellow-200 border-yellow-500 border"
                                : "bg-red-200 border-red-500 border")
                    }>
                        {flashMsg.msg}
                    </div>
                </div>
            }
            <FlashContext.Provider value={{
                flashMsg : (msg, status)=>{
                    setFlashMsg({msg,status})
                }
            }}>
                <Outlet />
            </FlashContext.Provider>
        </div>
    );
}