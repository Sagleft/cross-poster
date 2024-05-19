function showAlert(info) {
    new Noty({
        type: 'alert',
        layout: 'topRight',
        text: info,
        timeout: 4500,
    }).show();
}

function showError(info) {
    new Noty({
        type: 'error',
        layout: 'topRight',
        text: info,
        timeout: 10000,
    }).show();
}

function IsJsonString(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

function shutdownapp() {
    if(confirm("End the application?")) {
        $.get("/exit").always(function() {
            window.close();
        });
    }
}

var bar = document.getElementById('js-progressbar');

$( document ).ready(function() {
    UIkit.upload('.js-upload', {

        url: '/upload',
        multiple: false,
    
        beforeSend: function () {
            console.log('beforeSend', arguments);
        },
        beforeAll: function () {
            console.log('beforeAll', arguments);
        },
        load: function () {
            console.log('load', arguments);
        },
        error: function () {
            console.log('error', arguments);
        },
        complete: function () {
            console.log('complete', arguments);
        },
    
        loadStart: function (e) {
            console.log('loadStart', arguments);
    
            bar.removeAttribute('hidden');
            bar.max = e.total;
            bar.value = e.loaded;
        },
    
        progress: function (e) {
            console.log('progress', arguments);
    
            bar.max = e.total;
            bar.value = e.loaded;
        },
    
        loadEnd: function (e) {
            console.log('loadEnd', arguments);
    
            bar.max = e.total;
            bar.value = e.loaded;
        },
    
        completeAll: function () {
            console.log('completeAll', arguments);
    
            setTimeout(function () {
                bar.setAttribute('hidden', 'hidden');
            }, 1000);

            $("#isImageUploaded").val('1');
        }
    
    });

    new ClipboardJS('.btn-copy');

    $( "#main-form" ).submit(function( event ) {
        $.post( "/send", $( "#main-form" ).serialize(), function(response) {
            if(response.status == "error") {
                showError(response.error);
            } else {
                showAlert("The post was published");
                $("#isImageUploaded").val('0');
            }
        });
        event.preventDefault();
    });

    setTimeout( function() { checkAppStatus(); }, 2000);
});

function checkAppStatus() {
    $.post( "/check", {}, function(response) {
        if(response.status == "error") {
            showError(response.error);
        } else {
            showAlert("connected!");
        }
    });
}
