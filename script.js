var Component = {
    Email: document.querySelector('#email'),
    Password: document.querySelector('#password'),
    Message: document.querySelector('#message'),
    Form: document.querySelector('#myForm'),
    CloseButton: document.querySelector('#close'),
    Modal: document.querySelector('#myModal'),
    LastTimeclock: document.querySelector('#lastTimeclock')
};

var storage = {
    setItem: function (key, value) {
        var d = new Date();
        d.setTime(d.getTime() + (365*24*60*60*1000));
        var expires = 'expires=' + d.toUTCString();
        document.cookie = key + '=' + value + '; ' + expires + ';path=/';
    },
    getItem: function (key) {
        var c = document.cookie;
	var data = c.match(new RegExp(key + '=([^;]+)'));
	if (data) {
	    return data[1].split(', ')[1];
	}
	return '';
    }
};

window.onload = function () {
    Component.Email.focus();
    if (storage.getItem('lastTimeclock')) {
        showLastTimeclock();
    }
}

/*
window.onclick = function (event) {
    if (event.target === Component.Modal) {
        closeModal(); //Component.Modal.style.display = 'none';
    }
}
*/

function submitForm () {
    var email = Component.Email.value.trim() || '';
    var password = Component.Password.value.trim() || '';

    if (email.length > 0 && password.length > 0) {
        registerTimeclock({email: email, password: password});
    }
    else {
        showModal('Informe email e senha.');
    }
    Component.Email.focus();
    return false;
}

function registerTimeclock (credentials) {
    showModal('Enviando... Aguarde.');
    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'https://pontomenos.herokuapp.com/', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (this.status === 200) {
	    //var response = JSON.parse(this.responseText);
            //var result = response.message.replace(/(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})(\d{2})/, '$1-$2-$3T$4:45:$6');
	    var result = new Date().toLocaleString();
	    showModal('Ponto registrado.');
	    storage.setItem('lastTimeclock', result);
            showLastTimeclock();
	}
	else {
	    showModal('Erro ' + this.status + ': Falhou.');
	}
    };
    xhr.send(JSON.stringify(credentials));
}

function showModal (message) {
    Component.Message.innerHTML = message;
    Component.Modal.style.display = 'block';
}

function closeModal () {
    Component.Modal.style.display = 'none';
    Component.Email.focus();
}

function showLastTimeclock () {
    /*
    document.querySelector('footer').style.display = 'block';
    Component.LastTimeclock.innerHTML = '&Uacute;ltimo ponto: ' + storage.getItem('lastTimeclock');
    Component.LastTimeclock.style.display = 'block';
    */
}

