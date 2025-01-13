export function saveUserId(id){
    localStorage.setItem('user-id', id)
}

export function getUserId(){
    return localStorage.getItem('user-id') 
}

export function clearUserId(){
    localStorage.removeItem('user-id') 
}

 