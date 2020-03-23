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
							//"icon": "mdi mdi-menu",
							"submenu": [
								{ "value": "Категории", href:"category_list.html" },
								{ "value": "Бюджеты" },
								{ "value": "Движения" },
								{ "value": "Отчёт" }
							]
						}
					]
				},
				{ "view": "label", "label": "Balance", "align": "center", "gravity": 1 },
				{ "view": "label", "height": 0, "borderless": 1 }
			]
		},
		{
			"autoheight": false,
			"view": "form",
			"rows": [
				{
					"cols": [
						{ "label": "Движение", "view": "text", "height": 0 },
						{ "label": "Label", "view": "label", "height": 38, "align": "center" }
					]
				},
				{ "type": { "height": 100 },
					"template": "#rank#. #title#",
					"view": "dataview"
				}
			]
		}
	]
};

window.onload = function() {
	webix.ready(function(){
		webix.ui(desktop_view);
	});
}
