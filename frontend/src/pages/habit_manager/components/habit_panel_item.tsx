export function HabitPanelItem(){
    return <div className="bg-slate-100 border-gray-600 border p-2 pt-3 relative flex flex-col gap-2 w-56 group hover:outline outline-2 outline-blue-500">
        <div className="relative">
            <div className="absolute -top-3 -left-2 w-[calc(100%+(0.5rem*2))] bg-slate-500 h-2"/>
            <table>
                <tbody>
                    <tr>
                        <td>Min per Day</td>
                        <td>:</td>
                        <td>2 Hour(s)</td>
                    </tr>
                    <tr>
                        <td>Max Rest Day</td>
                        <td>:</td>
                        <td>2 per Month</td>
                    </tr>
                </tbody>
            </table>
            <div className='hidden group-hover:flex absolute top-full bg-slate-100 p-2 w-full border border-gray-600 flex-col gap-2 shadow-md'>
                <button className='bg-blue-500 w-full hover:bg-blue-600'>
                    <div className='text-white px-2'>Edit</div>
                </button>
                <button className='bg-slate-500 w-full hover:bg-slate-600'>
                    <div className='text-white px-2'>Archive</div>
                </button>
                <button className='bg-green-500 w-full hover:bg-green-600'>
                    <div className='text-white px-2'>Unarchive</div>
                </button>
                <button className='bg-red-500 w-full hover:bg-red-600'>
                    <div className='text-white px-2'>Delete</div>
                </button>
            </div>
        </div>
    </div>
}