"use strict";

function get_radio_value(radio_name) {

    var radios = document.getElementsByName(radio_name);

    for (var i = 0, length = radios.length; i < length; i++)
    {
        if (radios[i].checked)
        {
            return radios[i].value;
        }
    }

    return null;

}

function get_radio_value_to_id(radio_name, value) {

    var radios = document.getElementsByName(radio_name);

    for (var i = 0, length = radios.length; i < length; i++)
    {
        if (radios[i].value == value)
        {
            return radios[i].id;
        }
    }

    return null;

}

function trim(str) {
    return str.trim();
}

function html_post_json(url,json_str,onrecieve)
{
    var h = getXmlHttp();
    h.open('POST', url, true);
    h.setRequestHeader('Content-Type','application/json');
    h.onreadystatechange = function() {

        if(h.readyState != 4) {
            return;
        }

        if(h.status == 200)
        {
            var res = JSON.parse(h.responseText);
            if (res.status == 'error') {
                alert(res.message);
            }else if (res.status == 'warning') {
                alert(res.message);
            }else if (res.status != 'success') {
                alert('Неизвестный статус ответа: ' + res.status);
            }else{
                if (onrecieve != null) {
                    onrecieve(res.data);
                }
            }
        }
        else
        {
            alert('Ошибка HTTP-запроса! Код ошибки: ' + h.status + '. ' + h.responseText);
        }

    }

    h.send( json_str );
}

function getXmlHttp(){
    var xmlhttp;

    if (typeof XMLHttpRequest != 'undefined') {
        xmlhttp = new XMLHttpRequest();
    }else{
        xmlhttp = false;
    }

    return xmlhttp;
}
