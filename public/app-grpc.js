// grpc.js — Front statique, zéro lib

// ---- Config ----
const API_BASE = "https://localhost:50051"; // <- ton grpc-gateway HTTP/JSON
const AUTH_SVC = "auth.v1.AuthService";

// ---- DOM ----
const registerBtn = document.getElementById("register-btn");
const registerWebBtn = document.getElementById("register-webauthn-btn");
const loginBtn = document.getElementById("login-btn");
const loginWebBtn = document.getElementById("login-webauthn-btn");
const logoutBtn = document.getElementById("logout-btn");
const listUsersbtn = document.getElementById("list-users-btn");
const listUsers = document.getElementById("list-users");

// ---- Utils base64 / buffers ----
// gRPC-Gateway sérialise 'bytes' en base64 (RFC4648 non URL-safe)
function base64ToArrayBuffer(b64) {
  // normalise padding
  const pad = "=".repeat((4 - (b64.length % 4)) % 4);
  const normalized = (b64 + pad).replace(/-/g, "+").replace(/_/g, "/");
  const bin = atob(normalized);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return bytes.buffer;
}
function arrayBufferToBase64(buf) {
  const bytes = new Uint8Array(buf);
  let bin = "";
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
  return btoa(bin);
}

// Pour WebAuthn (certaines structures utilisent base64url)
function bufferToBase64url(buffer) {
  return arrayBufferToBase64(buffer).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
}

// ---- RPC helper (HTTP/JSON via gateway) ----
async function rpc(service, method, body) {
  const res = await fetch(`${API_BASE}/${service}/${method}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include", // utile si tu passes par cookies HttpOnly
    body: JSON.stringify(body ?? {}),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const msg = data?.message || data?.error || JSON.stringify(data);
    throw new Error(msg || `RPC ${method} failed`);
  }
  return data;
}

// ---- WebAuthn helpers (natifs, zéro lib) ----
function stringifyCredential(credential) {
  if (!credential || !credential.id || !credential.rawId || !credential.response) {
    throw new Error("Not a valid credential object");
  }
  const obj = {
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
  };
  return JSON.stringify(obj);
}

// Les réponses Start* RPC renvoient du JSON WebAuthn compacté en bytes (base64)
function preformatMakeCredReqFromRPC(startResp) {
  // startResp.creationOptionsJson : base64 -> JSON -> objet
  const json = new TextDecoder().decode(base64ToArrayBuffer(startResp.creationOptionsJson));
  const opts = JSON.parse(json);
  // Convertir les champs string(base64url) -> ArrayBuffer
  opts.publicKey.challenge = base64ToArrayBuffer(opts.publicKey.challenge);
  opts.publicKey.user.id = base64ToArrayBuffer(opts.publicKey.user.id);
  if (opts.publicKey.excludeCredentials) {
    opts.publicKey.excludeCredentials = opts.publicKey.excludeCredentials.map(c => ({
      ...c,
      id: base64ToArrayBuffer(c.id),
    }));
  }
  return opts.publicKey;
}
function preformatGetAssertReqFromRPC(startResp) {
  const json = new TextDecoder().decode(base64ToArrayBuffer(startResp.requestOptionsJson));
  const opts = JSON.parse(json);
  opts.publicKey.challenge = base64ToArrayBuffer(opts.publicKey.challenge);
  if (opts.publicKey.allowCredentials) {
    opts.publicKey.allowCredentials = opts.publicKey.allowCredentials.map(c => ({
      ...c,
      id: base64ToArrayBuffer(c.id),
    }));
  }
  return opts.publicKey;
}

// ---- Handlers ----

// Register (password) -> AuthService.RegisterPassword
registerBtn.addEventListener("click", async () => {
  const username = document.getElementById("reg-username").value.trim();
  const password = document.getElementById("reg-password").value;
  const email = ""; // adapte si ton proto l'exige
  if (!username || !password) return;

  try {
    const out = await rpc(AUTH_SVC, "RegisterPassword", { username, email, password });
    // selon ton impl, tu peux aussi recevoir un token ici ; sinon login derrière
    document.getElementById("register-msg").textContent = "Registered!";
    // Si ton serveur renvoie un token:
    if (out.accessToken) localStorage.setItem("token", out.accessToken);
    loadUser();
  } catch (e) {
    document.getElementById("register-msg").textContent = e.message || "Error";
  }
});

// Register (WebAuthn) -> StartRegistration + FinishRegistration
registerWebBtn.addEventListener("click", async () => {
  const username = document.getElementById("reg-username").value.trim();
  if (!username) return;

  try {
    const start = await rpc(AUTH_SVC, "StartRegistration", { userId: username });
    const publicKey = preformatMakeCredReqFromRPC(start);

    const credential = await navigator.credentials.create({ publicKey });
    const jsonString = stringifyCredential(credential);

    // credential_json est un champ bytes => encoder UTF-8 -> base64
    const b64 = arrayBufferToBase64(new TextEncoder().encode(jsonString));
    const end = await rpc(AUTH_SVC, "FinishRegistration", { credentialJson: b64 });

    if (end.accessToken) localStorage.setItem("token", end.accessToken);
    document.getElementById("register-msg").textContent = "Registered with WebAuthn!";
    loadUser();
  } catch (e) {
    document.getElementById("register-msg").textContent = e.message || "Error";
  }
});

// Login (password) -> AuthService.LoginPassword
loginBtn.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value.trim();
  const password = document.getElementById("login-password").value;
  if (!username || !password) return;

  try {
    const out = await rpc(AUTH_SVC, "LoginPassword", { usernameOrEmail: username, password });
    if (out.accessToken) localStorage.setItem("token", out.accessToken);
    document.getElementById("login-msg").textContent = "Logged in!";
    loadUser();
  } catch (e) {
    document.getElementById("login-msg").textContent = e.message || "Error";
  }
});

// Login (WebAuthn) -> StartLogin + FinishLogin
loginWebBtn.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value.trim();
  if (!username) return;

  try {
    const start = await rpc(AUTH_SVC, "StartLogin", { usernameOrEmail: username });
    const publicKey = preformatGetAssertReqFromRPC(start);

    const assertion = await navigator.credentials.get({ publicKey });
    const jsonString = stringifyCredential(assertion);

    const b64 = arrayBufferToBase64(new TextEncoder().encode(jsonString));
    const end = await rpc(AUTH_SVC, "FinishLogin", { username, assertionJson: b64 });

    if (end.accessToken) localStorage.setItem("token", end.accessToken);
    document.getElementById("login-msg").textContent = "Logged in with WebAuthn!";
    loadUser();
  } catch (e) {
    document.getElementById("login-msg").textContent = e.message || "Error";
  }
});

// Logout (client-side)
logoutBtn.addEventListener("click", () => {
  localStorage.removeItem("token");
  document.getElementById("user-section").style.display = "none";
  document.getElementById("users-section").style.display = "none";
});

// ---- /user & /user/all (temporaire) ----
// Tant que ton service User n'est pas exposé en gRPC+gateway, on garde ces endpoints.
// Remplace-les plus tard par p.ex. user.v1.UserService/GetMe & /ListUsers.
async function loadUser() {
  const token = localStorage.getItem("token");
  if (!token) return;
  const res = await fetch(`${API_BASE}/user`, {
    headers: { Authorization: `Bearer ${token}` },
    credentials: "include",
  });
  if (res.ok) {
    const user = await res.json();
    document.getElementById("user-info").textContent = JSON.stringify(user, null, 2);
    document.getElementById("user-section").style.display = "block";
  }
}

listUsersbtn.addEventListener("click", async () => {
  const token = localStorage.getItem("token");
  if (!token) return;
  const res = await fetch(`${API_BASE}/user/all`, {
    headers: { Authorization: `Bearer ${token}` },
    credentials: "include",
  });
  if (res.ok) {
    const users = await res.json();
    listUsers.textContent = JSON.stringify(users, null, 2);
    document.getElementById("users-section").style.display = "block";
  }
});

loadUser();
