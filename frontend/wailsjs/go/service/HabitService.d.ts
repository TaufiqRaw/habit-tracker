// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {domain} from '../models';
import {service} from '../models';

export function Create(arg1:domain.CreateHabitDTO):Promise<service.habitResult>;

export function Delete(arg1:number):Promise<any>;

export function Index(arg1:number,arg2:number,arg3:boolean):Promise<service.habitArrResult>;

export function ToggleArchived(arg1:number):Promise<service.habitNodeResult>;

export function Update(arg1:number,arg2:string):Promise<any>;
