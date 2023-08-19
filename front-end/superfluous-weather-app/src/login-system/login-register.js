function signup(username, password) {
  fetch("/api/signup", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username: username, password: password }),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.message === "success") {
        console.log("Signup successful");
      } else {
        console.log("Signup failed");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

function signin(username, password) {
  fetch("/api/signin", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username: username, password: password }),
  })
    .then((response) => {
      if (response.status === 200) {
        return response.json();
      } else {
        throw new Error("Authentication failed");
      }
    })
    .then((data) => {
      if (data.message === "success") {
        console.log("Signin successful");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}
