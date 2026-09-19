export const API='http://127.0.0.1:8080'
export const AUTH_API='http://127.0.0.1:8082'
export function token(){return localStorage.getItem('auth_token')||''}
export function userId(){return Number(localStorage.getItem('user_id')||0)}
export function logout(){localStorage.removeItem('auth_token');localStorage.removeItem('user_id');localStorage.removeItem('auth_user')}
export async function request(path:string,init:RequestInit={}){
 const headers=new Headers(init.headers||{});const t=token();if(t)headers.set('Authorization','Bearer '+t)
 const r=await fetch(API+path,{...init,headers});const j=await r.json().catch(()=>({}));if(!r.ok)throw new Error(j.message||j.error||'请求失败');return j
}
export async function authRequest(path:string,init:RequestInit={}){
 const headers=new Headers(init.headers||{});const t=token();if(t)headers.set('Authorization','Bearer '+t)
 const r=await fetch(AUTH_API+path,{...init,headers});const j=await r.json().catch(()=>({}));if(!r.ok)throw new Error(j.message||j.error||'请求失败');return j
}