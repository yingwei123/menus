function clickedItem(elem){
    var form = document.createElement("form");
    form.setAttribute("method", "get");
    form.setAttribute("action", "/dashboard/edit-menu-item/edit");

    var input = document.createElement("input");
    input.setAttribute("name","item")
    input.setAttribute("value",elem)
    input.setAttribute("type","hidden")

    form.appendChild(input)

    document.body.appendChild(form);
    form.submit();
}