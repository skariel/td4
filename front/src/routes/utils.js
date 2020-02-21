const basepath = "https://localhost:8081/"

export async function post(user, r, o) {
    const res = await fetch(basepath+'api/'+r, {
        method: 'POST',
        body: JSON.stringify(o),
        headers: new Headers({
            Authorization: 'Bearer '+user.jwt_auth,
        })
    });
    if (res.status == 200) {
        const data = await res.json();
        return {data: data, status: res.status}        
    }
    else {
        return {status: res.status}        
    }
}

export async function get(user, r) {
    const res = await fetch(basepath+'api/'+r);
    if (res.status == 200) {
        const data = await res.json();
        return {data: data, status: res.status}        
    }
    else {
        return {status: res.status}        
    }
}

export function staticpath(p) {
    return basepath+'static/'+p
}