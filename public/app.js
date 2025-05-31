const API_BASE = "https://localhost:3000";

const registerBtn = document.getElementById("register-btn");
const registerWebBtn = document.getElementById("register-webauthn-btn");
const loginBtn = document.getElementById("login-btn");
const loginWebBtn = document.getElementById("login-webauthn-btn");
const logoutBtn = document.getElementById("logout-btn");

function base64ToArrayBuffer(base64) {
  const binary = atob(base64.replace(/-/g, "+").replace(/_/g, "/"));
  const len = binary.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes.buffer;
}

function bufferToBase64url(buffer) {
  return btoa(String.fromCharCode(...new Uint8Array(buffer)))
    .replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function stringifyCredential(credential) {
  
  if (credential && credential.id && credential.rawId && credential.response) {
    const json = JSON.stringify({
      id: credential.id,
      rawId: bufferToBase64url(credential.rawId),
      type: credential.type,
      response: {
        clientDataJSON: bufferToBase64url(credential.response.clientDataJSON),
        attestationObject: credential.response.attestationObject
          ? bufferToBase64url(credential.response.attestationObject)
          : undefined,
        authenticatorData: credential.response.authenticatorData
          ? bufferToBase64url(credential.response.authenticatorData)
          : undefined,
        signature: credential.response.signature
          ? bufferToBase64url(credential.response.signature)
          : undefined,
        userHandle: credential.response.userHandle
          ? bufferToBase64url(credential.response.userHandle)
          : undefined,
      },
    });
    
    return json

  } else {
    throw new Error("Not a valid credential object");
  }
}


/**
 * 
 * @param {object} opts 
 * @returns 
 */
function preformatMakeCredReq(opts) {
  if ("Options" in opts) {
    opts = opts.Options
  }
  opts.publicKey.challenge = base64ToArrayBuffer(opts.publicKey.challenge);
  opts.publicKey.user.id = base64ToArrayBuffer(opts.publicKey.user.id);
  if (opts.publicKey.excludeCredentials) {
    opts.publicKey.excludeCredentials = opts.publicKey.excludeCredentials.map(
      (c) => ({
        id: base64ToArrayBuffer(c.id),
        type: c.type,
      }),
    );
  }
  return opts.publicKey;
}

function preformatGetAssertReq(opts) {
  opts.publicKey.challenge = base64ToArrayBuffer(opts.publicKey.challenge);
  if (opts.publicKey.allowCredentials) {
    opts.publicKey.allowCredentials = opts.publicKey.allowCredentials.map(
      (c) => ({
        id: base64ToArrayBuffer(c.id),
        type: c.type,
      }),
    );
  }
  return opts.publicKey;
}

function publicKeyCredentialToJSON(cred) {
  if (cred instanceof Array) {
    return cred.map((x) => publicKeyCredentialToJSON(x));
  }
  if (cred instanceof ArrayBuffer) {
    return bufferToBase64url(cred);
  }
  if (cred && typeof cred === "object") {
    const obj = {};
    for (let key in cred) {
      obj[key] = publicKeyCredentialToJSON(cred[key]);
    }
    return obj;
  }
  return cred;
}

registerBtn.addEventListener("click", async () => {
  const username = document.getElementById("reg-username").value.trim();
  const password = document.getElementById("reg-password").value;
  if (!username || !password) return;

  const res = await fetch(
    `${API_BASE}/register/password/${encodeURIComponent(username)}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ password }),
    },
  );
  const data = await res.json();
  if (res.ok) {
    localStorage.setItem("token", data.token);
    document.getElementById("register-msg").textContent = "Registered!";
    loadUser();
  } else {
    document.getElementById("register-msg").textContent = data.err || "Error";
  }
});

registerWebBtn.addEventListener("click", async () => {
  const username = document.getElementById("reg-username").value.trim();
  if (!username) return;

  const startRes = await fetch(
    `${API_BASE}/register/start/${encodeURIComponent(username)}`,
    { method: "POST" },
  );
  const startData = await startRes.json();
  if (!startRes.ok) {
    document.getElementById("register-msg").textContent =
      startData.err || "Error";
    return;
  }

  let credential;
  let json;
  try {
    credential = await navigator.credentials.create({
      publicKey: preformatMakeCredReq(startData),
    })


    
  } catch (err) {
    document.getElementById("register-msg").textContent =
      err.message || "WebAuthn error";
    return;
  }


  const endRes = await fetch(
    `${API_BASE}/register/end/${encodeURIComponent(username)}`,
    {
      method: "POST",
      body: stringifyCredential(credential),
    },
  );

  const endData = await endRes.json();
  if (endRes.ok) {
    localStorage.setItem("token", endData.token);
    document.getElementById("register-msg").textContent =
      "Registered with WebAuthn!";
    loadUser();
  } else {
    document.getElementById("register-msg").textContent =
      JSON.stringify(endData.err) || "Error";


    window.endData = endData
  }
});

loginBtn.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value.trim();
  const password = document.getElementById("login-password").value;
  if (!username || !password) return;

  const res = await fetch(
    `${API_BASE}/login/password/${encodeURIComponent(username)}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ password }),
    },
  );
  const data = await res.json();
  if (res.ok) {
    localStorage.setItem("token", data.token);
    document.getElementById("login-msg").textContent = "Logged in!";
    loadUser();
  } else {
    document.getElementById("login-msg").textContent = data.err || "Error";
  }
});

loginWebBtn.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value.trim();
  if (!username) return;

  const startRes = await fetch(
    `${API_BASE}/login/start/${encodeURIComponent(username)}`,
    { method: "POST" },
  );
  const startData = await startRes.json();
  if (!startRes.ok) {
    document.getElementById("login-msg").textContent = startData.err || "Error";
    return;
  }

  let credential;
  try {
    credential = await navigator.credentials.get({
      publicKey: preformatGetAssertReq(startData),
    });
  } catch (err) {
    document.getElementById("login-msg").textContent =
      err.message || "WebAuthn error";
    return;
  }

  const endRes = await fetch(
    `${API_BASE}/login/end/${encodeURIComponent(username)}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body:  stringifyCredential(credential),
    },
  );
  const endData = await endRes.json();
  if (endRes.ok) {
    localStorage.setItem("token", endData.token);
    document.getElementById("login-msg").textContent =
      "Logged in with WebAuthn!";
    loadUser();
  } else {
    document.getElementById("login-msg").textContent = endData.err || "Error";
  }
});

logoutBtn.addEventListener("click", () => {
  localStorage.removeItem("token");
  document.getElementById("user-section").style.display = "none";
});

async function loadUser() {
  const token = localStorage.getItem("token");
  if (!token) return;
  const res = await fetch(`${API_BASE}/user`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  if (res.ok) {
    const user = await res.json();
    document.getElementById("user-info").textContent = JSON.stringify(
      user,
      null,
      2,
    );
    document.getElementById("user-section").style.display = "block";
  }

}

loadUser();
