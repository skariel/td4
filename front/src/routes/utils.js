// TODO: list async fetch

const basepath = "https://127.0.0.1:8081/"

async function myfetch(r, method, body) {
    const res = await fetch(apipath(r), {
        method: method,
        body: body,
    });
    var data = []
    if (res.status == 200) {
        data = await res.json();
        if (data==null) {
            data = []
        }
    }
    return {data: data, status: res.status}        
}

export async function post(r, o) {
    const method = 'POST'
    const body = JSON.stringify(o)
    return myfetch(r, method, body)
}

export async function get(r) {
    const method = 'GET'
    const body = null
    return myfetch(r, method, body)
}

function apipath(p) {
    return basepath+'api/'+p
}

export function staticpath(p) {
    return basepath+'static/'+p
}