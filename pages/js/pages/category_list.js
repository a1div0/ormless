const desktop_view = {
	"rows": [
		{
			"css": "webix_dark",
			"view": "toolbar",
			"cols": [
				{ "layout": "x", "type": { "subsign": false },
					"view": "menu",
					"height": 0,
					"data": [
						{
							"value": "Меню",
							"submenu": [
								{ "value": "Категории", href:"category_list.html" },
								{ "value": "Бюджеты" },
								{ "value": "Движения" },
								{ "value": "Отчёт" }
							]
						}
					]
				},
				{ "view": "label", "label": "Категории", "align": "center", "gravity": 1 },
				{ "view": "label", "height": 0, "borderless": 1 }
			]
		},
		{
			"url": "cmd/entity.category_list",
			"scrollX": false,
            //autoConfig:true,
			"columns": [
				{ "id": "category_id"   , "header": "ИД", "width":50 },
                {
                    "id": "category_name" ,
                    "header": "Наименование",
                    "fillspace": true,
                    "template": "{common.treetable()} #category_name#"
                },
                { "id": "img_url"       , "header": "Значок", "fillspace": true },
                { "id": "inout_sign"    , "header": "Знак операции"},
                { "id": "sort"          , "header": "Сортировка"},
                { "id": "visible"       , "header": "Видимость"}
			],
			"view": "treetable"
		}
	]
};

window.onload = function() {
	webix.ready(function(){
		webix.ui(desktop_view);
	});
}
