let form = document.forms["registration"];

cookies.forEach(function(val, ind) {
    if (val.includes("token") && val.split("=")[1]) {
      window.location.replace("http://localhost:8080/");
    }
});

let input_name = form.elements["name"];
let input_password = form.elements["password"];

let submit_button = document.querySelector(".registration__button");

async function sendForm(event) {
    event.preventDefault();
    let user = {
        "name": input_name.value,
        "password": input_password.value
    }
    let response = await fetch("http://localhost:8080/registrate", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(user)
      });
    let r = await response.text();
    window.location.replace("http://localhost:8080/login");
}

form.addEventListener("submit", sendForm)