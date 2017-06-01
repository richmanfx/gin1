// Таблица
jQuery(document).ready(function($) {
    $("#hamTable").tablesorter({
        sortList:[[1,0]],               // Первоначальная сортировка - по позывным
        widgets:['zebra'],              // Серые строки через одну
        headers:{
            3:{sorter:false},
            7:{sorter:false},
            8:{sorter:false},
            9:{sorter:false},
            10:{sorter:false},
            11:{sorter:false},
            12:{sorter:false},
            13:{sorter:false}
        }
    });
});

// Если поле 'my_qra_cup' формы 'cup_form' пустое, алерт и по ссылкам не переходить
$(function(){
    $(".russian-cup").click(function() {
        if(document.cup_form.my_qra_cup.value.length === 0) {
            alert("Нужно ввести свой QRA!");
            return false;       //предотвращаем переход по ссылке href или data-target
        }
    });
});

// Если поле 'my_qra_fd' формы 'fd_form' пустое, алерт и по ссылкам не переходить
$(function(){
    $(".field-day").click(function() {
        if(document.fd_form.my_qra_fd.value.length === 0) {
            alert("Нужно ввести свой QRA!");
            return false;       //предотвращаем переход по ссылке href или data-target
        }
    });
});