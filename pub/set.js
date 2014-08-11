CKEDITOR.replace( 'editor', {
	uiColor: '#FAFAFA',
	allowedContent: true,
	height: 380,
	toolbar: [
		['Undo', 'Redo'],
		['Styles','Format'],
		['Bold', 'Italic', 'RemoveFormat'],
		['JustifyCenter', 'BulletedList', 'Blockquote'],
		['TextColor', 'BGColor'],
		['Source']
	],
	format_tags: 'p;pre;h3;h4'
});

var form = document.forms.namedItem("fileinfo");

form.addEventListener('submit', function(ev) {
	var oData = new FormData(document.forms.namedItem("fileinfo"));
	var progress = document.getElementById('uploadprogress');
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "upload", true);
	xhr.upload.onprogress = function (event) {
		if (event.lengthComputable) {
			var complete = (event.loaded / event.total * 100 | 0);
			progress.value = progress.innerHTML = complete;
		}
	};
	xhr.onload = function(oEvent) {
		if (xhr.status == 200) {
			var c = xhr.responseText;
//			var i = c.indexOf("<");
//			CKEDITOR.instances.editor.insertHtml(c.slice(i));

// 依据文件名添加。var c 是string，不是[]。
//			out = "kxmd";
	//		for (var i in c){
	//			str = c[i];
//				console.log(str);
	//			ext = str.slice(-4).toLowerCase();
	//			if (ext == ".png" || ext == ".jpg" || ext == ".gif") {
	//				out += '<img src="/pub/dbmy/' + str + '">';
	//			}
	//		}

			CKEDITOR.instances.editor.insertHtml(c);
			document.getElementById("file").value = "";
			progress.value = 0;
		}
	};
	xhr.send(oData);
	ev.preventDefault();
}, false); 
