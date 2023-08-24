// Select the button using its id
let forecastButton = document.getElementById("forecast-button");
let testButtton = document.getElementById("test-button");

if (
  !localStorage.getItem("token") ||
  isTokenExpired(localStorage.getItem("token"))
) {
  window.location.href = "/login.html";

}

function isTokenExpired(token) {
  // Decode token and check expiration here.
  // This assumes the JWT token's payload includes an "exp" claim.
  const decodedToken = JSON.parse(atob(token.split(".")[1]));
  const expirationDate = decodedToken.exp * 1000; // Convert to milliseconds
  return expirationDate < new Date().getTime();
}

// Add event listener to the button
forecastButton.addEventListener("click", function () {
  // URL for the GET request
  let url = "http://localhost:8080/weather-forecast?location=hartwell , GA";

  // Fetch data from the URL
  fetch(url)
    .then((response) => response.json()) // Transform the data into json
    .then((data) => {
      // Do something with the data
      console.log(data);
    })
    .catch((error) => {
      // Handle the error
      console.error("Error:", error);
    });
});

testButtton.addEventListener("click", function () {
  // URL for the GET request
  let url = "/api/hello-world";
  console.log("test button clicked");

  // Fetch data from the URL
  fetch(url)
    .then((response) => response.json()) // Transform the data into json
    .then((data) => {
      //Get json data
      let text = data.message;
      document.getElementById("result").innerHTML = text;
    })
    .catch((error) => {
      console.error("There was an error fetching the API:", error);
      document.getElementById("result").innerHTML =
        "An error occurred while fetching the data.";
    });
});



export function onRegisterButtonClick() {

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
};

export function onSigninButtonClick() {
  // URL for the POST request
  let url = "/api/signin";
  console.log("signin button clicked");
  let username = document.getElementById("login-username").value;
  let password = document.getElementById("login-password").value;

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
      console.log("login successful");

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
};
