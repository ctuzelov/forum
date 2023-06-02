var currUrl = new URL(window.location);
var currParams = new URLSearchParams(currUrl.search);
for (let currParam of currParams) {
    let id = currParam[1];
    var element = document.getElementById("cat" + id);
    element.classList.add("active")
}
function filter(button) {
    var form = document.createElement('form');
    var id = button.value;
    console.log(button)
    var url = new URL(window.location)
    var params = new URLSearchParams(url.search);
    var flag = false
    for (let param of params) {
        if (param[1] == id) {
            flag = true;
            continue;
        }
        form.innerHTML += '<input name="filter" value="' + param[1] + '">';
    }
    if (!flag) {
        form.innerHTML += '<input name="filter" value="' + id + '">';
    }
    form.action = "http://localhost:8080/";
    form.style.display = "none"
    document.body.append(form);
    form.submit();
}