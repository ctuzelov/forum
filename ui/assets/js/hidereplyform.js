var IsAuthorized = `{{.IsAuthorized}}`;
if (!(/true/).test(IsAuthorized)) {
    var elements = document.getElementsByClassName("replyform");
    for (var i = 0; i < elements.length; i++) {
        elements[i].style.display = "none";
    }
}