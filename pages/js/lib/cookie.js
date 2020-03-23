"use strict";

var cookie_db = {};

function cookieLoad() {

    let m = document.cookie.split(';');
    cookie_db = {};

    for (var i = 0; i < m.length; i++) {

        let name_value = trim(m[i]);

        if (name_value.length == 0) {continue;}

        let nv = name_value.split('=');

        switch (nv.length) {
            case 0: break;
            case 1: cookie_db[trim(nv[0])] = true; break;
            default:
                cookie_db[trim(nv[0])] = trim(nv[1]);
        }
    }

}

function cookieSave() {

    let expires_days = 365 * 100;

    let expires_dt = new Date();
    expires_dt.setTime(expires_dt.getTime() + expires_days*24*60*60*1000);

    let m = [];
    for (item_name in cookie_db) {
        if (item_name == '0') {continue;}

        document.cookie = item_name + '=' + cookie_db[item_name] + ';path=/;expires=' + expires_dt.toGMTString() + ';';
    }

}

function cookieGet(name) {
    return cookie_db[name];
}

function cookieSet(name, value) {
    cookie_db[name] = value;
}
