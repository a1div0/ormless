"use strict";

window.onload = function(){
    auth_dialog_init();
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
