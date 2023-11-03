/*

inspiration: 
https://dribbble.com/shots/2292415-Daily-UI-001-Day-001-Sign-Up

*/

document.addEventListener("DOMContentLoaded", function () {
  // this function runs when the DOM is ready, i.e. when the document has been parsed
  console.log("Dom is ready");

  let form = document.querySelector('form');
  let email = document.getElementById('email');
  let emailvalue = email.value;

  console.log(form);
  console.log(email)
  console.log(emailvalue);

  form.addEventListener('submit', (e) => {
    e.preventDefault();
    return false;
  });
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
    .then(response => response.json())
    .then(data => {
      // Handle the response from the server
      console.log(data);
    })
    .catch(error => {
      // Handle any errors that occurred during the request
      console.error(error);
    });
}

function OnSignIn() {
  const username = document.getElementById('username').value;
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  userData = {};
  userData['username'] = username;
  userData['email'] = email;
  userData['password'] = password;

  console.log(userData);

  fetch('/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  })
    .then(response => response.json())
    .then(data => {
      // Handle the response from the server
      console.log(data);
    })
    .catch(error => {
      // Handle any errors that occurred during the request
      console.error(error);
    });
}