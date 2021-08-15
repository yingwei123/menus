function login(){
    userName = document.getElementById("userName").value
    password = document.getElementById("password").value

    const data ={
        UserName:userName,
        Password:password
    }

    fetch('/login', {
    method: 'POST', // or 'PUT'
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
    })
    .then(response => response.status)
    .then(data => {
        if(data == 303|| data ==200){
            location.reload()
            return
        }
        document.getElementById("userName").value = ""
        document.getElementById("password").value = ""
        document.getElementById("error").innerHTML ="Incorrect Credentials"
        return
    })
    .catch((error) => {
    console.error('Error:', error);
    });
}