
cookies.forEach(function(val, ind) {
    if (val.includes("token") && val.split("=")[1]) {
      window.location.replace("http://localhost:8080/");
    }
});

let form = document.forms["login"];

let input_name = form.elements["name"];
let input_password = form.elements["password"];

let submit_button = document.querySelector(".login__button");

async function sendForm(event) {
    event.preventDefault();
    let user = {
        "name": input_name.value,
        "password": input_password.value
    }
    let response = await fetch("http://localhost:8080/log_in", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(user)
      });
    let resp = await response.json();
    let t = resp["Data"]; 
    console.log(t);
    document.cookie = "token="+t+"; path=/";
    window.location.replace("http://localhost:8080/");
}

form.addEventListener("submit", sendForm)