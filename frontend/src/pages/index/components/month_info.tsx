export function MonthInfo() {
    return <div className='grow p-2 relative border-gray-600 border bg-slate-100'>
        <div className='w-max font-black font-inter text-2xl border-r border-b border-gray-600 px-4 pt-2 absolute left-0 top-0 bg-white'>
            <div className='h-12 flex justify-center'>
                <div>
                    JANUARY
                </div>
            </div>
        </div>
        <div className='flex flex-col mt-14 gap-3'>
            <div className='flex flex-col'>
                <div>Gaming</div>
                <div className='flex'>
                    <div className='grow mt-2'>
                        <div className='h-2 flex rounded-full overflow-hidden bg-white border border-gray-600'>
                            <div className='bg-blue-500 w-1/2'></div>
                        </div>
                    </div>
                    <div className='ml-2'>
                        4/5
                    </div>
                </div>
            </div>
        </div>
    </div>
}