let cookie = document.cookie;
let cookies = cookie.split(";");
let token = "";
cookies.forEach(function(val, ind) {
    if (val.includes("token")) {
        token = val.split("=")[1];
    }
});

let navigation = document.querySelector(".navigation");

function initNavigation(logged) {
    if (logged) {
        navigation.innerHTML = `<button class="logout-button" onclick="logout()">Logout</button>` 
    } else {
        navigation.innerHTML = `
        <a href="/registration">Registration</a>
        <a href="/login">Login</a>`;
    }
}

async function updateContent() {
    let p = document.querySelector('.data');
    fetch("http://localhost:8080/api", {
        headers: {
            Authorization: 'Bearer '+token
        }
      }).then( function(response) {
        return response.text();
      } ).then(function(data) {
        p.textContent = data;
      }).catch(function(){
        p.textContent = "";
      })
    
}


function logout() {
    document.cookie = "token=";
    token="";
    updateContent();
    initNavigation(false);
}

initNavigation(token);
updateContent();
