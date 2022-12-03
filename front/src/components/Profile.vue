<template>
  <div>
    <div>
      <h1>Hello {{ user.Username }}</h1>
      <div>
        <p>email</p>
        <input v-model="email" />
      </div>
      <button @click="update">Update</button>
      <button @click="test">test</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from "@vue/reactivity";
import { watch } from "@vue/runtime-core";
import axios from "axios";

const user = ref({});
const email = ref();

axios.get(`http://localhost:3000/user`, {
    headers : {Authorization: "Bearer "+ sessionStorage.getItem("token")}
  }).then(({data})=> user.value=data).catch(() =>sessionStorage.clear())

const update = async () => {
  if (email.value) user.value.Email = email.value;

  axios.patch(
    "http://localhost:3000/user",
    {
      password: "",
      email: email.value,
    },
    {
      headers: { Authorization: "Bearer " + sessionStorage.token },
    }
  );
};
</script>

<style scoped lang="sass">
.exercices
  display: grid
  grid-template: repeat(1fr, 6)
  gap: 3rem
  margin-top: 2rem
.child
  grid-colomn: span 2
</style>
