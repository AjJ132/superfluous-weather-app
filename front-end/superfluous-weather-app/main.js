// Select the button using its id
let forecastButton = document.getElementById("forecast-button");
let testButtton = document.getElementById("test-button");

var previousSearch = "";

window.onSearchButtonClick = onSearchButtonClick;

if (
  !localStorage.getItem("token") ||
  isTokenExpired(localStorage.getItem("token"))
) {
  //window.location.href = "/login.html";

}

function isTokenExpired(token) {
  // Decode token and check expiration here.
  // This assumes the JWT token's payload includes an "exp" claim.
  const decodedToken = JSON.parse(atob(token.split(".")[1]));
  const expirationDate = decodedToken.exp * 1000; // Convert to milliseconds
  return expirationDate < new Date().getTime();
}

export function onSearchButtonClick(){
  //grab text from search button
  let search = document.getElementById("locationTB").value;

  //verify that search isnt empty
  if(search == ""){
    alert("Please enter a location");
    return;
  }

  //verify that search isnt the same as the previous search
  if(search == previousSearch){
   console.log("same search. Not doing anything");
    return;
  }

  //set previous search to current search
  previousSearch = search;

  //Call API
  getForecast(search);
};

async function getForecast(search) {
  try {
    let url = "/api/forecast?location=" + encodeURIComponent(search);
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error('Network response was not ok');
    }

    const data = await response.json();

    console.log("Forecast successful");
    console.log(data);

  } catch (error) {
    console.error("There was an error fetching the API:", error);
  }
}
