<template>

<div>
  <h1> Hello {{user.Username}}</h1>
  <div>
    <p> password</p>
    <input v-model="password">
  </div>
  <div>
    <p> email</p>
    <input v-model="email">
  </div>
  <button @click="update"> Update </button>
</div>
  
</template>

<script setup>
import { ref } from "@vue/reactivity";
import axios from "axios";

var user = ref({})
var password = ref()
var email = ref()

axios.get(`http://localhost:3000/user`, {
    headers : {Authorization: "Bearer "+ sessionStorage.getItem("token")}
  }).then(({data})=>{
    user.value = data
  })



const update = async () => {

  if(email.value)
    user.value.Email = email.value
  if(password.value)
    user.value.Password = password.value

  console.log(user.value);
  return
  axios.patch(`http://localhost:3000/user`, {
    headers : {Authorization: "Bearer "+ sessionStorage.getItem("token")},
    data: user
  })

}


</script>

<style>

</style>