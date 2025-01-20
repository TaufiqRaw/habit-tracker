import { HTMLProps, useEffect, useRef, useState } from "react"

export function TDHeight(props : {
    children : (h : number)=>React.ReactElement
} & Omit<HTMLProps<HTMLTableCellElement>, "children">){
    const tdRef = useRef<HTMLTableCellElement>(null) 
    const [height, setHeight] = useState(0)

    useEffect(()=>{
        if(!tdRef.current) return
        setHeight(tdRef.current.clientHeight)
    }, [])

    return <td ref={tdRef} {...props}>
        {props.children(height)}
    </td>
}