<template>
  <input v-model="username"/>
  <button @click="register">Register</button>
  <button @click="login">Login</button>
</template>

<script setup>
import HelloWorld from './components/HelloWorld.vue';
import axios from "axios"
import { ref } from 'vue';

const username = ref()


function bufferDecode(value) {
      return Uint8Array.from(atob(value), c => c.charCodeAt(0));
  }

function bufferEncode(value) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "")
}

const register = async () => {
  const { data } = await axios.post(`http://localhost:3000/register/start/${username.value}`)

  const options = data.Options
  options.publicKey.challenge = bufferDecode(options.publicKey.challenge);
  
  console.log(options);
  options.publicKey.user.id = bufferDecode(options.publicKey.user.id)

  if (options.publicKey.excludeCredentials) {
    for (var i = 0; i < options.publicKey.excludeCredentials.length; i++) {
      options.publicKey.excludeCredentials[i].id = bufferDecode(options.publicKey.excludeCredentials[i].id);
    }
  }

  options.publicKey.authenticatorSelection = {
    userVerification: "preferred" 
  }
    
  const credential =  await navigator.credentials.create({
    publicKey: options.publicKey,            
  })

  let attestationObject = credential.response.attestationObject;
  let clientDataJSON = credential.response.clientDataJSON;
  let rawId = credential.rawId;

  
  const cred= await axios.post(`http://localhost:3000/register/end/${username.value}`, {
    id: credential.id,
    rawId: bufferEncode(rawId),
    type: credential.type,
    response: {
      attestationObject: bufferEncode(attestationObject),
      clientDataJSON: bufferEncode(clientDataJSON)
  }})
  console.log(cred);

  sessionStorage.setItem("cred", cred.data.token)


}

const login = async () => {
  const { data } = await axios.post(`http://localhost:3000/login/start/${username.value}`)

  data.publicKey.challenge =bufferDecode(data.publicKey.challenge);
  data.publicKey.allowCredentials.forEach(function (listItem) {
            listItem.id = bufferDecode(listItem.id)
          });

  data.authenticatorSelection = {
    userVerification: "preferred" 
  }

  const assertion =  await navigator.credentials.get({
    publicKey: data.publicKey,            
  })

  let authData = assertion.response.authenticatorData;
  let clientDataJSON = assertion.response.clientDataJSON;
  let rawId = assertion.rawId;
  let sig = assertion.response.signature;
  let userHandle = assertion.response.userHandle;

  const cred= await axios.post(`http://localhost:3000/login/end/${username.value}`, {
    id: assertion.id,
    rawId: bufferEncode(rawId),
    type: assertion.type,
    response: {
      authenticatorData: bufferEncode(authData),
      clientDataJSON: bufferEncode(clientDataJSON),
      signature: bufferEncode(sig),
      userHandle: bufferEncode(userHandle),
    },
  })

  sessionStorage.setItem("token", cred.data.token)

  
}



</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
