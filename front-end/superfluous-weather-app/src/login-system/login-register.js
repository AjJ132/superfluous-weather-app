let signinButton = document.getElementById("signin-button");
let registerButton = document.getElementById("signup-button");

signinButton.addEventListener("click", function () {
  // URL for the GET request
  let url = "/api/signin";
  console.log("signin button clicked");
  let username = document.getElementById("login-username").value;
  let password = document.getElementById("login-password").value;

  //attach username and password to url as JSON
  url += "?username=" + username + "&password=" + password;

  //print url to console
  console.log(url);

  // Fetch data from the URL
  fetch(url)
    .then((response) => response.json()) // Transform the data into json
    .then((data) => {
      //Check if the request was 401 or 200

      if (data.status == 401) {
        document.getElementById("login-result").innerHTML =
          "Incorrect username or password.";
      } else if (data.status == 200) {
        console.log("login successful");

        //redirect to home page
        window.location.href = "/";
      }
    })
    .catch((error) => {
      console.error("There was an error fetching the API:", error);
      document.getElementById("result").innerHTML =
        "An error occurred while fetching the data.";
    });
});

registerButton.addEventListener("click", function () {
  // URL for the POST request
  let url = "/api/signup";
  console.log("signup button clicked");
  let username = document.getElementById("register-username").value;
  let password = document.getElementById("register-password").value;

  // Create JSON payload
  let payload = {
    username: username,
    password: password
  };

  // Fetch data using POST method and sending JSON in the body
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload)
  })
  .then((response) => response.json()) // Transform the data into json
  .then((data) => {
    // Check if the request was 401 or 200
    if (data.status == 401) {
      document.getElementById("login-result").innerHTML = "Something went wrong";
    } else if (data.status == 200) {
      console.log("register successful");

      // Save the token to local storage
      localStorage.setItem('token', data.token);

      // Redirect to home page
      window.location.href = "/";
    }
  })
  .catch((error) => {
    console.error("There was an error fetching the API:", error);
    document.getElementById("result").innerHTML = "An error occurred while fetching the data.";
  });
});

