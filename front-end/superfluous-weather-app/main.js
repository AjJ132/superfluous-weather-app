// Select the button using its id
let forecastButton = document.getElementById("forecast-button");
let testButtton = document.getElementById("test-button");

var previousSearch = "";

window.onSearchButtonClick = onSearchButtonClick;
window.logout = logout;

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

export function onSearchButtonClick() {
  //grab text from search button
  let search = document.getElementById("locationTB").value;

  //verify that search isnt empty
  if (search == "") {
    alert("Please enter a location");
    return;
  }
  //Call API
  getForecast(search);
}

async function getForecast(search) {
  try {
    let url = "/api/weather-forecast?location=" + encodeURIComponent(search);
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const data = await response.json();

    console.log("Forecast successful");
    console.log(data);

    // Extract the meta information
    const meta = data.meta;
    const responseTime = meta.responseTime;
    const fromCache = meta.fromCache;

    // Display meta information alongside the weather data
    document.getElementById("data-results").innerHTML = `<pre>${JSON.stringify(
      data,
      null,
      2
    )}</pre>`;
    document.getElementById(
      "responseTime"
    ).innerText = `Response Time: ${responseTime} ms`;
    document.getElementById("fromCache").innerText = `From Cache: ${fromCache}`;
  } catch (error) {
    console.error("There was an error fetching the API:", error);
  }
}

async function logout() {
  //remove token from local storage
  localStorage.removeItem("token");

  //redirect to login page
  window.location.href = "/login.html";
}
