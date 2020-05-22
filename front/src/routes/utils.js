// TODO: list async fetch

const basepath = "https://127.0.0.1:8081/"

async function myfetch(user, r, method, body) {
    const res = await fetch(apipath(r), {
        method: method,
        body: body,
        headers: {
            'Authorization': 'Bearer ' + user.jwt_auth,
        }
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

export async function post(user, r, o) {
    const method = 'POST'
    const body = JSON.stringify(o)
    return myfetch(user, r, method, body)
}

export async function get(user, r) {
    const method = 'GET'
    const body = null
    return myfetch(user, r, method, body)
}

function apipath(p) {
    return basepath+'api/'+p
}

export function staticpath(p) {
    return basepath+'static/'+p
}

export function init_location_change_event() {
    history.pushState = ( f => function pushState(){
        var ret = f.apply(this, arguments);
        window.dispatchEvent(new Event('pushstate'));
        window.dispatchEvent(new Event('locationchange'));
        return ret;
    })(history.pushState);

    history.replaceState = ( f => function replaceState(){
        var ret = f.apply(this, arguments);
        window.dispatchEvent(new Event('replacestate'));
        window.dispatchEvent(new Event('locationchange'));
        return ret;
    })(history.replaceState);

    window.addEventListener('popstate',()=>{
        window.dispatchEvent(new Event('locationchange'))
    });
}

export function getUser() {
    let user = {};
    let cookies = document.cookie.split(';')
    for (const ix in cookies) {
        const c = cookies[ix]
        const cs = c.split('=')
        user[cs[0].trim().replace('user_', '')] = cs[1];
    }
    return user
}
