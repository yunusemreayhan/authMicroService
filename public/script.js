/*

inspiration: 
https://dribbble.com/shots/2292415-Daily-UI-001-Day-001-Sign-Up

*/

document.addEventListener("DOMContentLoaded", function () {
  // this function runs when the DOM is ready, i.e. when the document has been parsed
  console.log("Dom is ready");

  let email = document.getElementById('email');
  if (email) {
    let emailvalue = email.value;
    console.log(emailvalue);
  }
  console.log(email)

  let form = document.querySelector('form');
  if (form) {
    console.log(form);
    form.addEventListener('submit', (e) => {
      e.preventDefault();
      return false;
    });
  }
}
);

function OnRegister() {
  const username = document.getElementById('username').value;
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  userData = {};
  userData['username'] = username;
  userData['email'] = email;
  userData['password'] = password;

  console.log(userData);

  fetch('/api/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  })
    .then(response => {
      try {
        console.log("response : ", response);
        if (response.status != 200) {
          response.text().then(text => {
            console.log("text : ", text);
            alert(text);
          })
        }
        console.log("reponse body : ", response.body);
        console.log("response status : ", response.status);
        return response.json()
      } catch (error) {
        console.log("response return error : ", error);
        return response
      }
    })
    .then(data => {
      // Handle the response from the server
      // if data type is json
      if (typeof data === 'object') {
        console.log(data);
        window.location.href = "./login.html";
      } else {
        console.log(typeof data)
        alert(data);
      }
    })
    .catch(error => {
      // Handle any errors that occurred during the request
      console.error(error);
    });
}

function OnSignIn() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  userData = {};
  userData['username'] = username;
  userData['password'] = password;

  console.log(userData);

  fetch('/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  })
    .then(response => {
      try {
        console.log("response : ", response);
        if (response.status != 200) {
          response.text().then(text => {
            console.log("text : ", text);
            alert(text);
          })
        }
        console.log("reponse body : ", response.body);
        console.log("response status : ", response.status);
        return response.json()
      } catch (error) {
        console.log("response return error : ", error);
        return response
      }
    })
    .then(data => {
      // Handle the response from the server
      // if data type is json
      if (typeof data === 'object') {
        console.log(data);
        localStorage.setItem('tokenAuthMicroService', data["voucher"]);
        window.location.href = "./index.html";
      } else {
        console.log(typeof data)
        alert(data);
      }
    })
    .catch(error => {
      // Handle any errors that occurred during the request
      console.error(error);
    });
}

function OnVerify() {
  console.log("Verify called!");
  userData = {};
  userData['voucher'] = localStorage.getItem('tokenAuthMicroService');

  console.log("verification data : ", userData);
  fetch('/api/verify', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  })
    .then(response => {
      try {
        console.log("response : ", response);
        if (response.status != 200) {
          response.text().then(text => {
            console.log("login verification failed : ", text);
            alert("login verification failed : " + text);
            window.location.href = "./login.html";
          })
        } else {
          alert ("login verification success");
        }
        console.log("reponse body : ", response.body);
        console.log("response status : ", response.status);
        return response.json()
      } catch (error) {
        console.log("response return error : ", error);
        return response
      }
    })
}