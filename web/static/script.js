document.addEventListener("DOMContentLoaded", () => {
  const token = localStorage.getItem("token");

  const signoutBtn = document.getElementById("signout-btn");
  if (signoutBtn) {
    signoutBtn.addEventListener("click", () => {
      localStorage.removeItem("token");
      window.location.reload();
    });
  }

  const createLink = document.getElementById("create-link");
  const loginLink = document.getElementById("login-link");
  const signupLink = document.getElementById("signup-link");
  if (createLink || loginLink || signupLink || signoutBtn) {
    if (token) {
      if (createLink) createLink.style.display = "inline";
      if (loginLink) loginLink.style.display = "none";
      if (signupLink) signupLink.style.display = "none";
      if (signoutBtn) signoutBtn.style.display = "inline";
    } else {
      if (createLink) createLink.style.display = "none";
      if (loginLink) loginLink.style.display = "inline";
      if (signupLink) signupLink.style.display = "inline";
      if (signoutBtn) signoutBtn.style.display = "none";
    }
  }

  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    if (token) {
      alert("You are already logged in.");
      window.location.href = "/";
      return;
    }
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const login = loginForm.login.value;
      const password = loginForm.password.value;
      try {
        const res = await fetch("/api/v1/auth/sign-in", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ login, password }),
        });
        const data = await res.json();
        if (res.ok) {
          localStorage.setItem("token", data.result);
          window.location.href = "/";
        } else {
          alert(data.error || "Login failed");
        }
      } catch (err) {
        alert("Error: " + err.message);
      }
    });
  }

  const signupForm = document.getElementById("signup-form");
  if (signupForm) {
    if (token) {
      alert("You are already logged in.");
      window.location.href = "/";
      return;
    }
    signupForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const login = signupForm.login.value;
      const password = signupForm.password.value;
      try {
        const res = await fetch("/api/v1/auth/sign-up", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ login, password }),
        });
        const data = await res.json();
        if (res.ok) {
          localStorage.setItem("token", data.result);
          window.location.href = "/";
        } else {
          alert(data.error || "Signup failed");
        }
      } catch (err) {
        alert("Error: " + err.message);
      }
    });
  }

  const createForm = document.getElementById("create-form");
  if (createForm) {
    if (!token) {
      alert("Please login first");
      window.location.href = "/login";
      return;
    }
    createForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const title = createForm.title.value;
      const description = createForm.description.value;
      const dateValue = createForm.date.value;
      const date = dateValue ? new Date(dateValue).toISOString() : "";
      const seats = parseInt(createForm.seats.value);
      const booking_ttl = createForm.booking_ttl.value;
      try {
        const res = await fetch("/api/v1/events", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            title,
            description,
            date,
            seats,
            booking_ttl,
          }),
        });
        const data = await res.json();
        if (res.ok) {
          window.location.href = "/";
        } else {
          alert(data.error || "Creation failed");
        }
      } catch (err) {
        alert("Error: " + err.message);
      }
    });
  }

  const eventIdElement = document.querySelector("[data-event-id]");
  if (eventIdElement) {
    const eventId = eventIdElement.dataset.eventId;
    const bookBtn = document.getElementById("book-btn");
    const confirmBtn = document.getElementById("confirm-btn");
    const freeSeatsP = document.getElementById("free-seats");

    if (!token) {
      if (bookBtn) bookBtn.style.display = "none";
      if (confirmBtn) confirmBtn.style.display = "none";
    }

    const updateFreeSeats = async () => {
      try {
        const res = await fetch(`/api/v1/events/${eventId}`);
        const data = await res.json();
        if (res.ok) {
          freeSeatsP.textContent = `Available Seats: ${data.result.seats || "Not available"}`;
        }
      } catch (err) {
        console.error("Failed to update free seats");
      }
    };

    updateFreeSeats();
    setInterval(updateFreeSeats, 5000);

    if (bookBtn) {
      bookBtn.addEventListener("click", async () => {
        try {
          const res = await fetch(`/api/v1/events/${eventId}/book`, {
            method: "POST",
            headers: { Authorization: `Bearer ${token}` },
          });
          const data = await res.json();
          if (res.ok) {
            alert("Booking created: " + data.result);
            updateFreeSeats();
          } else {
            alert(data.error || "Booking failed");
          }
        } catch (err) {
          alert("Error: " + err.message);
        }
      });
    }

    if (confirmBtn) {
      confirmBtn.addEventListener("click", async () => {
        try {
          const res = await fetch(`/api/v1/events/${eventId}/confirm`, {
            method: "POST",
            headers: { Authorization: `Bearer ${token}` },
          });
          const data = await res.json();
          if (res.ok) {
            alert("Booking confirmed");
            updateFreeSeats();
          } else {
            alert(data.error || "Confirmation failed");
          }
        } catch (err) {
          alert("Error: " + err.message);
        }
      });
    }
  }
});
