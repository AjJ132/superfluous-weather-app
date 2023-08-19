let signinButton = document.getElementById("signin-button");
let registerButton = document.getElementById("signup-button");

signinButton.addEventListener("click", function () {
  // URL for the GET request
  let url = "/api/signin";
  console.log("signin button clicked");

  // Fetch data from the URL
  fetch(url)
    .then((response) => response.json()) // Transform the data into json
    .then((data) => {
      //Check if the request was 401 or 200

      if (data.status == 401) {
      } else if (data.status == 200) {
      }
    })
    .catch((error) => {
      console.error("There was an error fetching the API:", error);
      document.getElementById("result").innerHTML =
        "An error occurred while fetching the data.";
    });
});

registerButton.addEventListener("click", function () {
  // URL for the GET request
  let url = "/api/signup";
  console.log("signup button clicked");

  // Fetch data from the URL
  fetch(url)
    .then((response) => response.json()) // Transform the data into json
    .then((data) => {
      //Check if the request was 401 or 200

      if (data.status != 200) {
      } else {
      }
    })
    .catch((error) => {
      console.error("There was an error fetching the API:", error);
      document.getElementById("result").innerHTML =
        "An error occurred while fetching the data.";
    });
});
