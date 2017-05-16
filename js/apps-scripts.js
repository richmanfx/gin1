// Таблица
jQuery(document).ready(function($) {
    $("#hamTable").tablesorter({
        sortList:[[1,0]],            // Первоначальная сортировка - по позывным
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
