let token = localStorage.getItem("token");
let isRegisterMode = false;

const authSection = document.getElementById("authSection");
const dashboard = document.getElementById("dashboard");
const logoutBtn = document.getElementById("logoutBtn");

const authTitle = document.getElementById("authTitle");
const authBtn = document.getElementById("authBtn");
const switchAuth = document.getElementById("switchAuth");
const switchText = document.getElementById("switchText");
const authMessage = document.getElementById("authMessage");

switchAuth.addEventListener("click", () => {
    isRegisterMode = !isRegisterMode;

    authTitle.textContent = isRegisterMode ? "Register" : "Login";
    authBtn.textContent = isRegisterMode ? "Register" : "Login";

    switchText.innerHTML = isRegisterMode
        ? 'Already have an account? <button id="switchAuth" class="link">Login</button>'
        : 'Don\'t have an account? <button id="switchAuth" class="link">Register</button>';

    document.getElementById("switchAuth").addEventListener("click", () => {
        switchAuth.click();
    });

    authMessage.textContent = "";
});

authBtn.addEventListener("click", async () => {
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;

    if (!email || !password) {
        showMessage(authMessage, "Email and password are required.", true);
        return;
    }

    const endpoint = isRegisterMode
        ? "/auth/register"
        : "/auth/login";

    try {
        const response = await fetch(endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email,
                password
            })
        });

        const data = await response.json();

        if (!response.ok) {
            showMessage(authMessage, data.error || "Request failed.", true);
            return;
        }

        if (isRegisterMode) {
            showMessage(authMessage, "Registration successful. You can now log in.", false);
            isRegisterMode = false;
            authTitle.textContent = "Login";
            authBtn.textContent = "Login";
            return;
        }

        token = data.token;
        localStorage.setItem("token", token);

        showDashboard();
        loadTickets();

    } catch (error) {
        showMessage(authMessage, "Could not connect to the server.", true);
    }
});

document.getElementById("createTicketBtn").addEventListener("click", createTicket);

document.getElementById("refreshBtn").addEventListener("click", loadTickets);

logoutBtn.addEventListener("click", logout);

function showDashboard() {
    authSection.classList.add("hidden");
    dashboard.classList.remove("hidden");
    logoutBtn.classList.remove("hidden");
}

function showAuth() {
    authSection.classList.remove("hidden");
    dashboard.classList.add("hidden");
    logoutBtn.classList.add("hidden");
}

async function createTicket() {
    const title = document.getElementById("ticketTitle").value.trim();
    const description = document.getElementById("ticketDescription").value.trim();
    const message = document.getElementById("ticketMessage");

    if (!title || !description) {
        showMessage(message, "Title and description are required.", true);
        return;
    }

    try {
        const response = await fetch("/tickets", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                title,
                description
            })
        });

        const data = await response.json();

        if (!response.ok) {
            showMessage(message, data.error || "Could not create ticket.", true);
            return;
        }

        document.getElementById("ticketTitle").value = "";
        document.getElementById("ticketDescription").value = "";

        showMessage(message, "Ticket created successfully.", false);

        loadTickets();

    } catch (error) {
        showMessage(message, "Could not connect to the server.", true);
    }
}

async function loadTickets() {
    const ticketList = document.getElementById("ticketList");

    ticketList.innerHTML = "Loading...";

    try {
        const response = await fetch("/tickets", {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 401) {
                logout();
                return;
            }

            ticketList.innerHTML = `<p class="error">${data.error || "Failed to load tickets."}</p>`;
            return;
        }

        if (!data.length) {
            ticketList.innerHTML = "<p>No tickets yet.</p>";
            return;
        }

        ticketList.innerHTML = data.map(ticket => `
            <div class="ticket">
                <h3>${escapeHtml(ticket.title)}</h3>
                <p>${escapeHtml(ticket.description)}</p>
                <span class="status">${ticket.status}</span>

                <div>
                    ${getStatusButton(ticket)}
                </div>
            </div>
        `).join("");

    } catch (error) {
        ticketList.innerHTML = '<p class="error">Could not connect to the server.</p>';
    }
}

function getStatusButton(ticket) {
    if (ticket.status === "open") {
        return `
            <button onclick="updateStatus('${ticket.id}', 'in_progress')">
                Move to In Progress
            </button>
        `;
    }

    if (ticket.status === "in_progress") {
        return `
            <button onclick="updateStatus('${ticket.id}', 'closed')">
                Close Ticket
            </button>
        `;
    }

    return "";
}

async function updateStatus(ticketId, status) {
    try {
        const response = await fetch(`/tickets/${ticketId}/status`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                status
            })
        });

        const data = await response.json();

        if (!response.ok) {
            alert(data.error || "Could not update ticket.");
            return;
        }

        loadTickets();

    } catch (error) {
        alert("Could not connect to the server.");
    }
}

function logout() {
    token = null;
    localStorage.removeItem("token");
    showAuth();
}

function showMessage(element, message, error) {
    element.textContent = message;
    element.className = error ? "message error" : "message success";
}

function escapeHtml(value) {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}

if (token) {
    showDashboard();
    loadTickets();
}