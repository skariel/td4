const api_basepath = "https://api.solvemytest.dev/"
const static_basepath = "https://solvemytest.dev/"
const invalidate_cache_ttl_ms = 5000

let invalidate_cache = false;

export function start_invalidate_cache() {
    invalidate_cache = true;
    setTimeout(()=>{invalidate_cache=false;}, invalidate_cache_ttl_ms);
}

async function myfetch(user, r, method, body) {
    try {
        let href = apipath(r)
        if ((invalidate_cache)&&(method=='GET')) {
            let u = new URL(href);
            u.searchParams.append('uncached','true');
            href = u.href;
            console.log('cache invalidated! '+href)
        }

        const res = await fetch(href, {
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
        console.log(e)
        alert(`problem fetching ${apipath(r)}. Error logged to console`)
        return {data: [], status:-1}
    }
}

export async function post(user, r, o) {
    const method = 'POST'
    const body = JSON.stringify(o)
    return myfetch(user, r, method, body)
}

export async function get(user, r) {
    const method = 'GET';
    const body = null;
    return myfetch(user, r, method, body);
}

export async function del(user, r) {
    const method = 'DELETE';
    const body = null;
    return myfetch(user, r, method, body);
}

function apipath(p) {
    return api_basepath+'api/'+p
}

export function loginpath() {
    return api_basepath+'auth/github'
}

export function staticpath(p) {
    return static_basepath+'/'+p
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
