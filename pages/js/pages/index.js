"use strict";

window.onload = function(){
    convert_markdown_for_page();
    auth_dialog_init();

    cookieLoad();
    let user_name = cookieGet('user_name');
    if ((typeof(user_name) == 'string') && (user_name.length > 0)) {
        window.location.href = THIS_APP_URL + '/desktop.html';
    }
}

function convert_markdown_for_page() {

    let converter = new showdown.Converter();
    //converter.setFlavor('github');
    const elements = document.getElementsByClassName('markdown');
    for (let index = 0; index < elements.length; ++index) {
        let e = elements[index];
        let md = e.innerHTML;
        md = md.replace(/    /g,'');
        e.innerHTML = converter.makeHtml(md);
    }

}

function auth_dialog_init() {

    var dialog = document.querySelector('dialog');
    if (! dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }
    var btn_login = document.querySelector('#btn_login');
    btn_login.addEventListener('click', function() {
      dialog.showModal();
    });
    dialog.querySelector('.close').addEventListener('click', function() {
      dialog.close();
    });
    var btn_yandex = document.querySelector('#btn_yandex');
    btn_yandex.addEventListener('click', function(){
        window.location.href = OAUTH_YANDEX_URL;
    });
    var btn_google = document.querySelector('#btn_google');
    btn_google.addEventListener('click', function(){
        window.location.href = OAUTH_GOOGLE_URL;
    });
    var btn_github = document.querySelector('#btn_github');
    btn_github.addEventListener('click', function(){
        window.location.href = OAUTH_GITHUB_URL;
    });

}
