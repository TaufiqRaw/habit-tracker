// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {domain} from '../models';

export function Create(arg1:domain.CreateHabitDTO):Promise<domain.Habit|any>;

export function Delete(arg1:number):Promise<any>;

export function GetNode(arg1:number):Promise<domain.HabitNode|any>;

export function Index(arg1:number,arg2:number,arg3:boolean):Promise<Array<domain.Habit>|any>;

export function ToggleArchived(arg1:number):Promise<domain.Habit|any>;

export function Update(arg1:domain.UpdateHabitDTO):Promise<domain.Habit|any>;
