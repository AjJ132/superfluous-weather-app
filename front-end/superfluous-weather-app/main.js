import './style.css'

// Select the button using its id
let forecastButton = document.getElementById('forecast-button');

// Add event listener to the button
forecastButton.addEventListener('click', function() {
  // URL for the GET request
  let url = 'http://localhost:8080/weather-forecast?location=hartwell , GA';

  // Fetch data from the URL
  fetch(url)
    .then(response => response.json()) // Transform the data into json
    .then(data => {
      // Do something with the data
      console.log(data);
    })
    .catch(error => {
      // Handle the error
      console.error('Error:', error);
    });
});


