import { ChangeEvent, useEffect, useRef, useState } from "react";

// function insertTextAtSelection(div : HTMLDivElement, txt : string) {
//     //get selection area so we can position insert
//     let sel = window.getSelection();
//     let text = div.textContent;
//     if(!text || !sel) return
//     let before = Math.min(sel.focusOffset, sel.anchorOffset);
//     let after = Math.max(sel.focusOffset, sel.anchorOffset);
//     //ensure string ends with \n so it displays properly
//     let afterStr = text.substring(after);
//     if (afterStr == "") afterStr = "\n";
//     //insert content
//     div.textContent = text.substring(0, before) + txt + afterStr;
//     //restore cursor at correct position
//     sel.removeAllRanges();
//     let range = document.createRange();
//     //childNodes[0] should be all the text
//     range.setStart(div.childNodes[0], before + txt.length);
//     range.setEnd(div.childNodes[0], before + txt.length);
//     sel.addRange(range);
// }

//TODO: if overflow, make left text hidden instead of growing

type GuardType = { 
    keyDown ?: (key : string, value : string)=>boolean, 
    keyUp ?: (value : string)=>string
} 

export function inputGuardInt(min ?: number, max ?: number) : GuardType {
    if(min && max) {
        if (min > max) throw Error("inputGuardInt : Invalid args")
    }
    return {
        keyDown(k : string) : boolean {
            const n = parseInt(k)
            if(Number.isNaN(n)) return false
            else return true
        },
        keyUp(v) : string {
            let n = parseInt(v)
            if(Number.isNaN(n)) return ""
            else {
                if(min !== undefined) {
                    n = Math.max(min, n)
                }
                if(max !== undefined) {
                    n = Math.min(max, n)
                }
                return n + ""
            }
        }
    }
}

export function Input(props :{
    value : string,
    placeholder : string,
    onChange : (s : string)=>void
    guard ?: GuardType
}) {
    const [isFocus, setIsFocus] = useState(false)
    const divRef = useRef<HTMLDivElement>(null)

    function getCaretPosition() {
        var caretPos = 0,
            sel : Selection | null, range : Range;
        if (divRef.current) {
            sel = window.getSelection();
            if (sel?.rangeCount) {
                range = sel.getRangeAt(0);
                if (range.commonAncestorContainer.parentNode == divRef.current) {
                    caretPos = range.endOffset;
                }
            }
        } 
        return caretPos;
    }

    function setCaret(pos : number) {
        var range = document.createRange()
        var sel = window.getSelection()

        if(!divRef.current || !sel || !divRef.current.childNodes[0]) return
        
        const maxLen = (divRef.current.textContent?.length ?? 0)
        if(pos > maxLen){
            range.setStart(divRef.current.childNodes[0], maxLen)
        }else{
            range.setStart(divRef.current.childNodes[0], pos)
        }
        range.collapse(true)
        
        sel.removeAllRanges()
        sel.addRange(range)
    }

    useEffect(()=>{
        if(!divRef.current) return
        const caretPos = getCaretPosition()
        divRef.current.innerText = props.value
        setCaret(caretPos)
    }, [props.value])

    function pasteGuard(e : React.ClipboardEvent<HTMLDivElement> & {originalEvent : ClipboardEvent} ){
        //TODO:allow pasting
        //Pasting disabled for now
        e.preventDefault();
        // let text = (e.originalEvent || e).clipboardData?.getData('text/plain');
        // if(divRef.current && text){
        //     insertTextAtSelection(divRef.current, text);
        //     props.onChange(divRef.current.innerText)
        // }
    }
    function keydownGuard(e : React.KeyboardEvent<HTMLDivElement>) {
        //TODO:check and prevent more irrelevant or bug inducing chars (if any)
        switch(e.key) {
            case "Enter" :
                e.preventDefault();
                e.stopPropagation();
                break;
        }
        if(props.guard && props.guard.keyDown) {
            const alwaysAccepted = ["Backspace", "Delete", "ArrowLeft", "ArrowRight"]
            if(!props.guard.keyDown(e.key, props.value) && !alwaysAccepted.some(v=>v===e.key)){
                e.preventDefault();
                e.stopPropagation();
            }
        }
    }
    function keyUpCapture(e : React.KeyboardEvent<HTMLDivElement>) {
        if(!divRef.current) return
        const text = divRef.current.innerText
        if(props.guard && props.guard.keyUp) {
            const prev = props.value
            const newVal = props.guard.keyUp(text)
            if(prev !== newVal) props.onChange(newVal)
            else {
                const caretPos = getCaretPosition()
                divRef.current.innerText = prev
                setCaret(caretPos)
            }
        }else {
            props.onChange(text)
        }
    }
    function onFocusClick(){
        setIsFocus(true)
        if(!divRef.current) setTimeout(onFocusClick, 20)
        else {
            divRef.current.innerText = props.value
            divRef.current!.focus()
        }
    }
    function onFocus(){
        setIsFocus(true)
    }
    function onBlur(){
        setIsFocus(false)
    }
    return <div className='border border-gray-600 px-1 w-full cursor-text'>
        { isFocus || props.value.trim() != ""
            ? <div className="text-gray-900" ref={divRef} onPaste={pasteGuard} onKeyDown={keydownGuard} onFocus={onFocus} onKeyUp={keyUpCapture} onBlur={onBlur} contentEditable suppressContentEditableWarning={true}/>
            : <div className="text-gray-400" onClick={onFocusClick}>
                {props.placeholder}
            </div>
        }
    </div>
}