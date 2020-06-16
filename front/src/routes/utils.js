// TODO: list async fetch

// const basepath = "https://127.0.0.1:8081/"
const basepath = "https://api.solvemytest.dev/"

async function myfetch(user, r, method, body) {
    try {
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
        else {
            data = await res.text();
            if (data==null) {
                data = ""
            }
            alert(data)
        }
        return {data: data, status: res.status}
    }
    catch (e) {
        alert(`problem fetching ${apipath(r)}`)
        return {data: [], status:-1}
    }
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

export function loginpath() {
    return basepath+'auth/github'
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

function check_cookie_name(name) {
    var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    if (match) {
        return match[2];
    }
    else{
        return null;
    }
}

function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

export function getUser() {
    let jwt_cookie = check_cookie_name("jwt_auth")
    if (jwt_cookie == null) {
        return {}
    }
    let user = parseJwt(jwt_cookie)
    user["jwt_auth"] = jwt_cookie
    return user
}
