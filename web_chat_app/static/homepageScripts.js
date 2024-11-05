async function Register() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  const response = await fetch("/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  const result = await response.json();
  alert(result.message);
}

async function Login() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  const response = await fetch("/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  const result = await response.json();
  alert(result.message);

  if (response.ok) {
    // Redirect to chat page if login is successful
    window.location.href = `./chatBox.html?Username=${username}`;
  } else {
    alert("Login failed or response not OK");
  }
}
