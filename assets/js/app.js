function showAlert(info) {
    new Noty({
        type: 'alert',
        layout: 'topRight',
        text: info,
        timeout: 1500,
    }).show();
}

function showError(info) {
    new Noty({
        type: 'error',
        layout: 'topRight',
        text: info,
        timeout: 1500,
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

            $("isImageUploaded").value('1');
        }
    
    });

    new ClipboardJS('.btn-copy');

    $( "#main-form" ).submit(function( event ) {
        $.post( "/send", $( "#main-form" ).serialize(), function(data) {
            if(data.status == "error") {
                showError(data.error);
            } else {
                showAlert("The post was published");
                $("isImageUploaded").value('0');
            }
        });
        event.preventDefault();
    });
});
