async function addItem() {
    var fileName=document.getElementById('image-upload-input').files[0];
    const itemPrice = document.getElementById("iprice").value
    const ingredients = document.getElementById("ingredients").value
    const description = document.getElementById("idescription").value
    const itemName = document.getElementById("iname").value
    const formData = new FormData();
    formData.append("image-upload-input", fileName);
    formData.append("item-price", itemPrice);
    formData.append("item-ingredients", ingredients);
    formData.append("item-description", description);
    formData.append("item-name", itemName);

    const options = {
        method: 'POST',
        body: formData,
    }

    try {
       await fetch('/dashboard/add-menu-item', options).then(response =>response.json() ).then((data) =>{
            console.log(data)
        })
        
    } catch (err) {
        console.error(err);
    }

 
}